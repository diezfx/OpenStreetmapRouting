package dijkstra_test

import (
	"OpenStreetmapRouting/data"
	"OpenStreetmapRouting/dijkstra"
	"testing"
)

var nodesSD = []data.Node{{ID: 0, Lat: 1.0, Lon: 1.0},
	{ID: 1, Lat: 1, Lon: 1.0},
	{ID: 2, Lat: 1.0, Lon: 10.0},
	{ID: 3, Lat: 3.0, Lon: 3.0},
	{ID: 4, Lat: 5.0, Lon: 5.0},
	{ID: 5, Lat: 10101010, Lon: 1010},
	{ID: 6, Lat: 5.0, Lon: 5.0},
	{ID: 7, Lat: 1000.0, Lon: 5.0}}
var edgesSD = []data.Edge{{Start: 0, End: 1, Cost: 2},
	{Start: 1, End: 2, Cost: 2},
	{Start: 2, End: 3, Cost: 2},
	{Start: 2, End: 7, Cost: 1},
	{Start: 3, End: 4, Cost: 2},
	{Start: 3, End: 6, Cost: 1},
	{Start: 4, End: 5, Cost: 2},
	{Start: 6, End: 3, Cost: 1},
	{Start: 7, End: 2, Cost: 1}}

var stationNodes = map[int64]data.Node{6: data.Node{ID: 6, Lat: 2.0, Lon: 1.0}, 7: data.Node{ID: 7, Lat: 10.0, Lon: 1.0}}

var stations = data.GasStations{Stations: stationNodes}

func TestStationDijkstra(t *testing.T) {

	Init(nodesSD, edgesSD)

	route, visitedStations, err := dijkstra.CalcStationDijkstraSnapshots(graph, &stations, graph.Nodes[0], graph.Nodes[5], 8.0)

	if route.TotalCost != 12 {
		t.Errorf("Expected cost of %d got %d", 12, route.TotalCost)
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(route.Route) != 8 {
		t.Errorf("Expected a way of length %d got %d", 8, len(route.Route))
	}

	if visitedStations[0].ID != stations.Stations[7].ID {
		t.Errorf("Expected station %d got %d", stations.Stations[7].ID, visitedStations[0].ID)
	}

	// test direct route
	route, visitedStations, err = dijkstra.CalcStationDijkstra(graph, &stations, graph.Nodes[0], graph.Nodes[5], 16)

	if route.TotalCost != 10 {
		t.Errorf("Expected cost of %d got %d", 12, route.TotalCost)
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(route.Route) != 6 {
		t.Errorf("Expected a way of length %d got %d", 6, len(route.Route))
	}

	if len(visitedStations) != 0 {
		t.Errorf("Expected #station %d got %d", 0, len(visitedStations))
	}

}
