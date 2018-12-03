package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"test/config"
	"test/data"
	"test/parsing"
	"time"
)

func main() {

	fmt.Println("Beginne")
	start := time.Now()
	config := config.LoadConfig("res/config.yaml")

	if config.OsmParse == 1 {
		parsing.Parse()
	}

	// load and init graph
	dat, err := ioutil.ReadFile(config.OutputFilename)

	if err != nil {
		log.Fatal(err.Error())
	}

	graphData := data.Graph{}

	err = graphData.Unmarshal(dat)

	if err != nil {
		log.Fatal(err.Error())
	}

	graph := data.GraphProd{Nodes: graphData.Nodes, Edges: graphData.Edges}

	graph.CalcOffsetList()

	fmt.Println("Ready!!")
	elapsed := time.Since(start)
	log.Printf("loading took %s", elapsed)
}
