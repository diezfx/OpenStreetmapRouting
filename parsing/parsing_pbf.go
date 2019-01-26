package parsing

import (
	"OpenStreetmapRouting/data"
	"fmt"
	"strconv"
	"strings"

	"github.com/thomersch/gosmparse"
)

// DataHandlerStep1 this handler will save all nodes that are fields and save all edges
// Streaming data will call those functions.
type DataHandlerStep1 struct {
	//Graph data.Graph

	Graph          *data.GraphRaw
	GasStationList *data.GasStations
}

type DataHandlerStep2 struct {
	//Graph data.Graph

	Graph          *data.GraphRaw
	GasStationList *data.GasStations
}

func (d *DataHandlerStep1) InitGraph() {

	nodeList := make([]data.Node, 0, 4400000)
	edgeList := make([]data.Edge, 0, 7500000)
	info := data.MetaInfo{RoadTypes: make(map[string]int, 0)}

	fuelList := data.GasStations{Stations: make(map[int64]data.Node, 0)}

	d.Graph = &data.GraphRaw{NodeIDs: make(map[int64]int64, 0), Nodes: nodeList, Edges: edgeList, Info: info}
	d.GasStationList = &fuelList

}

var unvalidRoadTypes = []string{"footway", "bridleway", "steps", "path", "cycleway", "construction", "track", "planned", "proposed"}

var gasStations_Charging = []string{"charging_station", "fuel"}

func (d *DataHandlerStep1) ReadNode(n gosmparse.Node) {}

//ReadNode filter nodes that are streets
func (d *DataHandlerStep2) ReadNode(n gosmparse.Node) {

	d.Graph.NodeIDMutex.Lock()
	if nodeID, ok := d.Graph.NodeIDs[n.ID]; ok == true {

		if nodeID == -1 {
			// needs testing  this solution or do the second field later

			node := data.Node{ID: int64(len(d.Graph.Nodes)), ID_Osm: n.ID, Lat: n.Lat, Lon: n.Lon}
			d.Graph.Nodes = append(d.Graph.Nodes, node)

			d.Graph.NodeIDs[n.ID] = node.ID

		}

	}

	// either is already there from way parsing; just add infos or create anew
	if node, ok := d.GasStationList.Stations[n.ID]; ok == true {

		node.Lat = n.Lat
		node.Lon = n.Lon

	} else if hTag, _ := n.Tags["amenity"]; contains(gasStations_Charging, hTag) {

		var fuelType data.NodeType

		if hTag == "charging_station" {
			fuelType = data.NodeType_ChargingStation
		} else {
			fuelType = data.NodeType_GasStation

		}

		node := data.Node{ID: n.ID, ID_Osm: n.ID, Lat: n.Lat, Lon: n.Lon, Type: fuelType}
		d.GasStationList.Stations[n.ID] = node
	}
	d.Graph.NodeIDMutex.Unlock()
}
func (d *DataHandlerStep1) ReadWay(w gosmparse.Way) {
	// only take streets
	if hTag, ok := w.Tags["highway"]; ok == true && !contains(unvalidRoadTypes, hTag) {

		d.Graph.NodeIDMutex.Lock()
		d.Graph.Info.RoadTypes[hTag]++
		for _, ID := range w.NodeIDs {
			// todo checken double count

			if _, ok := d.Graph.NodeIDs[ID]; ok == false {

				// placeholder value that no new val is set yet

				d.Graph.NodeIDs[ID] = -1

			}

		}
		d.Graph.NodeIDMutex.Unlock()

		speed := parseSpeed(w)

		for i, ID := range w.NodeIDs[:len(w.NodeIDs)-1] {
			edge := data.Edge{ID: w.ID, Start: ID, End: w.NodeIDs[i+1], Speed: speed}
			d.Graph.AddEdge(edge)

		}

		// if it's not oneway create edges the other way round as well
		if onewayTag, _ := w.Tags["oneway"]; onewayTag != "yes" {
			for i := len(w.NodeIDs) - 1; i > 1; i-- {
				edge := data.Edge{ID: w.ID, Start: w.NodeIDs[i], End: w.NodeIDs[i-1], Speed: speed}
				d.Graph.AddEdge(edge)
			}
		}

	}

	if hTag, _ := w.Tags["amenity"]; contains(gasStations_Charging, hTag) {

		//just add one of the nodes
		// todo find solution
		for _, ID := range w.NodeIDs[:len(w.NodeIDs)-1] {
			//ID is non unique Id(from way), ID_osm is unique id

			var fuelType data.NodeType

			if hTag == "charging_station" {
				fuelType = data.NodeType_ChargingStation
			} else {
				fuelType = data.NodeType_GasStation

			}

			node := data.Node{ID: w.ID, ID_Osm: ID, Type: fuelType}
			d.GasStationList.Stations[ID] = node

		}

	}

}
func (d *DataHandlerStep2) ReadWay(w gosmparse.Way)           {}
func (d *DataHandlerStep1) ReadRelation(r gosmparse.Relation) {}
func (d *DataHandlerStep2) ReadRelation(r gosmparse.Relation) {}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func parseSpeed(way gosmparse.Way) float64 {

	var speedKmh float64

	if speedS, ok := way.Tags["maxspeed"]; ok == false {
		//assume values for special cases

		//none means autobahn
		//convert from string to int or float?

		speedKmh = 50

	} else {

		speedKmhParse, parseError := stringToFloat(speedS)
		speedKmh = speedKmhParse
		// cover special cases

		if parseError != nil {

			//autobahn assume 130kmh
			if speedS == "none" {

				speedKmh = 130
				//schrittgeschwindigkeit
			} else if speedS == "walk" || speedS == "4-7" || speedS == "DE:walk" || speedS == "DE:living_street" || speedS == "Schrittgeschwindigkeit" {
				speedKmh = 7

			} else if speedS == "signals" || speedS == "variable" {
				// hard, maybe estimate based on street type
				//TODO
				// more than one value is indicated
			} else if strings.ContainsAny(speedS, ";,|") {
				//strange double values; just take the first one
				f := func(c rune) bool {
					if c == ';' || c == ',' || c == '|' {
						return true
					}
					return false
				}
				speedS = strings.FieldsFunc(speedS, f)[0]
				speedKmh, _ = stringToFloat(speedS)

				// the speed needs to be converted to kmh
			} else if strings.Contains(speedS, "mph") {
				speedS = strings.Fields(speedS)[0]
				speedMph, _ := stringToFloat(speedS)
				speedKmh = speedMph * 1.609344

				// speed is in kmh already so default
			} else if strings.Contains(speedS, "kph") {
				speedS = strings.Fields(speedS)[0]
				speedKmh, _ = stringToFloat(speedS)

			} else if speedS == "DE:urban" || speedS == "5ß" {
				speedKmh = 50
			} else if speedS == "zone:maxspeed=DE:30" || speedS == "DE:zone30" || speedS == "hgv=30" || speedS == "DE:zone:30" {

				speedKmh = 30

			} else if speedS == "DE:rural" {
				speedKmh = 100
			} else {

				fmt.Print(parseError.Error())

				speedKmh = 50
			}
		}

	}

	//just assume
	//todo better assumption
	if speedKmh <= 0 {

		speedKmh = 50
	}

	return speedKmh / 3.6
}

func stringToFloat(speedS string) (float64, error) {
	speed, err := strconv.Atoi(speedS)
	/*if err != nil {
		fmt.Println(err)

	}*/
	return float64(speed), err
}
