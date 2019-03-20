package controller

import (
	"OpenStreetmapRouting/data"
	"OpenStreetmapRouting/dijkstra"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

//FuelStationHandler returns all stations in the area given through the 4 query parameters
func (s *Server) FuelStationHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := r.URL.Query()

		neLatStr, ok := vars["nelat"]
		neLonStr, ok2 := vars["nelon"]

		if ok == false || ok2 == false {
			logrus.Warnf("invalid input %s,%s", neLatStr, neLonStr)
			http.Error(w, "invalid query parameters", 400)
			return
		}
		swLatStr, ok := vars["swlat"]
		swLonStr, ok2 := vars["swlon"]

		if ok == false || ok2 == false {
			http.Error(w, "invalid query parameters", 400)
			return
		}

		neLat, _ := strconv.ParseFloat(neLatStr[0], 64)
		neLon, _ := strconv.ParseFloat(neLonStr[0], 64)

		swLat, _ := strconv.ParseFloat(swLatStr[0], 64)
		swLon, _ := strconv.ParseFloat(swLonStr[0], 64)

		stationList := s.stationsGrid.GetNodesInArea(data.Coordinate{Lat: neLat, Lon: neLon}, data.Coordinate{Lat: swLat, Lon: swLon})

		stationsGet := data.FuelStationGet{Stations: stationList}

		stationsJson, err := json.Marshal(stationsGet)

		if err != nil || s.stations == nil {

			logrus.Error(err)
			http.Error(w, "no stations available", 500)
		}

		w.Write(stationsJson)

	}
}

//ReachableStationsHandler shows all reachable nodes in an area

func (s *Server) ReachableStationsHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := r.URL.Query()

		neLatStr, ok := vars["nelat"]
		neLonStr, ok2 := vars["nelon"]

		if ok == false || ok2 == false {
			logrus.Warnf("invalid input %s,%s", neLatStr, neLonStr)
			http.Error(w, "invalid query parameters", 400)
			return
		}
		swLatStr, ok := vars["swlat"]
		swLonStr, ok2 := vars["swlon"]

		if ok == false || ok2 == false {
			logrus.Warnf("invalid input %s,%s", swLatStr, swLonStr)
			http.Error(w, "invalid query parameters", 400)
			return
		}

		startLatStr, ok := vars["startlat"]
		startLonStr, ok2 := vars["startlon"]

		if ok == false || ok2 == false {
			logrus.Warnf("invalid input %s,%s", startLatStr, startLonStr)
			http.Error(w, "invalid query parameters", 400)
			return
		}
		/*
			neLat, _ := strconv.ParseFloat(neLatStr[0], 64)
			neLon, _ := strconv.ParseFloat(neLonStr[0], 64)

			swLat, _ := strconv.ParseFloat(swLatStr[0], 64)
			swLon, _ := strconv.ParseFloat(swLonStr[0], 64)
		*/

		startLat, _ := strconv.ParseFloat(startLatStr[0], 64)
		startLon, _ := strconv.ParseFloat(startLonStr[0], 64)

		reachable, unreachable := dijkstra.StationsReachable(data.Coordinate{Lat: startLat, Lon: startLon})

		reachableGet := data.FuelStationGet{Stations: reachable}
		unreachableGet := data.FuelStationGet{Stations: unreachable}

		type stationListGet struct {
			Reachable   data.FuelStationGet `json:"reachable"`
			Unreachable data.FuelStationGet `json:"unreachable"`
		}

		list := stationListGet{Reachable: reachableGet, Unreachable: unreachableGet}

		stationsJSON, err := json.Marshal(list)

		if err != nil || s.stations == nil {

			logrus.Error(err)

			mes := data.Message{Title: "no way found"}

			mesB, _ := json.Marshal(mes)

			w.Write(mesB)

			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(stationsJSON)

	}

}
