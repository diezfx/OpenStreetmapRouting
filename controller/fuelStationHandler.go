package controller

import (
	"OpenStreetmapRouting/data"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (s *Server) FuelStationHandler() http.HandlerFunc {

	stations := s.stations

	stationsGrid := data.Grid{}

	stationList := make([]data.Node, len(stations.Stations))

	i := 0
	for _, value := range stations.Stations {
		stationList[i] = value
		i++
	}

	stationsGrid.InitGrid(stationList, s.config)

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

		stationList := stationsGrid.GetNodesInArea(data.Coordinate{Lat: neLat, Lon: neLon}, data.Coordinate{Lat: swLat, Lon: swLon})

		stationsGet := data.FuelStationGet{Stations: stationList}

		stationsJson, err := json.Marshal(stationsGet)

		if err != nil || stations == nil {

			logrus.Error(err)
			http.Error(w, "no stations available", 500)
		}

		w.Write(stationsJson)

	}
}
