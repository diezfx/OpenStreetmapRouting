package controller

import (
	"OpenStreetmapRouting/data"
	"OpenStreetmapRouting/dijkstra"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

//FuelStationHandler returns all stations in the area given through the 4 query parameters
func (s *Server) FuelStationHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ne, sw, err := GetArea(r)
		if err != nil {
			logrus.Warnf("query wrong: %s", err)
			http.Error(w, "invalid query parameters", 400)
		}

		stationList := s.stationsGrid.GetNodesInArea(ne, sw)

		stationsGet := data.FuelStationGet{Stations: stationList}

		stationsJSON, err := json.Marshal(stationsGet)

		if err != nil || s.stations == nil {

			logrus.Error(err)
			http.Error(w, "no stations available", 500)
		}

		w.Write(stationsJSON)

	}
}

//ReachableStationsHandler shows all reachable nodes in an area(todo)
func (s *Server) ReachableStationsHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		/*
			ne, sw, err := GetArea(r)
			if err != nil {
				logrus.Warnf("query wrong: %s", err)
				http.Error(w, "invalid query parameters", 400)
			}
		*/

		start, err := GetStart(r)
		if err != nil {
			logrus.Warnf("query wrong: %s", err)
			http.Error(w, "invalid query parameters", 400)
		}

		reachable, unreachable := dijkstra.StationsReachable(s.graph, start)

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
			mesJSON, _ := json.Marshal(mes)
			w.Write(mesJSON)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(stationsJSON)

	}

}
