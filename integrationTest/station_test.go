// test if every gas station is reachable in the "big network"
package integrationTest

import (
	"github.com/diezfx/OpenStreetmapRouting/config"
	"github.com/diezfx/OpenStreetmapRouting/data"
	"github.com/diezfx/OpenStreetmapRouting/dijkstra"
	"github.com/diezfx/OpenStreetmapRouting/parsing"
	"math"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

var conf *config.Config
var graph *data.GraphProd

func TestMain(m *testing.M) {
	conf = config.LoadConfig("../res/config_test.yaml")
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)

	graph = InitGraphProd()
	data.InitGraphProdWithStations(graph, conf)

	os.Exit(m.Run())

}

func TestStationsReachable(t *testing.T) {

	stations := data.GetFuelStations()

	errorCount := 0
	startNode := graph.Grid.FindNextNode(48.739889600673365, 9.105295872250478, false)
	goalCosts, _ := dijkstra.CalcDijkstraToMany(graph, *startNode)

	for _, station := range stations.Stations {

		goalNode := graph.Grid.FindNextNode(station.Lat, station.Lon, false)

		if goalCosts[goalNode.ID].Cost >= math.MaxInt64 {
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("Expected all stations reachable got %d errors", errorCount)

	}

}

//Init calculates the offsetlist and creates the grid for the given graph
func InitGraphProd() *data.GraphProd {

	conf := config.GetConfig()
	graphData := parsing.ParseOrLoadGraph(conf)

	g := data.InitGraphProd(graphData, conf)

	return g

}
