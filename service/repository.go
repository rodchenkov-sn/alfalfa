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
}

type Rights struct {
	Write bool
	Read  bool
}

type Repository struct {
	Client        *mongo.Client
	Users         *mongo.Collection
	Measurements  *mongo.Collection
}

func (r *Repository) AddSupervisor(login string, supervisor common.Supervisor) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.Users.UpdateOne(
		ctx,
		bson.M{"login": login},
		bson.M{
			"$push": bson.M{"supervisors": supervisor.Login},
		},
	)
	return err
}

func (r *Repository) GetRights(issuer string, subject string) Rights {
	if issuer == subject {
		return Rights{Write: true, Read: true}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := r.Users.FindOne(ctx, bson.M{"login": subject})
	if user.Err() != nil {
		return Rights{Read: false, Write: false}
	}
	var realInfo common.UserInfo
	if err := user.Decode(&realInfo); err != nil {
		panic(err)
	}
	for _, supervisor := range realInfo.Supervisors {
		if supervisor == issuer {
			return Rights{Read: true, Write: true}
		}
	}
	return Rights{Read: false, Write: false}
}

func (r *Repository) Authenticate(credentials common.Credentials) AuthResult {
	return r.authenticate(credentials.Login, credentials.Password)
}

func (r *Repository) authenticate(login string, password string) AuthResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := r.Users.FindOne(ctx, bson.M{"login": login})
	if user.Err() != nil {
		return AuthResult{Login: login, Exist: false}
	}
	var realInfo common.UserInfo
	if err := user.Decode(&realInfo); err != nil {
		panic(err)
	}
	if ComparePasswords(realInfo.Password, password) {
		return AuthResult{Login: login, Exist: true, ValidPassword: true}
	} else {
		return AuthResult{Login: login, Exist: true, ValidPassword: false}
	}
}

func (r *Repository) AddUser(info common.UserInfo) error {
	hashedPassword, err := HashPassword(info.Password)
	if err != nil {
		return err
	}
	authResult := r.authenticate(info.Login, info.Password)
	if authResult.Exist {
		return common.UserAlreadyExistError{Login: info.Login}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	supervisors := []string{} // supervisors array must have address
	if info.Supervisors != nil {
		supervisors = info.Supervisors
	}
	_, err = r.Users.InsertOne(ctx, bson.M{
		"login":       info.Login,
		"password":    hashedPassword,
		"supervisors": supervisors,
	})
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
	cur, err := r.Measurements.Find(ctx, bson.M{"login": login}, options.Find().SetSort(bson.M{
		"timestamp": -1,
	}))
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