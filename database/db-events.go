package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Event struct {
	Name string
	Tmpl string
}

func Subscribe(ctx context.Context, id string) (<-chan Event, error) {
	events := make(chan Event)
	dbURL, minT, maxT := os.Getenv("DATABASE_URL"), 10*time.Second, time.Minute
	l := pq.NewListener(dbURL, minT, maxT, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	if err := l.Listen(strings.ToLower(id)); err != nil {
		return nil, errors.Wrap(err, "failed to create listener")
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				Unsubscribe(l, id)
				close(events)
				return
			case n := <-l.Notify:
				if n != nil {
					var e Event
					json.NewDecoder(strings.NewReader(n.Extra)).Decode(&e)
					events <- e
				}
			case <-time.After(90 * time.Second):
				go l.Ping()
			}
		}
	}()
	return events, nil
}

func Publish(id, name, tmpl string) error {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(&Event{name, tmpl})
	payload := strings.ReplaceAll(buf.String(), "\n", "")
	_, err := db.Exec(fmt.Sprintf("NOTIFY %s, '%s'", strings.ToLower(id), payload))
	return errors.Wrap(err, "failed to publish")
}

func Unsubscribe(l *pq.Listener, id string) error {
	return errors.Wrap(l.Unlisten(id), "failed to unsubscribe")
}
