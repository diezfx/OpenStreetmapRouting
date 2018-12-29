package data

import (
	"OpenStreetmapRouting/config"
	"math"

	"github.com/sirupsen/logrus"
)

//Grid contains a grid that helps finding the next node to a Lat long input
type Grid struct {
	Grid                             map[int][]*Node
	latMin, latMax, longMin, longMax float64

	LatSize  int
	LongSize int
}

//InitGrid every node is added to a cell
func (g *Grid) InitGrid(graph *GraphProd, config *config.Config) {

	g.Grid = make(map[int][]*Node, 0)

	g.LatSize = config.GridXSize
	g.LongSize = config.GridYSize

	g.latMin, g.latMax = findMinMaxLat(graph)

	g.longMin, g.longMax = findMinMaxLong(graph)

	//add all nodes
	//init list if doesnt exist
	for i, node := range graph.Nodes {
		x, y := g.CalculateGridPos(node.Lat, node.Lon)
		if list, ok := g.Grid[x*g.LatSize+y]; ok == true {

			list = append(list, &graph.Nodes[i])

		} else {
			list := make([]*Node, 0)

			list = append(list, &graph.Nodes[i])
			g.Grid[x*g.LatSize+y] = list

		}
	}

}

func findMinMaxLat(graph *GraphProd) (min, max float64) {

	max = -1.0
	min = math.MaxFloat64

	for _, node := range graph.Nodes {

		if node.Lat > max {
			max = node.Lat
		}
		if node.Lat < min {
			min = node.Lat
		}
	}

	return

}

func findMinMaxLong(graph *GraphProd) (min, max float64) {

	max = -1.0
	min = math.MaxFloat64

	for _, node := range graph.Nodes {

		if node.Lon > max {
			max = node.Lon
		}
		if node.Lon < min {
			min = node.Lon
		}
	}

	return
}

func (g *Grid) CalculateGridPos(lat, long float64) (x, y int) {
	x = -1
	y = -1

	deltaLat := g.latMax - g.latMin
	deltaLong := g.longMax - g.latMin

	//subtract the min to center then divide by delta
	latRelative := (lat - g.latMin) / deltaLat
	longRelative := (long - g.longMin) / deltaLong

	//
	x = int(latRelative * float64(g.LatSize))
	y = int(longRelative * float64(g.LongSize))

	return
}

//FindNextNode searches for the closest node for the given point
func (g *Grid) FindNextNode(lat, long float64) *Node {

	//first try the gridCell it is in
	x, y := g.CalculateGridPos(lat, long)

	candidates := make([]*Node, 0)
	if list, ok := g.Grid[x*g.LatSize+y]; ok == true {
		candidates = append(candidates, list...)
	}

	for dist := 0; true; dist++ {
		for i := 0; i <= dist; i++ {

			//check if grid can exist
			j := dist - i

			if x+i < g.LatSize && y+j < g.LongSize {
				if list, ok := g.Grid[(x+i)*g.LatSize+(y+j)]; ok == true {
					candidates = append(candidates, list...)
				}
			}
			if j == 0 && i == 0 {
				continue
			}
			if x+i < g.LatSize && y-j >= 0 && j != 0 {
				if list, ok := g.Grid[(x+i)*g.LatSize+(y-j)]; ok == true {
					candidates = append(candidates, list...)
				}
			}
			if x-i >= 0 && y+j < g.LongSize && x != 0 {
				if list, ok := g.Grid[(x-i)*g.LatSize+(y+j)]; ok == true {
					candidates = append(candidates, list...)
				}
			}
			if x-i >= 0 && y-j >= 0 && j != 0 && x != 0 {
				if list, ok := g.Grid[(x-i)*g.LatSize+(y-j)]; ok == true {
					candidates = append(candidates, list...)
				}
			}

		}

		//a candidate exists?
		//else add more cells
		if dist%2 == 0 {
			if len(candidates) > 0 {
				logrus.Debug(candidates)
				return findClosestNode(candidates, lat, long)

			}
		}
	}

	//should never be reached
	return nil

}

// find closest node
func findClosestNode(nodes []*Node, targetLat, targetLon float64) *Node {

	minDist := math.MaxFloat64
	pos := -1

	for i, node := range nodes {
		dist := CalcEuclidDist(math.Abs(targetLat-node.Lat), math.Abs(targetLon-node.Lon))
		if dist < minDist {
			pos = i
			minDist = dist
		}
	}

	return nodes[pos]

}

func CalcEuclidDist(x, y float64) float64 {

	return math.Sqrt(math.Pow(x+y, 2))
}
