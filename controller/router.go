package controller

import (
	"github.com/diezfx/OpenStreetmapRouting/config"
	"github.com/diezfx/OpenStreetmapRouting/data"

	"net/http"

	"github.com/gorilla/handlers"
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

	s.router.HandleFunc("/", HomeHandler)
	s.router.HandleFunc("/v1/route", s.RouteHandler)
	s.router.HandleFunc("/v1/route/node", s.GetNearestNode)
	s.router.HandleFunc("/v1/routewithstation", s.RouteStationHandler)
	s.router.HandleFunc("/v1/route/area", s.RouteAreaHandler)
	s.router.HandleFunc("/v1/route/areareachable", s.RouteAreaReachableHandler)

	s.router.HandleFunc("/v1/info", InfoHandler)
	s.router.HandleFunc("/v1/stations", s.FuelStationHandler())
	s.router.HandleFunc("/v1/reachablestations", s.ReachableStationsHandler())

	logrus.Infof("Server started at localhost:8000 ")
	http.ListenAndServe("localhost:8000", handlers.CORS()(s.router))
}
