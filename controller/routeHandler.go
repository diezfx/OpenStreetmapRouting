package controller

import (
	"OpenStreetmapRouting/data"
	"OpenStreetmapRouting/dijkstra"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// RouteAreaHandler take an area as argument and return all edges in this area
func (s *Server) RouteAreaHandler(w http.ResponseWriter, r *http.Request) {

	ne, sw, err := GetArea(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	nodeList := s.graph.Grid.GetNodesInArea(ne, sw)

	edgeList := make([]*data.Edge, 0)

	//get all edges that reach come from these nodes
	for _, node := range nodeList {
		nodeID := node.ID
		edgeBegin := s.graph.Offset[nodeID]
		edgeEnd := s.graph.Offset[nodeID+1]

		for i := edgeBegin; i < edgeEnd; i++ {

			edgeList = append(edgeList, &s.graph.Edges[i])

		}
	}

	areaGet := data.ConvertAreaToJSON(edgeList, s.graph)

	jsonData, err := json.Marshal(areaGet)

	if err != nil || jsonData == nil {

		logrus.Error(err)
		http.Error(w, "no area available", 500)
	}

	w.Write(jsonData)

}

func RouteHandler(w http.ResponseWriter, r *http.Request) {

	vars := r.URL.Query()

	startLatStr, ok := vars["startlat"]
	startLonStr, ok2 := vars["startlon"]

	if ok == false || ok2 == false {
		logrus.Warnf("invalid input %s,%s", startLatStr, startLonStr)
		http.Error(w, "invalid query parameters", 400)
		return
	}
	endLatStr, ok := vars["endlat"]
	endLonStr, ok2 := vars["endlon"]

	if ok == false || ok2 == false {
		http.Error(w, "invalid query parameters", 400)
		return
	}

	startLat, _ := strconv.ParseFloat(startLatStr[0], 64)
	startLon, _ := strconv.ParseFloat(startLonStr[0], 64)

	endLat, _ := strconv.ParseFloat(endLatStr[0], 64)
	endLon, _ := strconv.ParseFloat(endLonStr[0], 64)

	start := data.Coordinate{Lat: startLat, Lon: startLon}
	end := data.Coordinate{Lat: endLat, Lon: endLon}

	// send a dijkstra request

	route, err := dijkstra.GetRoute(start, end)

	if err != nil {
		logrus.Error(err)
		http.Error(w, "error calculating dijkstra", 500)
		return
	}

	getRoute := route.ConvertToJSON()
	routeRaw, err := json.Marshal(getRoute)

	// TODO change to simplified data
	if err != nil {
		logrus.Error(err)
		http.Error(w, "error ", 500)
		return
	}

	w.Write(routeRaw)
}

// RouteStationHandler returns a way with the stations on the way
// parameters: startlat, startlon,endlat,endlon,range(in km)
func (s *Server) RouteStationHandler(w http.ResponseWriter, r *http.Request) {

	vars := r.URL.Query()

	startLatStr, ok := vars["startlat"]
	startLonStr, ok2 := vars["startlon"]

	if ok == false || ok2 == false {
		logrus.Warnf("invalid input %s,%s", startLatStr, startLonStr)
		http.Error(w, "invalid query parameters", 400)
		return
	}
	endLatStr, ok := vars["endlat"]
	endLonStr, ok2 := vars["endlon"]

	if ok == false || ok2 == false {
		http.Error(w, "invalid query parameters", 400)
		return
	}

	rangeKmStr, ok := vars["rangeKm"]

	if ok == false {
		http.Error(w, "invalid query parameters", 400)
		return
	}

	startLat, _ := strconv.ParseFloat(startLatStr[0], 64)
	startLon, _ := strconv.ParseFloat(startLonStr[0], 64)

	endLat, _ := strconv.ParseFloat(endLatStr[0], 64)
	endLon, _ := strconv.ParseFloat(endLonStr[0], 64)

	rangeKm, _ := strconv.ParseFloat(rangeKmStr[0], 64)

	start := data.Coordinate{Lat: startLat, Lon: startLon}
	end := data.Coordinate{Lat: endLat, Lon: endLon}

	// send a dijkstra request

	route, stations, err := dijkstra.GetRouteWithStations(s.graph, s.stations, start, end, rangeKm)

	if err != nil {
		logrus.Error(err)
		http.Error(w, "error calculating dijkstra", 500)
		return
	}

	getRoute := route.ConvertToJSON()

	getRouteWithStations := data.GetRouteWithStations{Route: getRoute.Route, TotalCost: getRoute.TotalCost, Stations: stations}

	routeRaw, err := json.Marshal(getRouteWithStations)

	// TODO change to simplified data
	if err != nil {
		logrus.Error(err)
		http.Error(w, "error ", 500)
		return
	}

	w.Write(routeRaw)
}

func GetArea(r *http.Request) (ne data.Coordinate, sw data.Coordinate, err error) {
	vars := r.URL.Query()

	neLatStr, ok := vars["nelat"]
	neLonStr, ok2 := vars["nelon"]

	if ok == false || ok2 == false {
		logrus.Warnf("invalid input %s,%s", neLatStr, neLonStr)
		err = errors.New("invalid query parameters")
		return

	}
	swLatStr, ok := vars["swlat"]
	swLonStr, ok2 := vars["swlon"]

	if ok == false || ok2 == false {
		err = errors.New("invalid query parameters")
		return

	}

	neLat, _ := strconv.ParseFloat(neLatStr[0], 64)
	neLon, _ := strconv.ParseFloat(neLonStr[0], 64)

	swLat, _ := strconv.ParseFloat(swLatStr[0], 64)
	swLon, _ := strconv.ParseFloat(swLonStr[0], 64)

	ne = data.Coordinate{Lat: neLat, Lon: neLon}
	sw = data.Coordinate{Lat: swLat, Lon: swLon}

	return
}
