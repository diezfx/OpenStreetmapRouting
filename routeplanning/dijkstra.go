package dijkstra

import (
	"container/heap"
	"math"
	"test/data"
)

func calcDijkstra(g *data.GraphProd, start *data.Node) {

	pq := data.PriorityQueue{}

	//sets the edge that led to thid node
	prevs := make([]*data.Edge, len(g.Nodes))

	tentCost := make([]int64, len(g.Nodes))

	for i := range tentCost {
		tentCost[i] = math.MaxInt64
	}

	heap.Init(&pq)

	heap.Push(&pq, data.Item{Value: start.ID, Priority: 0})

	tentCost[start.ID] = 0

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*data.Item)
		if item.Priority >= prevs[g.Edges[item.Value].End].Cost {
			continue
		}
		// look at all reachable nodes
		edgeBegin := g.Offset[start.ID]
		edgeEnd := g.Offset[start.ID+1]
		for i := edgeBegin; i < edgeEnd; i++ {

			newItem := data.Item{Value: int64(i), Priority: item.Priority + g.Edges[i].Cost}

			// skip if cost is bigger then what we already know
			if newItem.Priority < tentCost[g.Edges[i].End] {
				heap.Push(&pq, newItem)
			}
		}
	}
}
