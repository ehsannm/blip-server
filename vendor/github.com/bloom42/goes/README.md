<p align="center">
  <h3 align="center">GOES</h3>
  <p align="center">Go Event Sourcing made easy</p>
</p>

--------

[![GoDoc](https://godoc.org/github.com/bloom42/goes?status.svg)](https://godoc.org/github.com/bloom42/goes)
[![GitHub release](https://img.shields.io/github/release/bloom42/goes.svg)](https://github.com/bloom42/goes/releases)
[![Build Status](https://travis-ci.org/bloom42/goes.svg?branch=master)](https://travis-ci.org/bloom42/goes)

`goes` is an opinionated transactional event sourcing / CQRS framework using PostgreSQL as both event
store and query store.
It handles all the event dispatching, serialization, deserialization, persistence and command execution
logic for you.


1. [Glossary](#glossary)
2. [Dataflow](#dataflow)
3. [Usage](#usage)
4. [Notes](#notes)
5. [Resources](#resources)
6. [Contributing](#contributing)
7. [License](#license)

-------------------


## Glossary

* **Commands**: Commands are responsible for: validating data, validating that the action can
be performed given the current state of the application and Building the event.
A `Command` returns 1 `Event` + optionnaly 1 non persisted event. The non persisted event
can be used to send non hashed tokens to a `SendEmail` reactor for example.

* **Events**: are the source of truth. They are applied to `Aggregates`

* **Aggregates**: represent the current state of the application. They are the read model.

* **Calculators**: are used to update the state of the application. This is the `Apply` method of `EventData`.

* **Reactors**: are used to trigger side effects as events happen. They are registered with the `On` Function. There is `Sync Reactors` which are called synchronously in the `Execute` function, and can fail the transaction if an error occur, and `Async Reactor` which are called asynchronously, and are not checked for error (fire and forget). They are not triggered by the `Apply` method but in the `Execute` function, thus they **are not** triggered when you replay events. You can triggers them when replaying by using `Dispatch(event)`.

* **Event store**: The write model, where events are persisted (PostgreSQL).

* **Query store**: The read model, where aggregates are persisted (PostgreSQL).


## Dataflow

![Dataflow diagram](/docs/dataflow.jpg)

## Usage

You can find the full example in the `examples/user` directory.

At the beggning there was the **noun**.

So we start by declaring an `Aggregate` (a read model).
```go
package main

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bloom42/goes"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// Aggregate definition                                                                           //
////////////////////////////////////////////////////////////////////////////////////////////////////

// User is our aggregate
type User struct {
	goes.BaseAggregate
	FirstName string
	LastName  string
	Addresses addresses `gorm:"type:jsonb;column:addresses"`
}

// AggregateType is our aggregate type
func (user *User) AggregateType() string {
	return "user"
}

// TableName is our postgres table
func (user *User) TableName() string {
	return "users"
}

// a subfield used as a JSONB column
type address struct {
	Country string `json:"country"`
	Region  string `json:"region"`
}

type addresses []address

// Value is used to serialize to SQL
func (a addresses) Value() (driver.Value, error) {
	j, err := json.Marshal(a)
	return j, err
}

// Scan is used to deserialize from SQL
func (a *addresses) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, a)

	}
	return errors.New(fmt.Sprint("failed to unmarshal JSONB from DB", src))
}
```

Then we should describe which kinds of actions (`Event`s) can happen to our `Aggregate`
and **what** this `Events` **change** to our `Aggregates`. Please welcome **verbs**.
`Event` actions are verb in the past tense.

The `Apply` mtehtods are our **Calculators**. They mutate the `Aggregate` states.
```go
////////////////////////////////////////////////////////////////////////////////////////////////////
// Events definition                                                                              //
////////////////////////////////////////////////////////////////////////////////////////////////////

// CreatedV1 is our first event
// json tags should be set because the struct will be serialized as JSON when saved in the eventstore
type CreatedV1 struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Apply our event to an user aggregate
func (eventData CreatedV1) Apply(agg goes.Aggregate, event goes.Event) {
	user := agg.(*User)
	user.ID = eventData.ID
	user.FirstName = eventData.FirstName
	user.LastName = eventData.LastName
	user.CreatedAt = event.Timestamp
	user.Addresses = addresses{}
}

// AggregateType is our target aggregate type
func (CreatedV1) AggregateType() string {
	return "user"
}

// Action is the performed action, in past tense
func (CreatedV1) Action() string {
	return "created"
}

// Version is the event's verion
func (CreatedV1) Version() uint64 {
	return 1
}

// FirstNameUpdatedV1 is our second event
type FirstNameUpdatedV1 struct {
	FirstName string `json:"first_name"`
}

// Apply our event to an user aggregate
func (eventData FirstNameUpdatedV1) Apply(agg goes.Aggregate, event goes.Event) {
	user := agg.(*User)
	user.FirstName = eventData.FirstName
}

// AggregateType is our target aggregate type
func (FirstNameUpdatedV1) AggregateType() string {
	return "user"
}

// Action is the performed action, in past tense
func (FirstNameUpdatedV1) Action() string {
	return "first_name_updated"
}

// Version is the event's verion
func (FirstNameUpdatedV1) Version() uint64 {
	return 1
}
```

then we should describe **how** we can perform these acions (`Event`s): this is our
`Command`s. They are responsible to validate the command against our current state and build the
event.
```go
////////////////////////////////////////////////////////////////////////////////////////////////////
// Commands definition                                                                            //
////////////////////////////////////////////////////////////////////////////////////////////////////

// ValidationError is a custom validation error type
type ValidationError error

// NewValidationError returns a new ValidationError
func NewValidationError(message string) ValidationError {
	return errors.New(message).(ValidationError)
}

func validateFirstName(firstName string) error {
	length := len(firstName)

	if length < 3 {
		return NewValidationError("FirstName is too short")
	} else if length > 42 {
		return NewValidationError("FirstName is too long")
	}
	return nil
}

// Create is our first command to create an user
type Create struct {
	FirstName string
	LastName  string
}

// Validate the command's validity against our business logic and the current application state
func (c Create) Validate(_ context.Context, tx goes.Tx, agg goes.Aggregate) error {
	// user := *agg.(*User)
	// _ = user
	return validateFirstName(c.FirstName)
}

// BuildEvent returns the CreatedV1 event
func (c Create) BuildEvent(context.Context) (goes.EventData, interface{}, error) {
	return CreatedV1{
		ID:        "0563019f-ade9-4cb1-81a7-4f1bb3213cb0",
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}, nil, nil
}

// AggregateType returns the target aggregate type
func (c Create) AggregateType() string {
	return "user"
}

// UpdateFirstName is our second command to update the user's firstname
type UpdateFirstName struct {
	FirstName string
}

// Validate the command's validity against our business logic and the current application state
func (c UpdateFirstName) Validate(_ context.Context, tx goes.Tx, agg goes.Aggregate) error {
	// user := agg.(*User)
	// _ = user
	return validateFirstName(c.FirstName)
}

// BuildEvent returns the FirstNameUpdatedV1 event
func (c UpdateFirstName) BuildEvent(context.Context) (goes.EventData, interface{}, error) {
	return FirstNameUpdatedV1{
		FirstName: c.FirstName,
	}, nil, nil
}

// AggregateType returns the target aggregate type
func (c UpdateFirstName) AggregateType() string {
	return "user"
}
```

then we can add some **reactors** to react to our events.
```go
// We register the User aggregate with all it's associated events
func registerUser() {
	goes.Register(
		&User{},
		FirstNameUpdatedV1{},
		CreatedV1{},
	)
}

func asyncReactorExample(context.Context, goes.Event) {
	time.Sleep(3 * time.Second)
	fmt.Println("hello from async reactor")
}

func syncReactorExample(_ context.Context, _ goes.Tx, event goes.Event) error {
	createdEvent := event.Data.(CreatedV1)
	fmt.Printf("User created: %s %s\n", createdEvent.FirstName, createdEvent.LastName)
	return nil
}
```

And finally we can query both our store and our event store using `goes.DB`.

```go
func queryUserEvents(userID string) {
	query := "SELECT * FROM users_events WHERE aggregate_id = ?;"
	events := []goes.StoreEvent{}
	err := goes.DB.Raw(query, userID).Scan(&events).Error
	if err != nil {
		panic(err)
	}

	ret := make([]goes.Event, len(events))
	for i, event := range events {
		ret[i], err = event.Deserialize()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(ret)
}


func main() {
	// configure the database
	err := goes.Init(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	goes.DB.LogMode(true)

	registerUser()
	goes.On(
		goes.MatchEvent(CreatedV1{}),
		[]goes.SyncReactor{syncReactorExample},
		[]goes.AsyncReactor{asyncReactorExample},
	)

	var user User

	command := Create{
		FirstName: "Sylvain",
		LastName:  "Kerkour",
	}
	metadata := goes.Metadata{
		"request_id": "my_request_id",
	}

	_, err = goes.Execute(context.Background(), command, &user, metadata)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
	// User {
	// 	ID: "0563019f-ade9-4cb1-81a7-4f1bb3213cb0",
	// 	FirstName: "Sylvain",
	// 	LastName: "Kerkour",
	// }

	queryUserEvents(user.ID)
	time.Sleep(5 * time.Second) // for the async reactor
}
```


## Notes

`Apply` methods should return a pointer
`Validate` methods take a pointer as input


## Resources

This implementation is sort of the Go implementation of the following event sourcing framework

* https://kickstarter.engineering/event-sourcing-made-simple-4a2625113224
* https://github.com/mishudark/eventhus
* https://github.com/looplab/eventhorizon


## Contributing

See [https://opensource.bloom.sh/contributing](https://opensource.bloom.sh/contributing)


## License

See `LICENSE.txt` and [https://opensource.bloom.sh/licensing](https://opensource.bloom.sh/licensing)
