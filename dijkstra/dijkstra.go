package dijkstra

import (
	"OpenStreetmapRouting/data"
	"container/heap"
	"errors"
	"math"
	"time"

	"github.com/sirupsen/logrus"
)

func GetRoute(start data.Coordinate, end data.Coordinate) (*data.NodeRoute, error) {

	graph := data.GetGraphProd()

	logrus.Debug(start, end)

	logrus.Infof("Find nodes close to Node")
	startTime := time.Now()
	startNode := graph.Grid.FindNextNode(start.Lat, start.Lon, false)
	endNode := graph.Grid.FindNextNode(end.Lat, end.Lon, false)

	gridTime := time.Since(startTime)

	logrus.Infof("Dijkstra started")

	result, err := CalcDijkstra(graph, startNode, endNode)
	dijkstraTime := time.Since(startTime) - gridTime
	endTime := time.Since(startTime)
	logrus.WithFields(logrus.Fields{
		"Time for Gridsearch": gridTime,
		"Time for dijkstra":   dijkstraTime,
		"Total time":          endTime}).Info("Dijkstra ended")

	return result, err

}

// CalcDijkstra takes a starting node and returns all edges on the way
// uses edges for the overview of cost and the way to the previous node
func CalcDijkstra(g *data.GraphProd, start *data.Node, target *data.Node) (*data.NodeRoute, error) {

	pq := make(data.PriorityQueue, 0, 10)

	if start.ID == target.ID {
		route := &data.NodeRoute{Route: make([]*data.Node, 0), TotalCost: 0}
		return route, nil
	}

	//sets the edge that led to this node
	prevs := make([]data.Edge, len(g.Nodes))

	for i := range prevs {

		edge := data.Edge{ID: -1, End: start.ID, Start: start.ID, Cost: math.MaxInt64}
		prevs[i] = edge

	}

	heap.Init(&pq)
	//edge for the begining
	edge := data.Edge{ID: -1, End: start.ID, Start: start.ID, Cost: 0}
	heap.Push(&pq, &data.Item{Value: edge, Priority: 0})

	for pq.Len() > 0 {

		item := heap.Pop(&pq).(*data.Item)

		currentEdge := item.Value.(data.Edge)

		if item.Priority >= prevs[currentEdge.End].Cost {

			continue

		}

		currentEdge.Cost = item.Priority
		prevs[currentEdge.End] = currentEdge
		// look at all reachable nodes
		edgeBegin := g.Offset[currentEdge.End]
		edgeEnd := g.Offset[currentEdge.End+1]

		for i := edgeBegin; i < edgeEnd; i++ {

			newItem := data.Item{Value: g.Edges[i], Priority: item.Priority + g.Edges[i].Cost}

			// skip if cost is bigger then what we already know
			if newItem.Priority < prevs[g.Edges[i].End].Cost {
				heap.Push(&pq, &newItem)

			}
		}
	}

	// add all nodes that are on the optimal way

	optWay := make([]*data.Node, 0)
	edge = prevs[target.ID]
	minCost := edge.Cost
	if edge.Cost == math.MaxInt64 {
		return nil, errors.New("no way found")
	}

	optWay = append(optWay, &g.Nodes[edge.End])

	for edge.End != start.ID {

		optWay = append(optWay, &g.Nodes[edge.Start])
		edge = prevs[edge.Start]

	}

	route := data.NodeRoute{Route: optWay, TotalCost: minCost}

	return &route, nil
}

// CalcDijkstra takes a starting node and returns all edges on the way
// uses edges for the overview of cost and the way to the previous node
func CalcDijkstraToMany(g *data.GraphProd, start *data.Node) ([]data.Edge, error) {

	pq := make(data.PriorityQueue, 0, 10)

	//sets the edge that led to this node
	prevs := make([]data.Edge, len(g.Nodes))

	for i := range prevs {

		edge := data.Edge{ID: -1, End: start.ID, Start: start.ID, Cost: math.MaxInt64}
		prevs[i] = edge

	}

	heap.Init(&pq)
	//edge for the begining
	edge := data.Edge{ID: -1, End: start.ID, Start: start.ID, Cost: 0}
	heap.Push(&pq, &data.Item{Value: edge, Priority: 0})

	for pq.Len() > 0 {

		item := heap.Pop(&pq).(*data.Item)

		currentEdge := item.Value.(data.Edge)

		if item.Priority >= prevs[currentEdge.End].Cost {

			continue

		}

		currentEdge.Cost = item.Priority
		prevs[currentEdge.End] = currentEdge
		// look at all reachable nodes
		edgeBegin := g.Offset[currentEdge.End]
		edgeEnd := g.Offset[currentEdge.End+1]

		for i := edgeBegin; i < edgeEnd; i++ {

			newItem := data.Item{Value: g.Edges[i], Priority: item.Priority + g.Edges[i].Cost}

			// skip if cost is bigger then what we already know
			if newItem.Priority < prevs[g.Edges[i].End].Cost {
				heap.Push(&pq, &newItem)

			}
		}
	}

	// add all nodes that are on the optimal way

	return prevs, nil
}
