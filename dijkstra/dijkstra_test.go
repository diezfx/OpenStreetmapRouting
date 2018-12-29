package dijkstra_test

import (
	"os"
	"test/config"
	"test/data"
	"test/dijkstra"
	"testing"

	"github.com/sirupsen/logrus"
)

var nodes = []data.Node{{ID: 0, Lat: 1.0, Lon: 1.0}, {ID: 1, Lat: 10.0, Lon: 1.0}, {ID: 2, Lat: 1.0, Lon: 10.0}, {ID: 3, Lat: 3.0, Lon: 3.0}, {ID: 4, Lat: 5.0, Lon: 5.0}}
var edges = []data.Edge{{Start: 0, End: 1, Cost: 9},
	{Start: 0, End: 2, Cost: 8},
	{Start: 0, End: 4, Cost: 7},
	{Start: 2, End: 0, Cost: 6},
	{Start: 2, End: 1, Cost: 5},
	{Start: 2, End: 4, Cost: 4},
	{Start: 3, End: 2, Cost: 3},
	{Start: 4, End: 1, Cost: 2},
	{Start: 4, End: 3, Cost: 1}}

var conf *config.Config

var graph *data.GraphProd

func TestMain(m *testing.M) {

	conf = config.LoadConfig("../res/config_test.yaml")

	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)

	graph = &data.GraphProd{Nodes: nodes, Edges: edges}
	graph.Init(conf)

	os.Exit(m.Run())

}

func TestDijkstraCostCalc(t *testing.T) {

	cost, optWay, _ := dijkstra.CalcDijkstra(graph, &graph.Nodes[0], &graph.Nodes[2])

	if cost != 8 {
		t.Errorf("Expected cost of %d got %d", 8, cost)
	}

	if len(optWay) != 2 {
		t.Errorf("Expected a way of length %d got %d", 2, len(optWay))
	}

	cost, optWay, _ = dijkstra.CalcDijkstra(graph, &graph.Nodes[0], &graph.Nodes[3])

	if cost != 8 {
		t.Errorf("Expected cost of %d got %d", 8, cost)
	}
	if len(optWay) != 3 {
		t.Errorf("Expected a way of length %d got %d", 2, len(optWay))
	}

}
