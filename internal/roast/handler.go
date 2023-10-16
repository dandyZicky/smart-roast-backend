package roast

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	row := db.QueryRow(`SELECT id, state FROM active_session WHERE id = $1`, rs)

	err := row.Scan(&sessionState.Id, &sessionState.State)

	if (fmt.Sprintf("%d", sessionState.Id) == rs) && (sessionState.State == 1) {
		http.Error(w, "session already started", 400)
		return
	}

	clientOptions := mqtt.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	client := mqtt.NewClient(clientOptions)

	if !client.Connect().WaitTimeout(time.Second * 20) {
		fmt.Fprint(w, "Took to long to connect to broker")
		return
	}

	defer client.Disconnect(1000)

	_, err = db.Exec(`INSERT INTO active_session (id, state) VALUES ($1, $2)`, rs, 1)
	if err != nil {
		http.Error(w, "cannot create roasting session", http.StatusBadRequest)
		return
	}

	// TODO: roaster_id acquired from Request
	// WARN: roast_session IS DIFFERENT from active_session, hence different ids

	// store roast_session.id
	var rsId int
	err = db.QueryRow(
		`INSERT INTO roast_sessions (roaster_id, user_id, roast_date) values ($1, $2, $3) RETURNING id`,
		1,
		1,
		time.Now(),
	).Scan(&rsId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	log.Printf("Active session id: %s\n", rs)
	log.Printf("Recorded session id: %d\n", rsId)

	stmt, err := db.Prepare(`INSERT INTO session_measurements (session_id, suhu) VALUES ($1, $2)`)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
	}

	mqttWait := make(chan SubscriberWait, 1)
	topic := fmt.Sprintf("tes_deh/benar/%s", rs)

	sub := client.Subscribe(topic, 1, roastCb(db, stmt, &rs, &rsId, w, mqttWait)).Wait()
	if !sub {
		http.Error(w, "can not subscribe", 400)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case msg := <-mqttWait:
			if msg.RoastDone {
				log.Printf("ROAST: %d, DONE", msg.SessionId)
				wg.Done()
			}
		}
	}()

	fmt.Fprint(w, "{message: `initialized connection`}")
	w.(http.Flusher).Flush()
	wg.Wait()
	fmt.Fprintln(w, "Session complete")
	db.Exec(`DELETE FROM active_session WHERE id = $1`, rs)
	return
}

func roastCb(
	db *sql.DB,
	stmt *sql.Stmt,
	rs *string,
	rsId *int,
	w http.ResponseWriter,
	state chan SubscriberWait,
) mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		msg := fmt.Sprintf("%s", m.Payload())
		if msg == "-1" {
			_, e := db.Exec(`DELETE FROM active_session WHERE id = $1`, *rs)
			state <- SubscriberWait{
				SessionId: *rsId,
				RoastDone: true,
			}

			if e != nil {
				log.Printf("Error occured: %s", e.Error())
			}

			c.Disconnect(1000)
			return
		}
		s := string(m.Payload())
		suhu, _ := strconv.ParseFloat(s, 64)
		fmt.Fprintf(w, `{"suhu": %f}`, suhu)
		w.(http.Flusher).Flush()
		// TODO: Probably use async (goroutine)
		_, err := stmt.Exec(*rsId, suhu)
		if err != nil {
			fmt.Fprint(w, err.Error())
			w.(http.Flusher).Flush()
		}
	}
}

type Session struct {
	Id    int `json:"id"`
	State int `json:"state"`
}

type SubscriberWait struct {
	SessionId int
	RoastDone bool
}
