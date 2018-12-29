package router

import (
	"OpenStreetmapRouting/server/controller"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", controller.HomeHandler)
	r.HandleFunc("/v1/route", controller.RouteHandler)

	logrus.Infof("Server startet at localhost:8000 ")
	http.ListenAndServe("localhost:8000", r)
}
