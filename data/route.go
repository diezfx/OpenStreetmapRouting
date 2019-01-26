package data

// reduced data for drawing
type GetRoute struct {
	Route GeoJsonRoute

	TotalCost int64
}

type Coordinate struct {
	Lat float64
	Lon float64
}

type NodeRoute struct {
	Route     []*Node
	TotalCost int64
}

// convert and reverse
func (nr *NodeRoute) ConvertToJson() *GetRoute {

	getRoute := GetRoute{Route: GeoJsonRoute{Type: "LineString", Coordinates: make([][]float64, 0)}, TotalCost: nr.TotalCost}

	for i := len(nr.Route) - 1; i >= 0; i-- {

		getRoute.Route.Coordinates = append(getRoute.Route.Coordinates, []float64{nr.Route[i].Lon, nr.Route[i].Lat})

	}

	return &getRoute

}

func ConvertAreaToJson(route []*Edge, g *GraphProd) *GeoJsonArea {

	getArea := GeoJsonArea{Type: "MultiLineString", Coordinates: make([][][]float64, 0)}

	for _, edge := range route {

		miniRoute := [][]float64{[]float64{g.Nodes[edge.Start].Lon, g.Nodes[edge.Start].Lat},
			[]float64{g.Nodes[edge.End].Lon, g.Nodes[edge.End].Lat}}
		getArea.Coordinates = append(getArea.Coordinates, miniRoute)

	}

	return &getArea

}

type GeoJsonRoute struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type GeoJsonArea struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}
