package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/diezfx/OpenStreetmapRouting/data"
	"github.com/diezfx/OpenStreetmapRouting/dijkstra"

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

	//get all edges that come from these nodes
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

func (s *Server) GetNearestNode(w http.ResponseWriter, r *http.Request) {

	vars := r.URL.Query()

	lat, err := parseCoord(vars, "lat")

	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad input", 400)
		return
	}
	lon, err := parseCoord(vars, "lon")
	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad input", 400)
		return
	}

	node := s.graph.Grid.FindNextNode(lat, lon, false)

	jsonData, err := json.Marshal(node)

	if err != nil {

		logrus.Error(err)
		http.Error(w, "oops", 500)
		return
	}

	w.Write(jsonData)

}

// RouteAreaReachableHandler take an area as argument and return all edges in this area red:unreachable, blue:reachable nodes
func (s *Server) RouteAreaReachableHandler(w http.ResponseWriter, r *http.Request) {

	ne, sw, err := GetArea(r)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	nodeList := s.graph.Grid.GetNodesInArea(ne, sw)

	edgeList := make([]*data.Edge, 0)

	start, err := GetStart(r)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), 400)
		logrus.Debug(start)
		return
	}

	startNode := s.graph.Grid.FindNextNode(start.Lat, start.Lon, false)

	edgesCosts, err := dijkstra.CalcDijkstraToMany(s.graph, *startNode)

	//get all edges that come from these nodes
	for _, node := range nodeList {
		nodeID := node.ID
		edgeBegin := s.graph.Offset[nodeID]
		edgeEnd := s.graph.Offset[nodeID+1]

		for i := edgeBegin; i < edgeEnd; i++ {

			edgeList = append(edgeList, &s.graph.Edges[i])

		}
	}

	areaGetReachable, areaGetUnreachable := data.ConvertAreaToJSONReachable(edgeList, s.graph, edgesCosts)

	areaGet := [2]data.GeoJSONArea{areaGetReachable, areaGetUnreachable}

	jsonData, err := json.Marshal(areaGet)

	if err != nil || jsonData == nil {

		logrus.Error(err)
		http.Error(w, "no area available", 500)
	}

	w.Write(jsonData)

}

func (s *Server) RouteHandler(w http.ResponseWriter, r *http.Request) {

	vars := r.URL.Query()

	start, err := GetStart(r)
	if err != nil {
		logrus.Warnf("invalid start: %s", err)
		http.Error(w, err.Error(), 400)
		return
	}

	endLat, err := parseCoord(vars, "endlat")

	if err != nil {
		logrus.Warnf("invalid query parameter: %s", err)
		http.Error(w, err.Error(), 400)
		return
	}

	endLon, err := parseCoord(vars, "endlon")
	if err != nil {
		logrus.Warnf("invalid query parameter: %s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	end := data.Coordinate{Lat: endLat, Lon: endLon}

	// send a dijkstra request

	route, err := dijkstra.GetRoute(s.graph, start, end)

	if err != nil {
		logrus.Error(err)
		mes := data.Message{Title: "no way found"}

		mesB, _ := json.Marshal(mes)
		w.WriteHeader(http.StatusInternalServerError)

		w.Write(mesB)
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

	start, err := GetStart(r)
	if err != nil {
		logrus.Warnf("invalid start: %s", err)
		http.Error(w, err.Error(), 400)
		return
	}

	endLat, err := parseCoord(vars, "endlat")

	if err != nil {
		logrus.Warnf("invalid query parameter: %s", err)
		http.Error(w, err.Error(), 400)
		return
	}

	endLon, err := parseCoord(vars, "endlon")
	if err != nil {
		logrus.Warnf("invalid query parameter: %s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	end := data.Coordinate{Lat: endLat, Lon: endLon}

	rangeKm, err := parseCoord(vars, "rangeKm")
	if err != nil {
		logrus.Warnf("invalid query parameter: %s", err)
		http.Error(w, err.Error(), 400)
		return
	}

	// send a dijkstra request

	route, stations, err := dijkstra.GetRouteWithStations(s.graph, s.stations, start, end, rangeKm, s.config)

	if err != nil {
		logrus.Error(err)
		http.Error(w, "error calculating dijkstra", 500)
		return
	}

	getRoute := route.ConvertToJSON()

	getRouteWithStations := data.GetRouteWithStations{Route: getRoute.Route, TotalCost: getRoute.TotalCost, Stations: stations}

	routeRaw, err := json.Marshal(getRouteWithStations)

	if err != nil {
		logrus.Error(err)
		http.Error(w, "error ", 500)
		return
	}

	w.Write(routeRaw)
}

func parseCoord(vars url.Values, varName string) (float64, error) {

	varStr, ok := vars[varName]

	if ok == false {
		errText := fmt.Sprintf("variable: %s doesn't exist", varName)
		return math.NaN(), errors.New(errText)
	}

	varFloat, err := strconv.ParseFloat(varStr[0], 64)
	if err != nil {
		return 0, err
	}

	return varFloat, nil

}

func GetStart(r *http.Request) (data.Coordinate, error) {

	vars := r.URL.Query()

	startLat, err := parseCoord(vars, "startlat")
	if err != nil {
		return data.Coordinate{}, err
	}
	startLon, err := parseCoord(vars, "startlon")
	if err != nil {
		return data.Coordinate{}, err
	}

	start := data.Coordinate{Lat: startLat, Lon: startLon}
	return start, nil

}

func GetArea(r *http.Request) (ne data.Coordinate, sw data.Coordinate, err error) {
	vars := r.URL.Query()

	neLon, err := parseCoord(vars, "nelon")
	if err != nil {
		return
	}
	neLat, err := parseCoord(vars, "nelat")
	if err != nil {
		return
	}

	swLon, err := parseCoord(vars, "swlon")
	if err != nil {
		return
	}
	swLat, err := parseCoord(vars, "swlat")
	if err != nil {
		return
	}

	ne = data.Coordinate{Lat: neLat, Lon: neLon}
	sw = data.Coordinate{Lat: swLat, Lon: swLon}

	return
}
