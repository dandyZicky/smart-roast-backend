package roast

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewRoastSession(
	rs string,
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	// check if there is already a session,
	var sessionState Session
	row := db.QueryRow("SELECT id, state FROM active_session WHERE id = $1", rs)

	row.Scan(&sessionState.Id, &sessionState.State)

	if (fmt.Sprintf("%d", sessionState.Id) == rs) && (sessionState.State == 1) {
		http.Error(w, "session already started", 400)
		return
	}

	db.Exec("INSERT INTO active_session (id, state) VALUES ($1, $2)", rs, 1)

	fmt.Printf("Session created with id: %s\n", rs)

	clientOptions := mqtt.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	client := mqtt.NewClient(clientOptions)

	con := client.Connect().Wait()
	if !con {
		fmt.Fprint(w, "connection failed to broker")
		return
	}

	topic := fmt.Sprintf("tes_deh/benar/%s", rs)

	var mqttWait bool = true
	sub := client.Subscribe(topic, 1, roastCb(db, &rs, w, &mqttWait)).Wait()
	if !sub {
		http.Error(w, "can not subscribe", 400)
		return
	}

	fmt.Fprint(w, "{message: `init connection`}")
	w.(http.Flusher).Flush()

	for mqttWait {
		time.Sleep(time.Millisecond)
	}

	fmt.Fprintln(w, "Session complete")
	db.Exec("DELETE FROM active_session WHERE id = $1", rs)
	return
	// fmt.Fprint(w, "Roasting session created")
	// w.(http.Flusher).Flush()
}

func roastCb(db *sql.DB, rs *string, w http.ResponseWriter, state *bool) mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		msg := fmt.Sprintf("%s", m.Payload())
		if msg == "-1" {
			_, e := db.Exec("DELETE FROM active_session WHERE id = $1", *rs)
			*state = false

			if e != nil {
				fmt.Printf("Error occured: %s", e.Error())
			}

			c.Disconnect(1000)
			return
		}
		fmt.Fprintf(w, `{"suhu": %s}`, m.Payload())
		w.(http.Flusher).Flush()
	}
}

type Session struct {
	Id    int `json:"id"`
	State int `json:"state"`
}
