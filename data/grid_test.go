package data_test

import (
	"test/data"
	"testing"
)

func TestGridGeneration(t *testing.T) {

	grid := data.Grid{}
	graph := &data.GraphProd{Nodes: nodes, Edges: edges}
	grid.InitGrid(graph, conf)

	node := grid.FindNextNode(1, 1)
	if node.ID != 0 {
		t.Errorf("Expected NodeId %d got %d", 0, node.ID)
	}

	node = grid.FindNextNode(10, 1)

	if node.ID != 1 {
		t.Errorf("Expected NodeId %d got %d", 1, node.ID)
	}

	node = grid.FindNextNode(4, 4.1)

	if node.ID != 4 {
		t.Errorf("Expected NodeId %d got %d", 4, node.ID)
	}

}
