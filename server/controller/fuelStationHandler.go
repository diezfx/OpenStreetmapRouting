package controller

import (
	"OpenStreetmapRouting/data"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func FuelStationHandler(w http.ResponseWriter, r *http.Request) {

	stations := data.GetFuelStations()

	stationsJson, err := json.Marshal(stations)

	if err != nil || stations == nil {

		logrus.Error(err)
		http.Error(w, "no stations available", 500)
	}

	w.Write(stationsJson)

}
