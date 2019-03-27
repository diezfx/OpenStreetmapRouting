package data

import (
	"sort"

	"github.com/umahmood/haversine"
)

//AddEdge adds an edge to the graph
func (g *GraphRaw) AddEdge(e Edge) {
	g.EdgeMutex.Lock()
	g.Edges = append(g.Edges, e)
	g.EdgeMutex.Unlock()
}

//calcEdgeCost get distance in cm
func calcEdgeCost(start, end *Node, e *Edge) int64 {

	_, dist := haversine.Distance(haversine.Coord{Lat: start.Lat, Lon: start.Lon}, haversine.Coord{Lat: end.Lat, Lon: end.Lon})

	return int64(dist * 1000 * 100)
}

//SortEdges sorts the edges depending on their start node
func SortEdges(edges []Edge) {

	IDs := func(i, j int) bool {
		return edges[i].Start < edges[j].Start

	}
	sort.Slice(edges, IDs)

}
