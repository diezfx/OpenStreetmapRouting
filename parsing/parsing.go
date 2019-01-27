package parsing

import (
	"OpenStreetmapRouting/config"
	"OpenStreetmapRouting/data"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thomersch/gosmparse"
)

func loadDec() *gosmparse.Decoder {

	config := config.GetConfig()

	var dec *gosmparse.Decoder
	if config.OsmLocation == "internet" {
		resp, err := http.Get(config.OsmFilename)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		dec = gosmparse.NewDecoder(resp.Body)
	}
	if config.OsmLocation == "local" {
		file, err := os.Open(config.OsmFilename)
		if err != nil {
			panic(err)
		}
		dec = gosmparse.NewDecoder(file)
	}
	return dec
}

//Parse parses the graph to our graph format
func parse() (*data.Graph, *data.MetaInfo) {
	log.Infof("Beginne")
	start := time.Now()

	config := config.GetConfig()

	log.Info(config)

	DataHandler := DataHandlerStep1{}
	DataHandler.InitGraph()

	// read graph data
	dec := loadDec()

	log.Info("Saving all edges")
	err := dec.Parse(&DataHandler)
	if err != nil {
		panic(err)
	}

	DataHandler2 := DataHandlerStep2{Graph: DataHandler.Graph, GasStationList: DataHandler.GasStationList}

	log.Info("Saving all nodes")
	dec = loadDec()
	err = dec.Parse(&DataHandler2)
	if err != nil {
		panic(err)
	}

	log.Info("Finished parsing.")
	log.WithFields(log.Fields{
		"Node count": len(DataHandler2.Graph.Nodes),
		"Edge count": len(DataHandler2.Graph.Edges)}).Info("Graph infos")

	/////////////////////////////////////
	// Converting parsed graph to better nodeIds starting from 0
	////////////////////////////////

	log.Info("Start converting graph")

	// transform to efficient graph
	graph := DataHandler2.Graph.UpdateIDs()

	info := &DataHandler2.Graph.Info
	info.EdgesTotal = len(DataHandler2.Graph.Edges)
	info.NodesTotal = len(DataHandler2.Graph.Nodes)

	DataHandler2.GasStationList.WriteFile(config)

	//fmt.Println(len(graph.Offset))

	log.Infof("parsing took %s", time.Since(start))
	return graph, info

}

func ParseOrLoadGraph(config *config.Config) *data.Graph {

	var graphData *data.Graph
	var info *data.MetaInfo
	if config.OsmParse == 1 {
		graphData, info = parse()
		graphData.WriteToFile(config)
		info.WriteToFile(config)

	} else {
		// load and init graph
		dat, err := ioutil.ReadFile(config.OutputFilename)
		if err != nil {
			log.Fatal(err.Error())
		}
		graphData = &data.Graph{}
		err = graphData.Unmarshal(dat)

		if err != nil {
			log.Fatal(err.Error())
		}

	}

	return graphData
}
