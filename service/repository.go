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

type AuthResult struct {
	Login         string
	Exist         bool
	ValidPassword bool
	Organization  bool
	Policy        *common.OrganizationPolicy
}

type Repository struct {
	Client        *mongo.Client
	Users         *mongo.Collection
	Measurements  *mongo.Collection
}

func (r *Repository) Authenticate(credentials common.Credentials) AuthResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := r.Users.FindOne(ctx, bson.M{"login": credentials.Login})
	if user.Err() != nil {
		return AuthResult{Login: credentials.Login, Exist: false}
	}
	var realInfo common.Credentials
	if err := user.Decode(&realInfo); err != nil {
		panic(err)
	}
	if ComparePasswords(realInfo.Password, credentials.Password) {
		return AuthResult{Login: credentials.Login, Exist: true, ValidPassword: true}
	} else {
		return AuthResult{Login: credentials.Login, Exist: true, ValidPassword: false}
	}
}

func (r *Repository) AddUser(credentials common.Credentials) error {
	hashedPassword, err := HashPassword(credentials.Password)
	if err != nil {
		return err
	}
	return r.insertIntoUsers(credentials, bson.M{
		"login":    credentials.Login,
		"password": hashedPassword,
	})
}

func (r* Repository) AddOrganization(organization common.Organization) error {
	hashedPassword, err := HashPassword(organization.Credentials.Password)
	if err != nil {
		return err
	}
	return r.insertIntoUsers(organization.Credentials, bson.M{
		"login":    organization.Credentials.Login,
		"password": hashedPassword,
		"policy":   organization.Policy,
	})
}

func (r *Repository) insertIntoUsers(credentials common.Credentials, document bson.M) error {
	authResult := r.Authenticate(credentials)
	if authResult.Exist {
		return common.UserAlreadyExistError{Login: credentials.Login}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.Users.InsertOne(ctx, document)
	return err
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
		Database(settings.UsersPath.Database).
		Collection(settings.UsersPath.Collection)
	measurements := client.
		Database(settings.MeasurementsPath.Database).
		Collection(settings.MeasurementsPath.Collection)
	return &Repository{
		Client:        client,
		Users:         users,
		Measurements:  measurements,
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
