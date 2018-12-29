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
func parse() *data.Graph {
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

	DataHandler2 := DataHandlerStep2{Graph: DataHandler.Graph}

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

	//fmt.Println(len(graph.Offset))

	log.Infof("parsing took %s", time.Since(start))
	return graph

}

func ParseOrLoadGraph(config *config.Config) *data.Graph {

	var graphData *data.Graph
	if config.OsmParse == 1 {
		graphData = parse()
		graphData.WriteToFile(config)
	} else {
		// load and init graph
		dat, err := ioutil.ReadFile(config.OutputFilename)
		//todo try to parse if file not exists
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
