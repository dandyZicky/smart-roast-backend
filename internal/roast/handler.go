package roast

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewRoastSession(rs string, w *http.ResponseWriter, db *sql.DB) {
	// check if there is already a session,
	var sessionState Session
	row := db.QueryRow("SELECT id, state FROM active_session WHERE id = $1", rs)

	row.Scan(&sessionState.Id, &sessionState.State)

	if (fmt.Sprintf("%d", sessionState.Id) == rs) && (sessionState.State == 1) {
		http.Error(*w, "session already started", 400)
		return
	}

	db.Exec("INSERT INTO active_session (id, state) VALUES ($1, $2)", rs, 1)

	fmt.Printf("Session created with id: %s\n", rs)

	clientOptions := mqtt.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	client := mqtt.NewClient(clientOptions)

	con := client.Connect().Wait()
	if !con {
		fmt.Fprint(*w, "connection failed to broker")
		return
	}

	topic := fmt.Sprintf("tes_deh/benar/%s", rs)

	go func() {
		sub := client.Subscribe(topic, 1, roastCb(db, &rs)).Wait()
		if !sub {
			fmt.Fprint(*w, "can not subscribe")
			return
		}

		for client.IsConnected() {
			time.Sleep(time.Millisecond)
		}
	}()

	fmt.Fprint(*w, "Roasting session created")
}

func roastCb(db *sql.DB, rs *string) mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		msg := fmt.Sprintf("%s", m.Payload())
		if msg == "-1" {
			_, e := db.Exec("DELETE FROM active_session WHERE id = $1", *rs)

			if e != nil {
				fmt.Printf("Error occured: %s", e.Error())
			}

			c.Disconnect(1000)
			return
		}
		fmt.Printf("Received message: %s from topic: %s\n", m.Payload(), m.Topic())
	}
}

type Session struct {
	Id    int `json:"id"`
	State int `json:"state"`
}
