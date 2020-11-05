package service

import (
	"context"
	"github.com/rodchenkov-sn/alfalfa/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Repository struct {
	Client       *mongo.Client
	Settings     RepositorySettings
	Users        *mongo.Collection
	Measurements *mongo.Collection
}

func (r *Repository) Authenticate(info common.AuthInfo) (exist bool, valid bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := r.Users.FindOne(ctx, bson.M{"login": info.Login})
	if user.Err() != nil {
		return false, false
	}
	var realInfo common.AuthInfo
	if err := user.Decode(&realInfo); err != nil {
		panic(err)
	}
	if ComparePasswords(realInfo.Password, info.Password) {
		return true, true
	} else {
		return true, false
	}
}

func (r *Repository) AddUser(info common.AuthInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if r.Users.FindOne(ctx, bson.M{"login": info.Login}).Err() != nil {
		hashedPassword, err := HashPassword(info.Password)
		if err != nil {
			return err
		}
		_, err = r.Users.InsertOne(ctx, bson.M{"login": info.Login, "password": hashedPassword})
		return err
	}
	return common.UserAlreadyExistError{}
}

func (r *Repository) AddMeasurement(login string, measurement common.Measurement) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.Measurements.InsertOne(ctx, bson.M{
		"login": login,
		"temperature": measurement.Temperature,
		"timestamp": measurement.Timestamp,
	})
	return err
}

func (r *Repository) GetMeasurements(login string) ([]common.Measurement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := r.Measurements.Find(ctx, bson.M{"login": login})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			panic(err)
		}
	}()
	var measurements []common.Measurement
	for cur.Next(ctx) {
		var measurement common.Measurement
		if err := cur.Decode(&measurement); err != nil {
			return nil, err
		}
		measurements = append(measurements, measurement)
	}
	return measurements, nil
}

func (r *Repository) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.Client.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}

func NewRepository(settings RepositorySettings) (*Repository, error) {

	client, err := connectToMongo(settings)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	users := client.
		Database(settings.UsersSettings.Database).
		Collection(settings.UsersSettings.Collection)
	measurements := client.
		Database(settings.MeasurementsSettings.Database).
		Collection(settings.MeasurementsSettings.Collection)
	return &Repository{
		Client:       client,
		Settings:     settings,
		Users:        users,
		Measurements: measurements,
	}, nil
}

func connectToMongo(settings RepositorySettings) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(settings.Uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
