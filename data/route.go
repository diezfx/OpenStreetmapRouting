package data

// GetRoute reduced data for drawing
type GetRoute struct {
	Route GeoJSONRoute

	TotalCost int64
}

// GetRouteWithStations .
type GetRouteWithStations struct {
	Route    GeoJSONRoute
	Stations []Node

	TotalCost int64
}

// Coordinate lat long coordinate
type Coordinate struct {
	Lat float64
	Lon float64
}

// NodeRoute route
type NodeRoute struct {
	Route     []*Node
	TotalCost int64
}

// ConvertToJSON takes the nodeRoute and returns a geojson linestring
func (nr *NodeRoute) ConvertToJSON() *GetRoute {

	getRoute := GetRoute{Route: GeoJSONRoute{Type: "LineString", Coordinates: make([][]float64, 0)}, TotalCost: nr.TotalCost}

	for i := range nr.Route {

		getRoute.Route.Coordinates = append(getRoute.Route.Coordinates, []float64{nr.Route[i].Lon, nr.Route[i].Lat})

	}

	return &getRoute

}

// ConvertAreaToJSON takes all route in an area and converts them to multilinestring geojson format
func ConvertAreaToJSON(route []*Edge, g *GraphProd) *GeoJSONArea {

	getArea := GeoJSONArea{Type: "MultiLineString", Coordinates: make([][][]float64, 0)}

	for _, edge := range route {

		miniRoute := [][]float64{[]float64{g.Nodes[edge.Start].Lon, g.Nodes[edge.Start].Lat},
			[]float64{g.Nodes[edge.End].Lon, g.Nodes[edge.End].Lat}}
		getArea.Coordinates = append(getArea.Coordinates, miniRoute)

	}

	return &getArea

}

// GeoJSONRoute route in geojson format
type GeoJSONRoute struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// GeoJSONArea all route in area in geojson format
type GeoJSONArea struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}
