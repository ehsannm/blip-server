package main

import (
	"context"
	"fmt"
	"git.ronaksoftware.com/blip/server/internal/tools"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
	"log"
)

var (
	_DG *dgo.Dgraph
)

type Post struct {
	UID       string   `json:"uid"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	GroupUIDs []string `json:"group_uids"`
}

func (p *Post) setQuery() []*api.NQuad {
	subject := "_:post"
	if p.UID != "" {
		subject = p.UID
	}
	nq := []*api.NQuad{
		{Subject: subject, Predicate: "title", ObjectId: "", ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: p.Title}}},
		{Subject: subject, Predicate: "body", ObjectId: "", ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: p.Body}}},
		{Subject: subject, Predicate: "type", ObjectId: "", ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: "POST"}}},
	}
	for idx := range p.GroupUIDs {
		nq = append(nq,
			&api.NQuad{Subject: subject, Predicate: "group", ObjectId: p.GroupUIDs[idx]},
		)
	}
	return nq
}
func (p *Post) Save(ctx context.Context) error {
	txn := _DG.NewTxn()
	defer txn.Discard(ctx)

	assigned, err := txn.Mutate(ctx, &api.Mutation{
		Set:       p.setQuery(),
		CommitNow: true,
	})
	if err != nil {
		return err
	}
	p.UID = assigned.Uids["post"]
	return nil
}

type Group struct {
	UID        string   `json:"uid"`
	Title      string   `json:"title"`
	MemberUIDs []string `json:"member_uids"`
}

func (g *Group) setQuery() []*api.NQuad {
	subject := "_:group"
	if g.UID != "" {
		subject = g.UID
	}
	nq := []*api.NQuad{
		{Subject: subject, Predicate: "title", ObjectId: "", ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: g.Title}}},
		{Subject: subject, Predicate: "type", ObjectId: "", ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: "GROUP"}}},
	}
	for idx := range g.MemberUIDs {
		nq = append(nq,
			&api.NQuad{Subject: subject, Predicate: "member", ObjectId: g.MemberUIDs[idx]},
		)
	}
	return nq
}
func (g *Group) Save(ctx context.Context) error {
	txn := _DG.NewTxn()
	defer txn.Discard(ctx)

	assigned, err := txn.Mutate(ctx, &api.Mutation{
		Set:       g.setQuery(),
		CommitNow: true,
	})
	if err != nil {
		return err
	}
	g.UID = assigned.Uids["group"]
	return nil
}
func (g *Group) AddMember(ctx context.Context, userID string) error {
	txn := _DG.NewTxn()
	defer txn.Discard(ctx)

	_, err := txn.Mutate(ctx, &api.Mutation{
		Set:       []*api.NQuad{{Subject: g.UID, Predicate: "<member>", ObjectId: userID}},
		CommitNow: true,
	})

	return err
}

type User struct {
	UID       string   `json:"uid"`
	Username  string   `json:"username"`
	Phone     string   `json:"phone"`
	GroupUIDs []string `json:"group_uids"`
}

func (u User) setQuery() []*api.NQuad {
	subject := "_:user"
	if u.UID != "" {
		subject = u.UID
	}
	nq := []*api.NQuad{
		{Subject: subject, Predicate: "username", ObjectId: "", ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: u.Username}}},
		{Subject: subject, Predicate: "phone", ObjectId: "", ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: u.Phone}}},
		{Subject: subject, Predicate: "type", ObjectId: "", ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: "USER"}}},
	}
	for idx := range u.GroupUIDs {
		nq = append(nq,
			&api.NQuad{Subject: subject, Predicate: "group", ObjectId: u.GroupUIDs[idx]},
		)
	}
	return nq
}
func (u *User) Save(ctx context.Context) error {
	txn := _DG.NewTxn()
	defer txn.Discard(ctx)

	assigned, err := txn.Mutate(ctx, &api.Mutation{
		Set:       u.setQuery(),
		CommitNow: true,
	})
	if err != nil {
		return err
	}
	u.UID = assigned.Uids["user"]
	return nil
}

func main() {
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_DG = dgo.NewDgraphClient(api.NewDgraphClient(conn))

	insert()

}

func insert() {
	for i := 0; i < 10; i++ {
		u := User{
			Username:  fmt.Sprintf("user(%d)", tools.RandomInt64(0)),
			Phone:     tools.RandomDigit(10),
			GroupUIDs: nil,
		}
		err := u.Save(context.Background())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(u)
		}
	}

	for i := 0; i < 10; i++ {
		g := Group{
			Title:      fmt.Sprintf("Group (%d)", tools.RandomInt64(0)),
			MemberUIDs: nil,
		}
		err := g.Save(context.Background())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(g)
		}
		// g.AddMember(context.Background())
	}

}

// func get(ctx context.Context) error {
// 	tx := _DG.NewTxn()
// 	defer tx.Discard(ctx)
//
// 	resp, err := tx.Query(ctx, `
// 	{
// 		getGroup(func:
// 	}
// 	`)
// 	if err != nil {
// 		return err
// 	}
//
//
//
// 	return nil
// }
