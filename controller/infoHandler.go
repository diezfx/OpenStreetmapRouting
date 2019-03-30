package controller

import (
	"github.com/diezfx/OpenStreetmapRouting/data"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	info := data.GetGraphInfo()

	infoGet := info.ConverToGetObject()

	infoJSON, err := json.Marshal(infoGet)

	if err != nil || info == nil {

		logrus.Error(err)
		http.Error(w, "no info available", 400)
	}

	w.Write(infoJSON)

}
