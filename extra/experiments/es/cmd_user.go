package main

import (
	"context"
	"github.com/bloom42/goes"
)

/*
   Creation Time: 2019 - Jul - 26
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type UpdateUser struct {
	FirstName string
	LastName  string
}

func (u UpdateUser) BuildEvent(ctx context.Context) (event goes.EventData, nonPersisted interface{}, err error) {
	return UpdateProfile{
		ID:        "",
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil, nil
}

func (u UpdateUser) Validate(ctx context.Context, tx goes.Tx, agg goes.Aggregate) error {
	return nil
}

func (u UpdateUser) AggregateType() string {
	return "USER"
}
