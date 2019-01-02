package router

import (
	"OpenStreetmapRouting/server/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", CorsHeader(controller.HomeHandler))
	r.HandleFunc("/v1/route", CorsHeader(controller.RouteHandler))

	logrus.Infof("Server startet at localhost:8000 ")
	http.ListenAndServe("localhost:8000", r)
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
