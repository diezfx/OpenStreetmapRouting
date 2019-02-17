package dijkstra_test

import (
	"OpenStreetmapRouting/config"
	"OpenStreetmapRouting/data"
	"OpenStreetmapRouting/dijkstra"
	"testing"

	"github.com/sirupsen/logrus"
)

var nodes = []data.Node{{ID: 0, Lat: 1.0, Lon: 1.0},
	{ID: 1, Lat: 10.0, Lon: 1.0},
	{ID: 2, Lat: 1.0, Lon: 10.0},
	{ID: 3, Lat: 3.0, Lon: 3.0},
	{ID: 5, Lat: 5.0, Lon: 5.0},
	{ID: 5, Lat: 10101010, Lon: 1010},
	{ID: 6, Lat: 5.0, Lon: 5.0}}
var edges = []data.Edge{{Start: 0, End: 1, Cost: 9},
	{Start: 0, End: 2, Cost: 8},
	{Start: 0, End: 4, Cost: 7},
	{Start: 2, End: 0, Cost: 6},
	{Start: 2, End: 1, Cost: 5},
	{Start: 2, End: 4, Cost: 4},
	{Start: 3, End: 2, Cost: 3},
	{Start: 4, End: 1, Cost: 2},
	{Start: 4, End: 3, Cost: 1},
	{Start: 4, End: 5, Cost: 2},
	{Start: 5, End: 6, Cost: 1}}

var conf *config.Config

var graph *data.GraphProd

func Init(nodes []data.Node, edges []data.Edge) {

	conf = config.LoadConfig("../res/config_test.yaml")

	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)

	graphData := data.Graph{Nodes: nodes, Edges: edges}
	graph = data.InitGraphProd(&graphData, conf)

}

func TestDijkstraCostCalc(t *testing.T) {

	Init(nodes, edges)

	route, _ := dijkstra.CalcDijkstra(graph, &graph.Nodes[0], &graph.Nodes[2])

	if route.TotalCost != 8 {
		t.Errorf("Expected cost of %d got %d", 8, route.TotalCost)
	}

	if len(route.Route) != 2 {
		t.Errorf("Expected a way of length %d got %d", 2, len(route.Route))
	}

	route, _ = dijkstra.CalcDijkstra(graph, &graph.Nodes[0], &graph.Nodes[3])

	if route.TotalCost != 8 {
		t.Errorf("Expected cost of %d got %d", 8, route.TotalCost)
	}
	if len(route.Route) != 3 {
		t.Errorf("Expected a way of length %d got %d", 3, len(route.Route))
	}

	route, _ = dijkstra.CalcDijkstra(graph, &graph.Nodes[2], &graph.Nodes[6])

	if len(route.Route) != 4 {
		t.Errorf("Expected a way of length %d got %d", 4, len(route.Route))
	}

}
