package parsing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"test/config"
	"test/data"
	"time"

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

func Parse() {
	fmt.Println("Beginne")
	start := time.Now()

	config := config.GetConfig()

	fmt.Println(config)

	DataHandler := DataHandlerStep1{}
	DataHandler.InitGraph()

	// read graph data
	dec := loadDec()

	fmt.Println("Saving all edges")
	err := dec.Parse(&DataHandler)
	if err != nil {
		panic(err)
	}

	DataHandler2 := DataHandlerStep2{Graph: DataHandler.Graph}

	fmt.Println("Saving all nodes")
	dec = loadDec()
	err = dec.Parse(&DataHandler2)
	if err != nil {
		panic(err)
	}

	/////////////////////////////////////
	// Converting parsed graph
	////////////////////////////////

	fmt.Println("Converting Graph to a more efficinet one")
	//graph := DataHandler.Graph.Convert()

	fmt.Println(len(DataHandler2.Graph.Nodes))
	fmt.Println(len(DataHandler2.Graph.Edges))
	//fmt.Println(len(graph.Edges))
	//fmt.Println(len(graph.Nodes))

	fmt.Println("Finished parsing.")

	// transform to efficient graph
	graph := DataHandler2.Graph.Convert()

	//fmt.Println(len(graph.Offset))

	fmt.Println(graph.Edges[5])

	elapsed := time.Since(start)
	log.Printf("parsing took %s", elapsed)

	writeToFile(graph)

}

func writeToFile(graph *data.Graph) {

	config := config.GetConfig()

	var encodedGraph []byte

	if config.OutputType == "json" {
		jsonGraph, err := json.Marshal(graph)
		encodedGraph = jsonGraph
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	} else {
		protoGraph, err := graph.Marshal()
		encodedGraph = protoGraph
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}

	ioutil.WriteFile(config.OutputFilename, encodedGraph, 0644)

}
