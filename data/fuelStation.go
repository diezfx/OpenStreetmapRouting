package data

import (
	"OpenStreetmapRouting/config"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/sirupsen/logrus"
)

type FuelStationGet struct {
	Stations []*Node
}

func GetFuelStations() *GasStations {
	if stations == nil {
		logrus.Warnf("Info not initialized")

		stations = &GasStations{}
		stations.LoadInfo(config.GetConfig())

	}
	return stations
}

func (s *GasStations) LoadInfo(conf *config.Config) {

	dat, err := ioutil.ReadFile(conf.FuelStationsFilename)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(dat, s)

	if err != nil {
		log.Fatal(err.Error())
	}

}

func (g *GasStations) ConverToGetObject() *FuelStationGet {

	stationJsonGet := FuelStationGet{make([]*Node, len(g.Stations))}
	i := 0

	for _, v := range g.Stations {

		stationJsonGet.Stations[i] = v
		i++
	}

	return &stationJsonGet

}
