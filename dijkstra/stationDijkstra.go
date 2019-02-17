package dijkstra

import (
	"OpenStreetmapRouting/data"
	"errors"
	"math"
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

	result, stationOnRoute, err := CalcStationDijkstra(graph, stations, startNode, endNode, rangeCm)
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
func CalcStationDijkstra(g *data.GraphProd, gasStations *data.GasStations, start *data.Node, target *data.Node, rangeCm int64) (route *data.NodeRoute, stations []data.Node, err error) {

	stations = make([]data.Node, 0)

	minWay := make([]*data.Node, 0)

	route = &data.NodeRoute{Route: minWay, TotalCost: 0}

	partialStart := start
	var partialTarget *data.Node
	var partialRoute *data.NodeRoute

	for partialTarget != target {

		wayCostEdges, err := CalcDijkstraToMany(g, partialStart)

		if err != nil {
			logrus.Fatalf("Error: %s", err.Error())
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

			//look for closest station to target(airdistance)

			var minStation int64 = -1
			var minDistance = math.MaxFloat64

			logrus.Debug("another station on the way is necessaray")

			for stationID := range gasStations.Stations {

				euclidDist := CalcEuclidDist(g.Nodes[stationID].Lat, g.Nodes[stationID].Lon, target.Lat, target.Lon)

				if wayCostEdges[stationID].Cost < rangeCm && euclidDist < minDistance {

					minStation = stationID
					minDistance = CalcEuclidDist(g.Nodes[stationID].Lat, g.Nodes[stationID].Lon, target.Lat, target.Lon)
				}

			}
			if minStation == -1 {
				logrus.Error("no good station found")

			}
			// add way to this station to route and station to list
			partialTarget = &g.Nodes[minStation]
			partialRoute, _ = findWayToGoal(partialStart, partialTarget, g, wayCostEdges)
			partialStart = &g.Nodes[minStation]
			//reverse route
			for i := len(partialRoute.Route) - 1; i >= 1; i-- {
				route.Route = append(route.Route, partialRoute.Route[i])
			}
			route.TotalCost += partialRoute.TotalCost
			stations = append(stations, *partialTarget)
		}

	}
	return route, stations, nil

}

func goalReachable(goal *data.Node, wayCostEdges []data.Edge, rangeCm int64) bool {

	logrus.Debug(wayCostEdges[goal.ID].Cost)

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
