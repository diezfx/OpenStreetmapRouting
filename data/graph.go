package data

import (
	"github.com/diezfx/OpenStreetmapRouting/config"

	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"github.com/golang-collections/collections/queue"
	"github.com/sirupsen/logrus"
)

// instance that is usable by dijkstra
//todo DI
var info *MetaInfo

// GraphRaw contains the node,edge lists and additionaly a map for the old/new id mapping
type GraphRaw struct {
	NodeIDs     map[int64]int64
	NodeIDMutex sync.Mutex

	EdgeMutex sync.Mutex
	Info      MetaInfo
	Nodes     []Node
	Edges     []Edge
}

// GraphProd Graph with the offset list to use in "production" with dijkstra
type GraphProd struct {
	Nodes  []Node
	Edges  []Edge
	Offset []int
	Grid   Grid
}

type MetaInfo struct {
	RoadTypes  map[string]int
	NodesTotal int
	EdgesTotal int
}

type RoadType struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type MetaInfoGet struct {
	RoadTypes  []RoadType `json:"roadTypes"`
	NodesTotal int        `json:"nodesTotal"`
	EdgesTotal int        `json:"edgesTotal"`
}

//UpdateIDs update the ids in the edges and calculate the cost
//  list after edges
// open channel and send them while still creating the raw node graph cooler
func (g *GraphRaw) UpdateIDs() *Graph {

	graph := Graph{Edges: g.Edges, Nodes: g.Nodes}

	// finded biggest connected component

	for i, edge := range g.Edges {
		edge.Start = int64(g.NodeIDs[edge.Start])
		edge.End = int64(g.NodeIDs[edge.End])
		edge.Cost = calcEdgeCost(&g.Nodes[edge.Start], &g.Nodes[edge.End], &edge)

		g.Edges[i] = edge
	}

	SortEdges(g.Edges)
	return &graph

}

// CalcConnectedComponent returns a list which give every node a component; returns the biggest component
// super non optimal
/*
func CalcConnectedComponent(g *GraphProd) []bool {

	//what is the nr of the biggest one?
	maxComponentNr := 0
	//how many nodes are in it?
	maxComponenCount := 0

	randomStarts := make([]int, 5)

	for i := 0; i < len(randomStarts); i++ {

		randomStarts[i] = rand.Intn(len(g.Nodes) - 1)
	}
	var visited []bool

	for i, start := range randomStarts {
		visited = make([]bool, len(g.Nodes))

		if visited[start] == false {
			count := bfs(g, visited, g.Nodes[start])

			if count >= maxComponenCount {

				maxComponentNr = i
				maxComponenCount = count
			}

		}
	}

	if maxComponenCount != maxComponentNr {
		visited = make([]bool, len(g.Nodes))
		bfs(g, visited, g.Nodes[randomStarts[maxComponentNr]])

	}

	logrus.Debugf("The maximum estimated number of connected nodes is %d", maxComponenCount)
	return visited
}
*/

// returns true if there is a connection to biggest component so far
func bfs(g *GraphProd, visited []bool, start Node) int {
	q := queue.New()

	q.Enqueue(start.ID)
	counter := 0

	for q.Len() > 0 {
		curr := g.Nodes[q.Dequeue().(int64)]

		edgeBegin := g.Offset[curr.ID]
		edgeEnd := g.Offset[curr.ID+1]

		for i := edgeBegin; i < edgeEnd; i++ {

			if visited[g.Edges[i].End] == false {
				q.Enqueue(g.Edges[i].End)
				visited[g.Edges[i].End] = true
				counter++

			}
		}
	}
	return counter

}

func (i *MetaInfo) WriteToFile(config *config.Config) {

	infoJSON, err := json.Marshal(i)

	info = i

	if err != nil {
		logrus.Fatal(err)
		return
	}

	ioutil.WriteFile(config.InfoFilename, infoJSON, 0644)

}

//ConverToGetObject conversion is needed for vue.js table
func (i *MetaInfo) ConverToGetObject() *MetaInfoGet {

	infoJSONGet := MetaInfoGet{RoadTypes: make([]RoadType, 0), NodesTotal: i.NodesTotal, EdgesTotal: i.EdgesTotal}

	for k, v := range i.RoadTypes {

		roadType := RoadType{Type: k, Count: v}
		infoJSONGet.RoadTypes = append(infoJSONGet.RoadTypes, roadType)
	}

	return &infoJSONGet

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

//CalcOffsetList calculates the offset list
func (g *GraphProd) CalcOffsetList() {

	currNodeID := 0

	g.Offset = make([]int, 0, len(g.Nodes))

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

//InitGraphProd add the offset list that is needed for dijkstra and the grid
func InitGraphProd(graphData *Graph, conf *config.Config) *GraphProd {

	g := &GraphProd{Nodes: graphData.Nodes, Edges: graphData.Edges}

	g.CalcOffsetList()

	grid := Grid{}
	grid.InitGrid(g.Nodes, conf)

	g.Grid = grid

	return g

}

// InitGraphProdWithStations add stations to nodeList
// connect station to main graph
// calculate offsetlist, sort edges again
func InitGraphProdWithStations(graphProd *GraphProd, conf *config.Config) *GasStations {

	//visited := CalcConnectedComponent(graphProd)

	//graphProd.Grid.connectedComponent = visited

	stations := GetFuelStations()

	// add stations to graph
	// update stations map with new ids

	stationsNew := make(map[int64]Node)

	for _, station := range stations.Stations {

		station.ID = int64(len(graphProd.Nodes))
		graphProd.Nodes = append(graphProd.Nodes, station)

		//connect station to closest mainroad noad
		connectNode := graphProd.Grid.FindNextNode(station.Lat, station.Lon, false)

		newEdge := Edge{ID: int64(len(graphProd.Edges)), Start: station.ID, End: connectNode.ID, Speed: 5, Cost: 10}
		newBackEdge := Edge{ID: int64(len(graphProd.Edges) + 1), Start: connectNode.ID, End: station.ID, Speed: 5, Cost: 10}

		graphProd.Edges = append(graphProd.Edges, newEdge, newBackEdge)
		stationsNew[station.ID] = station

		if station.ID > 10000000 {
			logrus.Debug(station.ID)
		}

	}
	stations.SetStations(stationsNew)

	// order is important
	// node order should be correct already
	SortEdges(graphProd.Edges)
	graphProd.CalcOffsetList()

	return stations
}

func GetGraphInfo() *MetaInfo {
	if info == nil {
		logrus.Errorf("Info not initialized")

		info = &MetaInfo{}
		info.LoadInfo(config.GetConfig())

	}
	return info
}

func (i *MetaInfo) LoadInfo(conf *config.Config) {

	dat, err := ioutil.ReadFile(conf.InfoFilename)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(dat, i)

	if err != nil {
		log.Fatal(err.Error())
	}

}
