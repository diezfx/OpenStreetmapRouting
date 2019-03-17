package data

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	// needed for sqlx
)

// DB wraps sql DB to extend it with own SQL queries
type DB struct {
	*mongo.Client
}

// NewDB creates a wrapper around sqlx to shorten common sql requests
//todo use config
/*
func NewDB(config *config.StationDistanceCalcConfig) (*DB, error) {
	db, err := sqlx.Open("postgres", "user=root password=root host=localhost port=5432 database=Routing sslmode=disable")

	if err != nil {
		return nil, err
	}

	return &DB{db}, err
}
*/
func NewMongoDB() (*DB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:8081"))

	if err != nil {
		logrus.Fatal(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	return &DB{client}, err
}
