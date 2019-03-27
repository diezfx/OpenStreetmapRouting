package dijkstra

import (
	"OpenStreetmapRouting/data"
	"errors"
	"math"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

func GetRouteWithStations(graph *data.GraphProd, stations *data.GasStations, start data.Coordinate, end data.Coordinate, rangeKm float64) (*data.NodeRoute, []data.Node, error) {

	logrus.Debug(start, end)

	logrus.Infof("Find nodes close to Node")
	startTime := time.Now()
	startNode := graph.Grid.FindNextNode(start.Lat, start.Lon, true)
	endNode := graph.Grid.FindNextNode(end.Lat, end.Lon, true)
	rangeCm := int64(rangeKm * 1000 * 100)

	gridTime := time.Since(startTime)

	logrus.Infof("Dijkstra started")

	result, stationOnRoute, err := CalcStationDijkstra(graph, stations, *startNode, *endNode, rangeCm)
	dijkstraTime := time.Since(startTime) - gridTime
	endTime := time.Since(startTime)
	logrus.WithFields(logrus.Fields{
		"Time for Gridsearch": gridTime,
		"Time for dijkstra":   dijkstraTime,
		"Total time":          endTime}).Info("Dijkstra ended")

	return result, stationOnRoute, err

}

// CalcStationDijkstra basic idea:
// Calculate dijkstra to goal, if reachable
//	if not take reachable station that is closest to target(air distance)
func CalcStationDijkstra(g *data.GraphProd, gasStations *data.GasStations, start data.Node, target data.Node, rangeCm int64) (route *data.NodeRoute, stations []data.Node, err error) {

	stations = make([]data.Node, 0)

	minWay := make([]*data.Node, 0)

	route = &data.NodeRoute{Route: minWay, TotalCost: 0}

	partialStart := start
	var partialTarget data.Node
	var partialRoute *data.NodeRoute

	for partialTarget != target {

		wayCostEdges, err := CalcDijkstraToMany(g, partialStart)

		if err != nil {
			logrus.Fatalf("Error: %s", err.Error())
			return nil, nil, err
		}

		if wayCostEdges[target.ID].Cost == math.MaxInt64 {
			return nil, nil, errors.New("no way found")
		}

		// is goal reachable?

		if goalReachable(target, wayCostEdges, rangeCm) {
			partialTarget = target
			partialRoute, _ = findWayToGoal(partialStart, target, g, wayCostEdges)

			//reverse route
			for i := len(partialRoute.Route) - 1; i >= 0; i-- {

				route.Route = append(route.Route, partialRoute.Route[i])

			}
			route.TotalCost += partialRoute.TotalCost

		} else {

			minStation, err := getMinDistanceStation(g.Nodes, wayCostEdges, gasStations, start, target, rangeCm)
			if err != nil {
				return nil, nil, err
			}

			// add way to this station to route and station to list
			partialTarget = g.Nodes[minStation]
			partialRoute, _ = findWayToGoal(partialStart, partialTarget, g, wayCostEdges)
			partialStart = g.Nodes[minStation]
			//add route in reverse
			for i := len(partialRoute.Route) - 1; i >= 1; i-- {
				route.Route = append(route.Route, partialRoute.Route[i])
			}
			route.TotalCost += partialRoute.TotalCost
			stations = append(stations, partialTarget)
		}

	}
	return route, stations, nil

}

// save partialRoute,stations(so far),partialStart
type snapshot struct {
	stations     []data.Node
	partialRoute data.NodeRoute
	partialStart data.Node
}

// CalcStationDijkstraSnapshots basic idea:
// Calculate dijkstra to goal, if reachable
//	if not take reachable station that is closest to target(air distance)
// in case of no found way to target take other station
// if that doesnt work go back one station and try from there
func CalcStationDijkstraSnapshots(g *data.GraphProd, gasStations *data.GasStations, start data.Node, target data.Node, rangeCm int64) (route *data.NodeRoute, stations []data.Node, err error) {

	stations = make([]data.Node, 0)

	minWay := make([]*data.Node, 0)

	route = &data.NodeRoute{Route: minWay, TotalCost: 0}

	partialStart := start
	var partialTarget data.Node
	var partialRoute *data.NodeRoute

	// already visited stations don't take again
	//blacklist := make(map[int]struct{}, 0)
	//snapshots := make([]snapshot, 0)

	for partialTarget != target {

		wayCostEdges, err := CalcDijkstraToMany(g, partialStart)

		if err != nil {
			logrus.Fatalf("Error: %s", err.Error())
			return nil, nil, err
		}

		if wayCostEdges[target.ID].Cost == math.MaxInt64 {
			return nil, nil, errors.New("no way found")
		}

		// is goal reachable?

		if goalReachable(target, wayCostEdges, rangeCm) {
			partialTarget = target
			partialRoute, _ = findWayToGoal(partialStart, target, g, wayCostEdges)

			//reverse route
			for i := len(partialRoute.Route) - 1; i >= 0; i-- {

				route.Route = append(route.Route, partialRoute.Route[i])
			}
			route.TotalCost += partialRoute.TotalCost
		} else {

			minStations, err := getMinDistanceStations(g.Nodes, wayCostEdges, gasStations, start, target, rangeCm, 5)
			if len(minStations) <= 0 {
				return nil, nil, err
			}

			//create snapshots

			// add way to this station to route and station to list
			partialTarget = g.Nodes[minStations[0]]
			partialRoute, _ = findWayToGoal(partialStart, partialTarget, g, wayCostEdges)
			partialStart = g.Nodes[minStations[0]]
			//add route in reverse
			for i := len(partialRoute.Route) - 1; i >= 1; i-- {
				route.Route = append(route.Route, partialRoute.Route[i])
			}
			route.TotalCost += partialRoute.TotalCost
			stations = append(stations, partialTarget)
		}

	}
	return route, stations, nil

}

// getMinDistanceStations returns the count closest stations(air distance)
func getMinDistanceStations(nodes []data.Node, wayCostEdges []data.Edge, stations *data.GasStations, start data.Node, target data.Node, rangeCm int64, count int) ([]int64, error) {

	//look for closest station to target(airdistance)

	type idDist struct {
		id   int
		dist float64
	}

	distTable := make([]idDist, len(wayCostEdges))

	for i := range wayCostEdges {
		distTable[i] = idDist{id: i, dist: CalcEuclidDist(nodes[i].Lat, nodes[i].Lon, target.Lat, target.Lon)}
	}

	//todo:not worth to sort
	sort.Slice(distTable, func(i, j int) bool {
		if wayCostEdges[i].Cost < rangeCm && wayCostEdges[j].Cost < rangeCm {
			return distTable[i].dist < distTable[j].dist
		} else if wayCostEdges[i].Cost < rangeCm && wayCostEdges[j].Cost >= rangeCm {
			return true
		} else if wayCostEdges[i].Cost >= rangeCm && wayCostEdges[j].Cost < rangeCm {
			return false
		} else {
			// both unreachable; doesnt matter
			return distTable[i].dist < distTable[j].dist
		}
	})

	clostestStations := make([]int64, 0)

	for i := 0; i < count; i++ {
		if wayCostEdges[distTable[i].id].Cost < rangeCm {
			clostestStations = append(clostestStations, int64(distTable[i].id))
		} else {
			return clostestStations, errors.New("no more reachable stations")
		}

	}

	return clostestStations, nil

}

func getMinDistanceStation(nodes []data.Node, wayCostEdges []data.Edge, stations *data.GasStations, start data.Node, target data.Node, rangeCm int64) (int64, error) {

	//look for closest station to target(airdistance)

	var minStation int64 = -1
	var minDistance = math.MaxFloat64

	logrus.Debug("another station on the way is necessaray")

	for stationID := range stations.Stations {

		euclidDist := CalcEuclidDist(nodes[stationID].Lat, nodes[stationID].Lon, target.Lat, target.Lon)

		if wayCostEdges[stationID].Cost < rangeCm && euclidDist < minDistance {

			minStation = stationID
			minDistance = CalcEuclidDist(nodes[stationID].Lat, nodes[stationID].Lon, target.Lat, target.Lon)
		}

	}
	if minStation == -1 {
		return -1, errors.New("no good station found")

	}
	return minStation, nil

}

func goalReachable(goal data.Node, wayCostEdges []data.Edge, rangeCm int64) bool {

	if wayCostEdges[goal.ID].Cost > rangeCm {
		return false
	}
	return true
}

func CalcEuclidDist(x1, x2, y1, y2 float64) float64 {

	d1 := math.Abs(x1 - y1)
	d2 := math.Abs(x2 - y2)
	return math.Sqrt(math.Pow(d1, 2) + math.Pow(d2, 2))
}

func StationsReachable(graph *data.GraphProd, start data.Coordinate) (Reachable []*data.Node, Unreachable []*data.Node) {

	stations := data.GetFuelStations()

	Reachable = make([]*data.Node, 0)
	Unreachable = make([]*data.Node, 0)
	errorCount := 0

	startNode := graph.Grid.FindNextNode(start.Lat, start.Lon, false)
	goalCosts, _ := CalcDijkstraToMany(graph, *startNode)

	for _, station := range stations.Stations {

		goalNode := graph.Grid.FindNextNode(station.Lat, station.Lon, false)

		if goalCosts[goalNode.ID].Cost >= math.MaxInt64 {
			Unreachable = append(Unreachable, goalNode)
			errorCount++
		} else {
			Reachable = append(Reachable, goalNode)
		}
	}

	if errorCount > 0 {
		logrus.Errorf("Expected all stations reachable got %d errors", errorCount)

	}
	return

}
