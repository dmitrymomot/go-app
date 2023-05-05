package main

import (
	"database/sql"
	"time"

	"github.com/dmitrymomot/go-app/pkg/eventstore"
	"github.com/dmitrymomot/go-utils"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sirupsen/logrus"
)

const dbConnString = "postgresql://pguser:pgpass@127.0.0.1/pgdb?sslmode=disable"

func main() {
	logger := logrus.WithFields(logrus.Fields{
		"app":       "eventstore-example",
		"component": "main",
	})
	defer func() { logger.Info("Server successfully shutdown") }()

	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Create a context with a timeout and set the Server's context
	ctx, cancel := utils.NewContextWithCancel(logger.WithField("component", "context"))
	defer cancel()

	// Init event store
	es, err := eventstore.NewEventStore(ctx, db, "users2")
	if err != nil {
		panic(err)
	}

	user := NewUser()

	// Create user and save it to the event store
	_, err = es.AppendEvent(ctx, user.ID, UserCreatedEvent{
		Name:      "John Doe",
		Email:     "johndoe@mail.dev",
		Status:    "active",
		CreatedAt: user.CreatedAt,
	})
	if err != nil {
		panic(err)
	}

	// Update user name and save it to the event store
	_, err = es.AppendEvent(ctx, user.ID, UserNameUpdatedEvent{
		Name:      "John Doe Jr.",
		UpdatedAt: time.Now().UnixNano(),
	})
	if err != nil {
		panic(err)
	}

	// Update user email and save it to the event store
	_, err = es.AppendEvent(ctx, user.ID, UserEmailUpdatedEvent{
		Email:     "john@doe.dev",
		UpdatedAt: time.Now().UnixNano(),
	})
	if err != nil {
		panic(err)
	}

	// Update user status and save it to the event store
	_, err = es.AppendEvent(ctx, user.ID, UserStatusUpdatedEvent{
		Status:    "inactive",
		UpdatedAt: time.Now().UnixNano(),
	})
	if err != nil {
		panic(err)
	}

	startTime := time.Now()
	// Load user from the event store
	if err := es.LoadCurrentState(ctx, user.ID, user); err != nil {
		panic(err)
	}
	// Print user data
	utils.PrettyPrint("load user current state without snaphot for ", time.Since(startTime).String(), user)

	// Create a snapshot of the user
	if err := es.StoreSnapshot(ctx, user.ID, user); err != nil {
		panic(err)
	}

	startTime = time.Now()
	// Load user from the snapshot
	if err := es.LoadCurrentState(ctx, user.ID, user); err != nil {
		panic(err)
	}

	// Print user data
	utils.PrettyPrint("load user from snapshot for ", time.Since(startTime).String(), user)
}
