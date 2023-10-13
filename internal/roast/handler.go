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
	var id int
	row := db.QueryRow("SELECT id FROM active_session WHERE id = $1", rs)

	row.Scan(&id)
	fmt.Printf("Session id: %d", id)

	if fmt.Sprintf("%d", id) == rs {
		http.Error(*w, "session already started", 400)
		return
	}

	clientOptions := mqtt.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	client := mqtt.NewClient(clientOptions)

	con := client.Connect().Wait()
	if !con {
		fmt.Fprint(*w, "connection failed to broker")
		return
	}

	topic := fmt.Sprintf("tes_deh/benar/%s", rs)

	go func() {
		db.Exec("INSERT INTO active_session (id, state) VALUES ($1, $2)", rs, 1)
		sub := client.Subscribe(topic, 1, roastCallback).Wait()
		if !sub {
			fmt.Fprint(*w, "can not Subscribe")
			return
		}

		for client.IsConnected() {
			time.Sleep(time.Millisecond)
		}
	}()

	fmt.Fprint(*w, "Roasting session created")
}

var roastCallback mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	msg := fmt.Sprintf("%s", m.Payload())
	if msg == "-1" {
		c.Disconnect(1000)
		return
	}
	fmt.Printf("Received message: %s from topic: %s\n", m.Payload(), m.Topic())
}

type Session struct {
	Id    int `json:"id"`
	State int `json:"state"`
}
