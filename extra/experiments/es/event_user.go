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

type UpdateProfile struct {
	ID        string
	FirstName string
	LastName  string
}

func (u UpdateProfile) AggregateType() string {
	return "USER"
}

func (u UpdateProfile) Action() string {
	return "User Updated"
}

func (u UpdateProfile) Version() uint64 {
	return 1
}

func (u UpdateProfile) Apply(agg goes.Aggregate, event goes.Event) {
	user := agg.(*User)
	user.FirstName = u.FirstName
	user.LastName = u.LastName
}

type Create struct {
	ID        string
	FirstName string
	LastName  string
}

func (c Create) AggregateType() string {
	return "USER"
}

func (c Create) Action() string {
	return "User Created"
}

func (c Create) Version() uint64 {
	return 1
}

func (c Create) Apply(agg goes.Aggregate, event goes.Event) {
	user := agg.(*User)
	user.FirstName = c.FirstName
	user.LastName = c.LastName
}
