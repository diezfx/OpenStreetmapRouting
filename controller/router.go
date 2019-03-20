package controller

import (
	"OpenStreetmapRouting/config"
	"OpenStreetmapRouting/data"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router       *mux.Router
	config       *config.Config
	graph        *data.GraphProd
	stations     *data.GasStations
	stationsGrid *data.Grid
}

func Start(graph *data.GraphProd, stations *data.GasStations, stationsGrid *data.Grid) {

	s := Server{config: config.GetConfig(), graph: graph, stations: stations, stationsGrid: stationsGrid}

	s.router = mux.NewRouter()

	s.router.HandleFunc("/", CorsHeader(HomeHandler))
	s.router.HandleFunc("/v1/route", CorsHeader(RouteHandler))
	s.router.HandleFunc("/v1/routewithstation", CorsHeader(s.RouteStationHandler))
	s.router.HandleFunc("/v1/route/area", CorsHeader(s.RouteAreaHandler))
	s.router.HandleFunc("/v1/info", CorsHeader(InfoHandler))
	s.router.HandleFunc("/v1/stations", CorsHeader(s.FuelStationHandler()))
	s.router.HandleFunc("/v1/reachablestations", CorsHeader(s.ReachableStationsHandler()))

	logrus.Infof("Server started at localhost:8000 ")
	http.ListenAndServe("localhost:8000", s.router)
}

func CorsHeader(request http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO real logging
		log.Println(r.RequestURI)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		request.ServeHTTP(w, r)

	})
}
