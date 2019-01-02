package data

import (
	"OpenStreetmapRouting/config"
	"encoding/json"
	"io/ioutil"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/umahmood/haversine"
)

// instance that is usable by dijkstra
var graphProd *GraphProd

// GraphRaw contains the node,edge lists and additionaly a map for the old/new id mapping
type GraphRaw struct {
	NodeIDs     map[int64]int64
	NodeIDMutex sync.Mutex

	EdgeMutex sync.Mutex

	Nodes []Node
	Edges []Edge
}

// GraphProd Graph with the offset list to use in "production" with dijkstra
type GraphProd struct {
	Nodes  []Node
	Edges  []Edge
	Offset []int
	Grid   Grid
}

//UpdateIDs update the ids in the edges and calculate the cost
// sort list after edges
// open channel and send them while still creating the raw node graph cooler
func (g *GraphRaw) UpdateIDs() *Graph {

	graph := Graph{Edges: g.Edges, Nodes: g.Nodes}

	for i, edge := range g.Edges {
		edge.Start = int64(g.NodeIDs[edge.Start])
		edge.End = int64(g.NodeIDs[edge.End])
		edge.Cost = calcEdgeCost(&g.Nodes[edge.Start], &g.Nodes[edge.End], &edge)

		g.Edges[i] = edge
	}

	sortEdges(g.Edges)
	return &graph

}

//WriteToFile write the graph to a file depending on the config json|protobuf
func (g *Graph) WriteToFile(config *config.Config) {

	var encodedGraph []byte

	if config.OutputType == "json" {
		jsonGraph, err := json.Marshal(g)
		encodedGraph = jsonGraph
		if err != nil {
			logrus.Fatal(err)
			return
		}

	} else {
		protoGraph, err := g.Marshal()
		encodedGraph = protoGraph
		if err != nil {
			logrus.Fatal(err)
		}

	}

	ioutil.WriteFile(config.OutputFilename, encodedGraph, 0644)

}

//AddEdge adds an edge to the graph
func (g *GraphRaw) AddEdge(e Edge) {
	g.EdgeMutex.Lock()
	g.Edges = append(g.Edges, e)
	g.EdgeMutex.Unlock()
}

//calcEdgeCost get distance then divide by the speed to get the cost for the edge
func calcEdgeCost(start, end *Node, e *Edge) int64 {

	_, dist := haversine.Distance(haversine.Coord{Lat: start.Lat, Lon: start.Lon}, haversine.Coord{Lat: end.Lat, Lon: end.Lon})

	return int64((dist * 1000 / (e.Speed / 10000)))
}

func sortEdges(edges []Edge) {

	sortIDs := func(i, j int) bool {
		return edges[i].Start < edges[j].Start

	}
	sort.Slice(edges, sortIDs)

}

//CalcOffsetList calculates the offset list
func (g *GraphProd) CalcOffsetList() {

	currNodeID := 0

	g.Offset = append(g.Offset, 0)

	for i := 0; i < len(g.Edges); i++ {

		edge := g.Edges[i]

		if int64(currNodeID) != edge.Start {

			// check if some nodes have no outgoing edges
			for j := currNodeID; int64(j) < edge.Start; j++ {
				g.Offset = append(g.Offset, i)
				currNodeID++
			}

		}

	}
	for j := currNodeID; j < len(g.Nodes); j++ {
		g.Offset = append(g.Offset, len(g.Edges))

	}

}

// //add the offset list that is needed for dijkstra and the grid
func InitGraphProd(graphData *Graph, conf *config.Config) *GraphProd {

	g := &GraphProd{Nodes: graphData.Nodes, Edges: graphData.Edges}

	g.CalcOffsetList()

	grid := Grid{}
	grid.InitGrid(g, conf)

	g.Grid = grid

	graphProd = g

	return g

}

func GetGraphProd() *GraphProd {
	if graphProd != nil {
		return graphProd
	} else {
		logrus.Fatal("Graph is not initialized")
	}
	return nil
}
