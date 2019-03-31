package dijkstra

import (
	"errors"
	"math"
	"sort"
	"time"

	"github.com/diezfx/OpenStreetmapRouting/config"
	"github.com/diezfx/OpenStreetmapRouting/data"

	"github.com/sirupsen/logrus"
)

func GetRouteWithStations(graph *data.GraphProd, stations *data.GasStations, start data.Coordinate, end data.Coordinate, rangeKm float64, config *config.Config) (*data.NodeRoute, []data.Node, error) {

	logrus.Debug(start, end)

	logrus.Infof("Find nodes close to Node")
	startTime := time.Now()
	startNode := graph.Grid.FindNextNode(start.Lat, start.Lon, true)
	endNode := graph.Grid.FindNextNode(end.Lat, end.Lon, true)
	rangeCm := int64(rangeKm * 1000 * 100)

	gridTime := time.Since(startTime)

	logrus.Infof("Dijkstra started")

	result, stationOnRoute, err := CalcStationDijkstraSnapshots(graph, stations, *startNode, *endNode, rangeCm, config)
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

			minStation, err := getMinDistanceStation(g.Nodes, wayCostEdges, gasStations, partialStart, target, rangeCm)
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

// CalcStationDijkstraSnapshots basic idea:
// Calculate dijkstra to goal, if reachable
//	if not take reachable station that is closest to target(air distance)
// in case of no found way to target take other station
// if that doesnt work go back one station and try from there
func CalcStationDijkstraSnapshots(g *data.GraphProd, gasStations *data.GasStations, start data.Node, target data.Node, rangeCm int64, config *config.Config) (route *data.NodeRoute, stations []data.Node, err error) {

	stations = make([]data.Node, 0)
	minWay := make([]*data.Node, 0)

	route = &data.NodeRoute{Route: minWay, TotalCost: 0}

	type S struct{}
	// already visited stations don't take again
	blacklist := make(map[int64]struct{}, 0)
	snapshotStack := data.NewStack()
	stationsVisitedCounter := 0

	startSnap := data.Snapshot{PartialStart: start, PartialRoute: *route, Stations: stations}
	snapshotStack.Push(startSnap)
	blacklist[start.ID] = S{}

	for {

		snap, err := snapshotStack.Pop()
		if err != nil {
			return nil, nil, err
		}
		partialStart := snap.PartialStart
		route = &snap.PartialRoute
		stations = snap.Stations

		wayCostEdges, err := CalcDijkstraToMany(g, partialStart)

		if err != nil {
			logrus.Fatalf("Error: %s", err.Error())
			return nil, nil, err
		}

		if wayCostEdges[target.ID].Cost == math.MaxInt64 {
			continue
		}

		// is goal reachable?

		if goalReachable(target, wayCostEdges, rangeCm) {

			partialRoute, _ := findWayToGoal(partialStart, target, g, wayCostEdges)

			//reverse route
			for i := len(partialRoute.Route) - 1; i >= 1; i-- {
				route.Route = append(route.Route, partialRoute.Route[i])
			}
			route.TotalCost += partialRoute.TotalCost
			return route, stations, nil
		}

		stationsVisitedCounter++
		if stationsVisitedCounter > config.DijkstraMaxStations {
			return nil, stations, errors.New("too many stations visited already")
		}

		minStations, err := getMinDistanceStations(g.Nodes, wayCostEdges, gasStations, target, rangeCm, config.DijkstraMaxstationsPerStep)

		for _, node := range minStations {
			logrus.Debug(wayCostEdges[node])
		}

		// no stations reachable -> go back to the previous station and try another one
		if len(minStations) <= 0 {
			continue
		}

		//save snapshots in  stack, continue with the closest one
		// add way to this station to route and station to list
		for i := len(minStations) - 1; i >= 0; i-- {
			node := minStations[i]

			if _, ok := blacklist[node]; ok == true {
				continue
			}

			partialTarget := g.Nodes[node]

			partialRoute, _ := findWayToGoal(partialStart, partialTarget, g, wayCostEdges)

			//add route in reverse

			newRouteList := route.Route
			for i := len(partialRoute.Route) - 1; i >= 1; i-- {
				newRouteList = append(newRouteList, partialRoute.Route[i])
			}
			logrus.Debug(len(newRouteList))
			newCost := route.TotalCost + partialRoute.TotalCost
			newRoute := data.NodeRoute{TotalCost: newCost, Route: newRouteList}
			newstations := append(stations, partialTarget)

			newPartialStart := g.Nodes[node]

			snapshot := data.Snapshot{PartialStart: newPartialStart, PartialRoute: newRoute, Stations: newstations}

			snapshotStack.Push(snapshot)

			blacklist[node] = S{}
		}

	}

}

// getMinDistanceStations returns the count closest stations(air distance)
func getMinDistanceStations(nodes []data.Node, wayCostEdges []data.Edge, stations *data.GasStations, target data.Node, rangeCm int64, count int) ([]int64, error) {

	distTable := make(map[int64]float64, len(stations.Stations))

	for i := range stations.Stations {
		distTable[i] = CalcEuclidDist(nodes[i].Lat, nodes[i].Lon, target.Lat, target.Lon)
	}

	clostestStations := make([]int64, 0)

	for stationID := range stations.Stations {

		if wayCostEdges[stationID].Cost < rangeCm {

			clostestStations = updateMinList(clostestStations, stationID, distTable, rangeCm, count)

		}

	}

	if len(clostestStations) < count {
		return clostestStations, errors.New("not enough stations found")
	}

	return clostestStations, nil

}

func updateMinList(minDistanceStations []int64, id int64, distTable map[int64]float64, rangeCm int64, count int) []int64 {

	if len(minDistanceStations) < count {

		minDistanceStations := append(minDistanceStations, id)

		sort.Slice(minDistanceStations, func(i, j int) bool {
			return distTable[minDistanceStations[i]] < distTable[minDistanceStations[j]]
		})

		return minDistanceStations
	}

	// check if the new distance is really smaller then replace biggest
	maxMinDistanceStation := minDistanceStations[len(minDistanceStations)-1]
	if distTable[maxMinDistanceStation] < distTable[id] {
		return minDistanceStations
	}

	minDistanceStations[len(minDistanceStations)-1] = id
	sort.Slice(minDistanceStations, func(i, j int) bool {
		return distTable[minDistanceStations[i]] < distTable[minDistanceStations[j]]
	})

	return minDistanceStations

}

func getMinDistanceStation(nodes []data.Node, wayCostEdges []data.Edge, stations *data.GasStations, start data.Node, target data.Node, rangeCm int64) (int64, error) {

	//look for closest station to target(airdistance)

	var minStation int64 = -1
	var minDistance = math.MaxFloat64

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

// CalcEuclidDist .
func CalcEuclidDist(x1, x2, y1, y2 float64) float64 {

	d1 := math.Abs(x1 - y1)
	d2 := math.Abs(x2 - y2)
	return math.Sqrt(math.Pow(d1, 2) + math.Pow(d2, 2))
}

// StationsReachable checks if all stations are smaller infinite value to reach
func StationsReachable(graph *data.GraphProd, start data.Coordinate, reachCm int64) (Reachable []*data.Node, Unreachable []*data.Node) {

	stations := data.GetFuelStations()

	Reachable = make([]*data.Node, 0)
	Unreachable = make([]*data.Node, 0)
	errorCount := 0

	startNode := graph.Grid.FindNextNode(start.Lat, start.Lon, false)
	goalCosts, _ := CalcDijkstraToMany(graph, *startNode)

	for _, station := range stations.Stations {

		goalNode := graph.Grid.FindNextNode(station.Lat, station.Lon, false)

		if goalCosts[goalNode.ID].Cost >= reachCm {
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
