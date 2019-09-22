package main

import (
	"context"
	"fmt"
	"github.com/bloom42/goes"
	"log"
)

func main() {
	err := goes.Init("./db")
	if err != nil {
		log.Fatal(err)
	}
	goes.DB.LogMode(true)

	registerUser()
	goes.On(
		goes.MatchEvent(Create{}),
		[]goes.SyncReactor{Log},
		nil,
	)

	var user User

	command := UpdateUser{
		FirstName: "Ehsan",
		LastName:  "Noureddin Moosa",
	}
	metadata := goes.Metadata{
		"request_id": "my_request_id",
	}

	_, err = goes.Execute(context.Background(), command, &user, metadata)
	if err != nil {
		log.Fatal(err)
	}

}

func registerUser() {
	goes.Register(
		&User{},
		UpdateProfile{}, Create{},
	)
}

func Log(ctx context.Context, tx goes.Tx, event goes.Event) error {
	fmt.Println(event.ID, event.Action, event.AggregateType)
	return nil
}
