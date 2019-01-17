package controller

import (
	"OpenStreetmapRouting/data"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	info := data.GetGraphInfo()

	infoGet := info.ConverToGetObject()

	infoJson, err := json.Marshal(infoGet)

	if err != nil || info == nil {

		logrus.Error(err)
		http.Error(w, "no info available", 400)
	}

	w.Write(infoJson)

}
