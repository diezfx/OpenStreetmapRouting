package data

import (
	"github.com/diezfx/OpenStreetmapRouting/config"
	"math"
)

//Grid contains a grid that helps finding the next node to a Lat long input
type Grid struct {
	Grid                             map[int][]*Node
	latMin, latMax, longMin, longMax float64

	LatSize  int
	LongSize int
}

//InitGrid every node is added to a cell
func (g *Grid) InitGrid(nodeList []Node, config *config.Config) {

	g.Grid = make(map[int][]*Node, 0)

	g.LatSize = config.GridXSize
	g.LongSize = config.GridYSize

	g.latMin, g.latMax = findMinMaxLat(nodeList)

	g.longMin, g.longMax = findMinMaxLong(nodeList)

	//add all nodes
	//init list if doesnt exist
	for i, node := range nodeList {
		x, y := g.CalculateGridPos(node.Lat, node.Lon)
		if list, ok := g.Grid[x*g.LatSize+y]; ok == true {

			g.Grid[x*g.LatSize+y] = append(list, &nodeList[i])

		} else {
			list := make([]*Node, 0)

			list = append(list, &nodeList[i])
			g.Grid[x*g.LatSize+y] = list

		}
	}

}

func findMinMaxLat(nodeList []Node) (min, max float64) {

	max = -1.0
	min = math.MaxFloat64

	for _, node := range nodeList {

		if node.Lat > max {
			max = node.Lat
		}
		if node.Lat < min {
			min = node.Lat
		}
	}

	return

}

func findMinMaxLong(nodeList []Node) (min, max float64) {

	max = -1.0
	min = math.MaxFloat64

	for _, node := range nodeList {

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

	deltaLat := g.latMax - g.latMin
	deltaLong := g.longMax - g.longMin

	//subtract the min to center then divide by delta
	latRelative := (lat - g.latMin) / deltaLat
	longRelative := (long - g.longMin) / deltaLong

	//
	x = int(math.Round(latRelative * float64(g.LatSize)))
	y = int(math.Round(longRelative * float64(g.LongSize)))

	return
}

// idea for sending only data from an area
// send all get gridpos north-east and south-west and send all nodes in this grid

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *Grid) GetNodesInArea(northEast, southWest Coordinate) []*Node {

	xNE, yNE := g.CalculateGridPos(northEast.Lat, northEast.Lon)
	xSW, ySW := g.CalculateGridPos(southWest.Lat, southWest.Lon)

	xStart := xNE
	xEnd := xSW
	if xNE > xSW {
		xStart = xSW
		xEnd = xNE

	}
	yStart := yNE
	yEnd := ySW
	if yNE > ySW {
		yStart = ySW
		yEnd = yNE
	}

	// get all grids between the rectangle spanned from ne,sw
	// iterate over all grids
	// at least one grid looked at

	nodes := make([]*Node, 0)
	for i := xStart; i <= xEnd; i++ {

		for j := yStart; j <= yEnd; j++ {

			nodes = append(nodes, g.Grid[g.convert2DTo1D(i, j)]...)

		}
	}

	return nodes

}

func (g *Grid) convert2DTo1D(x, y int) int {

	return x*g.LatSize + y

}

//FindNextNode searches for the closest node for the given point
func (g *Grid) FindNextNode(lat, long float64, mainComponent bool) *Node {

	//first try the gridCell it is in
	x, y := g.CalculateGridPos(lat, long)

	candidates := make([]*Node, 0)
	if list, ok := g.Grid[x*g.LatSize+y]; ok == true {

		candidates = append(candidates, list...)
	}

	//don't waste too much time out of bounds
	// not completely correct
	if g.latMin > lat {
		lat = g.latMin
	}
	if g.latMax < lat {
		lat = g.latMax
	}

	if g.longMin > long {
		long = g.longMin
	}
	if g.longMax < long {
		long = g.longMax
	}

	//start := time.Now()
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
		/*
			if mainComponent {

				newCandidates := make([]*Node, 0)

				//filter candidates if only maincomp is allowed too lazy to do it at top
				for _, candidate := range candidates {
					if g.connectedComponent[candidate.ID] {
						newCandidates = append(newCandidates, candidate)

					}
				}
				candidates = newCandidates

			}
		*/

		if len(candidates) > 0 {

			return findClosestNode(candidates, lat, long)
			//elapsed := time.Since(start)
			//logrus.Debugf("Elapsed time %s with distance: %d", elapsed, dist)

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
		dist := CalcEuclidDist(targetLat, targetLon, node.Lat, node.Lon)
		if dist < minDist {
			pos = i
			minDist = dist
		}
	}

	return nodes[pos]

}

func CalcEuclidDist(x1, x2, y1, y2 float64) float64 {

	d1 := math.Abs(x1 - y1)
	d2 := math.Abs(x2 - y2)
	return math.Sqrt(math.Pow(d1, 2) + math.Pow(d2, 2))
}
