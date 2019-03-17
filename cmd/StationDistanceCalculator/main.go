package main

import (
	"OpenStreetmapRouting/config"
	"OpenStreetmapRouting/data"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

var edges = []data.Edge{{ID: 0, Start: 0, End: 1, Cost: 9},
	{ID: 1, Start: 0, End: 2, Cost: 8},
	{ID: 2, Start: 0, End: 4, Cost: 7},
	{ID: 3, Start: 2, End: 0, Cost: 6},
	{Start: 2, End: 1, Cost: 5}}

func main() {

	// not used so far
	_ = &config.StationDistanceCalcConfig{}

	client, err := data.NewMongoDB()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Connected to DB.")

	collection := client.Database("Routing").Collection("StationDists")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err = collection.InsertOne(ctx, edges[1])

	if err != nil {
		logrus.Fatal(err)
	}

}
