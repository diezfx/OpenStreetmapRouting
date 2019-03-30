package data_test

import (
	"github.com/diezfx/OpenStreetmapRouting/config"
	"github.com/diezfx/OpenStreetmapRouting/data"
	"github.com/diezfx/OpenStreetmapRouting/parsing"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

var nodes = []data.Node{{ID: 0, Lat: 1.0, Lon: 1.0}, {ID: 1, Lat: 10.0, Lon: 1.0}, {ID: 2, Lat: 1.0, Lon: 10.0}, {ID: 3, Lat: 3.0, Lon: 3.0}, {ID: 4, Lat: 5.0, Lon: 5.0}, {ID: 5, Lat: -1.0, Lon: -1.0}}
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

func TestMain(m *testing.M) {
	conf = config.LoadConfig("../res/config_test.yaml")
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)

	os.Exit(m.Run())

}
func TestOffseGeneration(t *testing.T) {

	graph := data.GraphProd{Nodes: nodes, Edges: edges}

	graph.CalcOffsetList()

	expectedOffsetList := []int{0, 3, 3, 6, 7, 9}

	for i := 0; i < len(expectedOffsetList); i++ {

		if expectedOffsetList[i] != graph.Offset[i] {
			t.Errorf("Wrong offset, should be %d instead of %d", expectedOffsetList[i], graph.Offset[i])
		}

	}
}

// first creates small graph, saves it then reads it again
func TestGraphReading(t *testing.T) {

	graph := data.Graph{Nodes: nodes, Edges: edges}
	graph.WriteToFile(conf)

	graphData := parsing.ParseOrLoadGraph(conf)

	if cmp.Equal(graph, *graphData) == false {
		t.Errorf("Wrong graphData, should be %v+ instead of %v+", graph, *graphData)
	}
}
