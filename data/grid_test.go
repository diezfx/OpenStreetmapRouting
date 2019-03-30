package data_test

import (
	"github.com/diezfx/OpenStreetmapRouting/data"
	"testing"
)

func TestGridGeneration(t *testing.T) {

	grid := data.Grid{}
	graph := &data.GraphProd{Nodes: nodes, Edges: edges}
	grid.InitGrid(graph.Nodes, conf)

	node := grid.FindNextNode(1, 1, false)
	if node.ID != 0 {
		t.Errorf("Expected NodeId %d got %d", 0, node.ID)
	}

	node = grid.FindNextNode(10, 1, false)

	if node.ID != 1 {
		t.Errorf("Expected NodeId %d got %d", 1, node.ID)
	}

	node = grid.FindNextNode(4, 4.1, false)

	if node.ID != 4 {
		t.Errorf("Expected NodeId %d got %d", 4, node.ID)
	}

}

func TestGetArea(t *testing.T) {

	grid := data.Grid{}
	graph := &data.GraphProd{Nodes: nodes, Edges: edges}
	grid.InitGrid(graph.Nodes, conf)

	// area to get
	ne := data.Coordinate{Lat: 5.0, Lon: 1.0}
	sw := data.Coordinate{Lat: 10.0, Lon: 10.0}

	result := grid.GetNodesInArea(ne, sw)

	if len(result) != 2 {
		t.Errorf("Expected result size of 2 got %d", len(result))
	}

}
