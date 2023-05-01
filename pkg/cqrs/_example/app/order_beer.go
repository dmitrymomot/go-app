package app

import (
	"context"
	"log"
	"math/rand"

	"github.com/dmitrymomot/go-app/pkg/cqrs"
	"github.com/pkg/errors"
)

// OrderBeerHandler is a command handler, which handles OrderBeer command and emits BeerOrdered.
// BeerOrdered is not handled by any event handler, but we may use persistent Pub/Sub to handle it in the future.
type OrderBeerHandler struct {
	eventBus cqrs.EventBus
}

// NewOrderBeerHandler implements cqrs.CommanfHandlerFactory interface.
func NewOrderBeerHandler() cqrs.CommanfHandlerFactory {
	return func(cb cqrs.CommandBus, eb cqrs.EventBus) cqrs.CommandHandler {
		return &OrderBeerHandler{
			eventBus: eb,
		}
	}
}

func (o OrderBeerHandler) HandlerName() string {
	return "OrderBeerHandler"
}

func (o OrderBeerHandler) NewCommand() interface{} {
	return &OrderBeer{}
}

func (o OrderBeerHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*OrderBeer)

	if rand.Int63n(10) == 0 {
		// sometimes there is no beer left, command will be retried
		return errors.Errorf("no beer left for room %s, please try later", cmd.RoomId)
	}

	if err := o.eventBus.Publish(ctx, &BeerOrdered{
		RoomId: cmd.RoomId,
		Count:  cmd.Count,
	}); err != nil {
		return err
	}

	log.Printf("%d beers ordered to room %s", cmd.Count, cmd.RoomId)
	return nil
}

// OrderBeerOnRoomBooked is a event handler, which handles RoomBooked event and emits OrderBeer command.
type OrderBeerOnRoomBooked struct {
	commandBus cqrs.CommandBus
}

// NewOrderBeerOnRoomBooked implements cqrs.EventHandlerFactory interface.
func NewOrderBeerOnRoomBooked() cqrs.EventHandlerFactory {
	return func(cb cqrs.CommandBus, eb cqrs.EventBus) cqrs.EventHandler {
		return &OrderBeerOnRoomBooked{
			commandBus: cb,
		}
	}
}

func (o OrderBeerOnRoomBooked) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "OrderBeerOnRoomBooked"
}

func (OrderBeerOnRoomBooked) NewEvent() interface{} {
	return &RoomBooked{}
}

func (o OrderBeerOnRoomBooked) Handle(ctx context.Context, e interface{}) error {
	event := e.(*RoomBooked)

	orderBeerCmd := &OrderBeer{
		RoomId: event.RoomId,
		Count:  rand.Int63n(10) + 1,
	}

	return o.commandBus.Send(ctx, orderBeerCmd)
}
