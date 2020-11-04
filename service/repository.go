package service

import (
	"context"
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

func (r *Repository) AddUser(info AuthInfo) error {
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
	return UserAlreadyExistError{}
}

func (r *Repository) AddMeasurement(measurement MeasurementWithAuth) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Authenticate(ctx, *measurement.User); err != nil {
		return err
	}
	_, err := r.Measurements.InsertOne(ctx, bson.M{
		"login": measurement.User.Login,
		"temperature": measurement.Temperature,
		"timestamp": measurement.Timestamp,
	})
	return err
}

func (r *Repository) GetMeasurements(info AuthInfo) ([]Measurement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Authenticate(ctx, info); err != nil {
		return nil, err
	}
	cur, err := r.Measurements.Find(ctx, bson.M{"login": info.Login})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			panic(err)
		}
	}()
	var measurements []Measurement
	for cur.Next(ctx) {
		var measurement Measurement
		if err := cur.Decode(&measurement); err != nil {
			return nil, err
		}
		measurements = append(measurements, measurement)
	}
	return measurements, nil
}

func (r *Repository) Authenticate(ctx context.Context, info AuthInfo) error {
	user := r.Users.FindOne(ctx, bson.M{"login": info.Login})
	if user.Err() != nil {
		return UserNotfoundError{Login: info.Login}
	}
	var realInfo AuthInfo
	if err := user.Decode(&realInfo); err != nil {
		panic(err)
	}
	if ComparePasswords(realInfo.Password, info.Password) {
		return nil
	} else {
		return InvalidPasswordError{Login: info.Login}
	}
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
