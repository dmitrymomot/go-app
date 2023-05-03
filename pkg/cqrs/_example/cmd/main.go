package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrymomot/go-app/pkg/cqrs"
	"github.com/dmitrymomot/go-app/pkg/cqrs/_example/app"
	"github.com/dmitrymomot/go-env"
	"github.com/dmitrymomot/go-utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var simulateTraffic = env.GetBool("SIMULATE_TRAFFIC", false)

func main() {
	logger := logrus.WithFields(logrus.Fields{
		"app":       "cqrs-example",
		"component": "main",
	})
	defer func() { logger.Info("Server successfully shutdown") }()

	redisOptions, err := redis.ParseURL("redis://redis:6379/0")
	if err != nil {
		logger.WithError(err).Fatal("Cannot parse redis URL")
	}
	redisClient := redis.NewClient(redisOptions)
	defer redisClient.Close()

	commandsPublisher, err := cqrs.NewPublisher(redisClient,
		cqrs.NewLogrusWrapper(logger.WithField("component", "cqrs-commands-publisher")),
	)
	if err != nil {
		logger.WithError(err).Fatal("Cannot create commands publisher")
	}
	defer commandsPublisher.Close()

	commandsSubscriber, err := cqrs.NewSubscriber(redisClient, "example-commands",
		cqrs.NewLogrusWrapper(logger.WithField("component", "cqrs-commands-subscriber")),
	)
	if err != nil {
		logger.WithError(err).Fatal("Cannot create commands subscriber")
	}
	defer commandsSubscriber.Close()

	// Events will be published to PubSub configured Redis, because they may be consumed by multiple consumers.
	// (in that case BookingsFinancialReport and OrderBeerOnRoomBooked).
	eventsPublisher, err := cqrs.NewPublisher(redisClient,
		cqrs.NewLogrusWrapper(logger.WithField("component", "cqrs-events-publisher")),
	)
	if err != nil {
		logger.WithError(err).Fatal("Cannot create events publisher")
	}
	defer eventsPublisher.Close()

	router, err := cqrs.NewRouter(cqrs.NewLogrusWrapper(logger.WithField("component", "cqrs-router")), 10)
	if err != nil {
		logger.WithError(err).Fatal("Cannot create router")
	}
	defer router.Close()

	// cqrs.Facade is facade for Command and Event buses and processors.
	// You can use facade, or create buses and processors manually (you can inspire with cqrs.NewFacade)
	cqrsFacade, err := cqrs.NewFacade(
		redisClient,
		cqrs.NewLogrusWrapper(logger.WithField("component", "cqrs-facade")),
		router,
		commandsPublisher, eventsPublisher, commandsSubscriber,
		[]cqrs.CommanfHandlerFactory{
			app.NewBookRoomHandler(),
			app.NewOrderBeerHandler(),
		}, []cqrs.EventHandlerFactory{
			app.NewBookingsFinancialReport(),
			app.NewOrderBeerOnRoomBooked(),
		},
	)
	if err != nil {
		logger.WithError(err).Fatal("Cannot create cqrs facade")
	}

	if simulateTraffic {
		// publish BookRoom commands every second to simulate incoming traffic
		go publishCommands(cqrsFacade.CommandBus(), logger.WithField("component", "publishCommands"))
	}

	// processors are based on router, so they will work when router will start
	if err := router.Run(context.Background()); err != nil {
		logger.WithError(err).Fatal("Cannot run router")
	}
}

func publishCommands(commandBus cqrs.CommandBus, logger *logrus.Entry) func() {
	i := 0
	for {
		i++

		startDate := time.Now().Add(time.Hour * 24 * 2)
		endDate := startDate.Add(time.Hour * 24 * 3)

		bookRoomCmd := &app.BookRoom{
			RoomId:    fmt.Sprintf("%d", i),
			GuestName: "John",
			StartDate: utils.Pointer(startDate),
			EndDate:   utils.Pointer(endDate),
		}
		if err := commandBus.Send(context.Background(), bookRoomCmd); err != nil {
			logger.WithError(err).Error("Cannot send BookRoom command")
		}

		time.Sleep(time.Second)
	}
}
