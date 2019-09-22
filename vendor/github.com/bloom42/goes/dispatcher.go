package goes

import (
	"context"
	"strconv"
)

// AsyncReactor are reactors which don't care about the event's insertion transaction, they
// are executed asynchronously (in they own goroutine)
type AsyncReactor = func(context.Context, Event)

// SyncReactor are reactors that execute in the same transaction than the event's one and thus can
// fail it in case of error.
type SyncReactor = func(context.Context, Tx, Event) error

// EventMatcher is a func that can match event to a criteria.
type EventMatcher func(Event) bool

// eventBusSubscription is a subscription in the eventbus
type eventBusSubscription struct {
	matcher EventMatcher
	sync    []SyncReactor
	async   []AsyncReactor
}

// eventBus is a global in-memory bus within each event flow before getting saved in event store
var eventBus = []eventBusSubscription{}

// On is used to register `SyncReactor` and `AsyncReactor` to react to `Event`s
func On(matcher EventMatcher, sync []SyncReactor, async []AsyncReactor) {
	if sync == nil {
		sync = []SyncReactor{}
	}
	if async == nil {
		async = []AsyncReactor{}
	}

	subscription := eventBusSubscription{
		matcher: matcher,
		sync:    sync,
		async:   async,
	}

	eventBus = append(eventBus, subscription)
}

func dispatch(ctx context.Context, tx Store, event Event) error {
	for _, subscription := range eventBus {

		if subscription.matcher(event) {
			// dispatch sync reactor synchronously
			// it can be something like a projection
			for _, syncReactor := range subscription.sync {
				if err := syncReactor(ctx, tx, event); err != nil {
					return err
				}
			}

			// dispatch async reactors asynchronously
			for _, asyncReactor := range subscription.async {
				go asyncReactor(ctx, event)
			}
		}
	}
	return nil
}

// MatchEvent matches a specific event type, nil events never match.
func MatchEvent(t EventData) EventMatcher {
	eventType := t.AggregateType() +
		"." + t.Action() +
		"." + strconv.FormatUint(t.Version(), 10)
	return func(e Event) bool {
		eType := e.AggregateType +
			"." + e.Action +
			"." + strconv.FormatUint(e.Version, 10)
		return eventType == eType
	}
}

// MatchAny matches any event
func MatchAny() EventMatcher {
	return func(event Event) bool {
		return true
	}
}

// MatchAggregate matches a specific aggregate type, nil events never match
func MatchAggregate(t Aggregate) EventMatcher {
	return func(event Event) bool {
		data := event.Data.(EventData)
		return data.AggregateType() == t.AggregateType()
	}
}

// MatchAnyOf matches if any of several matchers matches
func MatchAnyOf(matchers ...EventMatcher) EventMatcher {
	return func(e Event) bool {
		for _, m := range matchers {
			if m(e) {
				return true
			}
		}
		return false
	}
}

// MatchAnyEventOf matches if any of several matchers matches
func MatchAnyEventOf(events ...EventData) EventMatcher {
	return func(e Event) bool {
		for _, t := range events {
			if MatchEvent(t)(e) {
				return true
			}
		}
		return false
	}
}
