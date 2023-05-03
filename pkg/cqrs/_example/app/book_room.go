package app

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/dmitrymomot/go-app/pkg/cqrs"
	"github.com/dmitrymomot/go-app/pkg/uuid"
)

// BookRoomHandler is a command handler, which handles BookRoom command and emits RoomBooked.
//
// In CQRS, one command must be handled by only one handler.
// When another handler with this command is added to command processor, error will be retuerned.
type BookRoomHandler struct {
	eventBus    cqrs.EventBus
	bookedRooms *sync.Map
}

// NewBookRoomHandler implements cqrs.CommanfHandlerFactory interface.
func NewBookRoomHandler() cqrs.CommanfHandlerFactory {
	return func(cb cqrs.CommandBus, eb cqrs.EventBus) cqrs.CommandHandler {
		return &BookRoomHandler{
			eventBus:    eb,
			bookedRooms: &sync.Map{},
		}
	}
}

func (b BookRoomHandler) HandlerName() string {
	return "BookRoomHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b BookRoomHandler) NewCommand() interface{} {
	return &BookRoom{}
}

func (b BookRoomHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*BookRoom)

	if _, ok := b.bookedRooms.Load(cmd.RoomId); ok {
		// Room is already booked, we can't book it twice
		// In production, you probably will return some error here
		return nil
	}

	// some random price, in production you probably will calculate in wiser way
	price := (rand.Int63n(40) + 1) * 10

	log.Printf(
		"Booked %s for %s from %s to %s",
		cmd.RoomId,
		cmd.GuestName,
		cmd.StartDate.Format(time.RFC3339),
		cmd.EndDate.Format(time.RFC3339),
	)

	// RoomBooked will be handled by OrderBeerOnRoomBooked event handler,
	// in future RoomBooked may be handled by multiple event handler
	if err := b.eventBus.Publish(ctx, &RoomBooked{
		ReservationId: uuid.New().String(),
		RoomId:        cmd.RoomId,
		GuestName:     cmd.GuestName,
		Price:         price,
		StartDate:     cmd.StartDate,
		EndDate:       cmd.EndDate,
	}); err != nil {
		return err
	}

	return nil
}
