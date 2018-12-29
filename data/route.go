package data

// reduced data for drawing
type GetRoute struct {
	Route []Coordinate

	TotalCost int64
}

type Coordinate struct {
	Lat float64
	Lon float64
}

type NodeRoute struct {
	Route []*Node

	TotalCost int64
}
