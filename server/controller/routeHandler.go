package controller

import (
	"OpenStreetmapRouting/data"
	"OpenStreetmapRouting/dijkstra"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
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
v	
	if err != nil {
		logrus.Error(err)
		http.Error(w, "error calculating dijkstra", 500)
		return
	}

	getRoute := route.ConvertToJson()



	routeRaw, err := json.Marshal(getRoute)

	// TODO change to simplified data
	if err != nil {
		logrus.Error(err)
		http.Error(w, "error ", 500)
		return
	}

	w.Write(routeRaw)
}
