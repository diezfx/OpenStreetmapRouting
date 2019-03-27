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

	stations := &GasStations{}
	stations.LoadInfo(config.GetConfig())

	return stations
}

func (s *GasStations) GetStationsAsList() []Node {
	stationList := make([]Node, len(s.Stations))
	i := 0
	for _, value := range s.Stations {
		stationList[i] = value
		i++
	}
	return stationList

}

func (s *GasStations) SetStations(stations map[int64]Node) {
	s.Stations = stations
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

func (s *GasStations) ConverToGetObject() *FuelStationGet {

	stationJSONGet := FuelStationGet{make([]*Node, len(s.Stations))}
	i := 0

	for _, v := range s.Stations {

		stationJSONGet.Stations[i] = &v
		i++
	}

	return &stationJSONGet

}

func (g *GasStations) WriteFile(config *config.Config) {
	var encodedStations []byte

	jsonGraph, err := json.Marshal(g)
	encodedStations = jsonGraph
	if err != nil {
		logrus.Fatal(err)
		return
	}

	ioutil.WriteFile(config.FuelStationsFilename, encodedStations, 0644)

}
