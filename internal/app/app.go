package app

import (
	"context"
	"jwt-mongo-auth/internal/controller"
	"jwt-mongo-auth/internal/repository/mongodb"
	"jwt-mongo-auth/internal/service"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Run() {

	log := logrus.New()

	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Fatal("Не удалось поднять базу", err)
	}
	defer db.Disconnect(context.Background())

	if db.Ping(context.Background(), readpref.Primary()); err != nil {
		logrus.Fatal("Не удалось законнектиться с базой", err)
	}

	dataBase := db.Database("test")
	userRepo := mongodb.NewTokenRepo(dataBase)

	serv := service.NewService(log, userRepo)

	c := controller.NewController(log, serv)

	c.Run()

}
