package data

import (
	"sort"
	"sync"

	"github.com/umahmood/haversine"
)

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
}

//Convert update the ids in the edges and calculate the cost
// sort list after edges
// open channel and send them while still creating the raw node graph cooler
func (graphR *GraphRaw) Convert() *Graph {

	graph := Graph{Edges: graphR.Edges, Nodes: graphR.Nodes}

	for i, edge := range graphR.Edges {
		edge.Start = int64(graphR.NodeIDs[edge.Start])
		edge.End = int64(graphR.NodeIDs[edge.End])
		edge.Cost = calcEdgeCost(&graphR.Nodes[edge.Start], &graphR.Nodes[edge.End], &edge)

		graphR.Edges[i] = edge
	}

	sortEdges(graph.Edges)
	//graph.CalcOffsetList()
	return &graph

}

func (g *GraphRaw) AddEdge(e Edge) {
	g.EdgeMutex.Lock()
	g.Edges = append(g.Edges, e)
	g.EdgeMutex.Unlock()
}

//calcEdgeCost get distance then divide by the speed to get the cost for the edge
func calcEdgeCost(start, end *Node, e *Edge) int64 {

	_, dist := haversine.Distance(haversine.Coord{Lat: start.Lat, Lon: start.Lon}, haversine.Coord{Lat: end.Lat, Lon: end.Lon})

	return int64(dist / e.Speed)
}

func sortEdges(edges []Edge) {

	sortIDs := func(i, j int) bool {
		return edges[i].Start < edges[j].Start

	}
	sort.Slice(edges, sortIDs)

}

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
