package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/dmitrymomot/go-app/pkg/eventstore"
	"github.com/dmitrymomot/go-utils"
	faker "github.com/go-faker/faker/v4"
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
	es, err := eventstore.NewEventStore(ctx, db, "user")
	if err != nil {
		panic(err)
	}

	user := NewUser()

	// Create user and save it to the event store
	_, err = es.AppendEvent(ctx, eventstore.NewEvent(
		user.ID,
		"UserCreatedEvent",
		UserCreatedEvent{
			Name:      "John Doe",
			Email:     "johndoe@mail.dev",
			Status:    "active",
			CreatedAt: user.CreatedAt,
		}))
	if err != nil {
		panic(err)
	}

	// Update user name and save it to the event store
	_, err = es.AppendEvent(ctx, eventstore.NewEvent(
		user.ID,
		"UserNameUpdatedEvent",
		UserNameUpdatedEvent{
			Name:      "John Doe Jr.",
			UpdatedAt: time.Now().UnixNano(),
		},
	))
	if err != nil {
		panic(err)
	}

	// Update user email and save it to the event store
	_, err = es.AppendEvent(ctx, eventstore.NewEvent(
		user.ID,
		"UserEmailUpdatedEvent",
		UserEmailUpdatedEvent{
			Email:     "john@doe.dev",
			UpdatedAt: time.Now().UnixNano(),
		},
	))
	if err != nil {
		panic(err)
	}

	// Update user status and save it to the event store
	_, err = es.AppendEvent(ctx, eventstore.NewEvent(
		user.ID,
		"UserStatusUpdatedEvent",
		UserStatusUpdatedEvent{
			Status:    "inactive",
			UpdatedAt: time.Now().UnixNano(),
		},
	))
	if err != nil {
		panic(err)
	}

	// Load user from the event store
	latestState, err := es.LoadCurrentState(ctx, user)
	if err != nil {
		panic(err)
	}

	// Print user data
	utils.PrettyPrint("latestState", latestState)

	// Create a snapshot of the user
	snapshot, err := es.StoreSnapshot(ctx, user)
	if err != nil {
		panic(err)
	}

	// Print snapshot data
	utils.PrettyPrint("snapshot", snapshot)

	// Simulate events
	simulateEvents(ctx, es, 10000, 30)

	startTime := time.Now()
	// Load user from the snapshot
	latestState, err = es.LoadCurrentState(ctx, user)
	if err != nil {
		panic(err)
	}

	// Print user data
	utils.PrettyPrint("load user from snapshot for ", time.Since(startTime).String(), latestState)
}

// simulate events
func simulateEvents(ctx context.Context, es *eventstore.EventStore, usersN, eventsN int64) {
	eventsMap := map[string]func() interface{}{
		"UserNameUpdatedEvent": func() interface{} {
			return UserNameUpdatedEvent{
				Name:      faker.Name(),
				UpdatedAt: time.Now().UnixNano(),
			}
		},
		"UserEmailUpdatedEvent": func() interface{} {
			return UserEmailUpdatedEvent{
				Email:     faker.Email(),
				UpdatedAt: time.Now().UnixNano(),
			}
		},
		"UserStatusUpdatedEvent": func() interface{} {
			return UserStatusUpdatedEvent{
				Status:    faker.Word(),
				UpdatedAt: time.Now().UnixNano(),
			}
		},
	}

	users := make([]*User, 0, usersN)
	for i := int64(0); i < usersN; i++ {
		user := NewUser()
		users = append(users, user)

		// Create user and save it to the event store
		_, err := es.AppendEvent(ctx, eventstore.NewEvent(
			user.ID,
			"UserCreatedEvent",
			UserCreatedEvent{
				Name:      faker.Name(),
				Email:     faker.Email(),
				Status:    "active",
				CreatedAt: user.CreatedAt,
			}))
		if err != nil {
			panic(err)
		}
	}

	for i := int64(0); i < usersN*eventsN; i++ {
		userIdx := rand.Intn(len(users))
		user := users[userIdx]

		event := getRandomMapKey(eventsMap)
		_, err := es.AppendEvent(ctx, eventstore.NewEvent(
			user.ID,
			event,
			eventsMap[event](),
		))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Simulated events:", usersN*eventsN)

	idx := rand.Intn(len(users))
	startTime := time.Now()
	// Load user from the snapshot
	_, err := es.LoadCurrentState(ctx, users[idx])
	if err != nil {
		panic(err)
	}
	fmt.Println("Load random user from the snapshot:", time.Since(startTime).String())
	utils.PrettyPrint(users[idx])
}

// returns random map key from the given map
func getRandomMapKey(m map[string]func() interface{}) string {
	// Create a slice to store the map keys
	keys := make([]string, 0, len(m))

	// Fill the slice with the map keys
	for k := range m {
		keys = append(keys, k)
	}

	// Generate a random index
	idx := rand.Intn(len(keys))

	// Return the key at the generated position
	return keys[idx]
}
