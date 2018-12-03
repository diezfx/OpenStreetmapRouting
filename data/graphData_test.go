package data_test

import (
	"test/data"
	"testing"
)

func TestOffseGeneration(t *testing.T) {

	nodes := []data.Node{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 3}}

	edges := []data.Edge{{Start: 0, End: 1, Cost: 9},
		{Start: 0, End: 2, Cost: 8},
		{Start: 0, End: 4, Cost: 7},
		{Start: 2, End: 0, Cost: 6},
		{Start: 2, End: 1, Cost: 5},
		{Start: 2, End: 4, Cost: 4},
		{Start: 3, End: 2, Cost: 3},
		{Start: 4, End: 1, Cost: 2},
		{Start: 4, End: 3, Cost: 1}}

	graph := data.Graph{Nodes: nodes, Edges: edges}

	graph.CalcOffsetList()

	expectedOffsetList := []int{0, 3, 3, 6, 7, 9}

	for i := 0; i < len(expectedOffsetList); i++ {

		if expectedOffsetList[i] != graph.Offset[i] {
			t.Errorf("Wrong offset, should be %d instead of %d", expectedOffsetList[i], graph.Offset[i])
		}

	}
}
