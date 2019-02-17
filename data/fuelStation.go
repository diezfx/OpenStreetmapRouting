package data

import (
	"OpenStreetmapRouting/config"
	"encoding/json"
	"io/ioutil"
	"log"
)

type FuelStationGet struct {
	Stations []*Node
}

func GetFuelStations() *GasStations {

	stations := &GasStations{}
	stations.LoadInfo(config.GetConfig())

	return stations
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
