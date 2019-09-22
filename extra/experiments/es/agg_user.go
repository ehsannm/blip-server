package main

import "github.com/bloom42/goes"

/*
   Creation Time: 2019 - Jul - 26
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type User struct {
	goes.BaseAggregate
	FirstName string
	LastName  string
}

func (u User) AggregateType() string {
	return "USER"
}

func (u User) TableName() string {
	return "USERS"
}
