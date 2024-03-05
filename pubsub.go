package planner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type event struct {
	Name string
	Tmpl string
}

func Subscribe(id string, out chan event) (*pq.Listener, error) {
	minT, maxT := 10*time.Second, time.Minute
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
			case n := <-l.Notify:
				if n != nil {
					var e event
					json.NewDecoder(strings.NewReader(n.Extra)).Decode(&e)
					out <- e
				}
			case <-time.After(90 * time.Second):
				go l.Ping()
			}
		}
	}()
	return l, nil
}

func Publish(id, name, tmpl string) error {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(&event{name, tmpl})
	payload := strings.ReplaceAll(buf.String(), "\n", "")
	_, err := db.Exec(fmt.Sprintf("NOTIFY %s, '%s'", strings.ToLower(id), payload))
	return errors.Wrap(err, "failed to publish")
}

func Unsubscribe(l *pq.Listener, id string) error {
	return errors.Wrap(l.Unlisten(id), "failed to unsubscribe")
}

//// OLD pre-postgres
// var eventBus = struct {
// 	sync.Mutex
// 	Events map[string][]chan event
// }{Events: map[string][]chan event{}}

// func Subscribe(id string, out chan event) {
// 	eventBus.Lock()
// 	defer eventBus.Unlock()
// 	eventBus.Events[id] = append(eventBus.Events[id], out)
// }

// func Publish(id, name, tmpl string) {
// 	eventBus.Lock()
// 	defer eventBus.Unlock()
// 	for _, l := range eventBus.Events[id] {
// 		l <- event{name, tmpl}
// 	}
// }

// func Unsubscribe(id string, out chan event) {
// 	eventBus.Lock()
// 	defer eventBus.Unlock()
// 	for i, l := range eventBus.Events[id] {
// 		if l == out {
// 			eventBus.Events[id] = append(eventBus.Events[id][:i], eventBus.Events[id][i+1:]...)
// 			break
// 		}
// 	}
// }
