package roast

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewRoastSession(
	rs string,
	roasterId string,
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

	clientOptions := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER"))
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

	log.Printf("Active session id: %s\n", rs)

	// TODO: roaster_id acquired from Request
	// WARN: roast_session IS DIFFERENT from active_session, hence different ids

	// store roast_session.id
	var rsId int
	err = db.QueryRow(
		`INSERT INTO roast_sessions (roaster_id, user_id, roast_date) values ($1, $2, $3) RETURNING id`,
		roasterId,
		rs,
		time.Now(),
	).Scan(&rsId)

	if err != nil {
		db.Exec(`DELETE FROM active_session WHERE id = $1`, rs)
		log.Printf("Active session id: %s aborted", rs)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	log.Printf("Recorded session id: %d\n", rsId)

	stmt, err := db.Prepare(
		`INSERT INTO session_measurements (session_id, suhu, timestamp) VALUES ($1, $2, $3)`,
	)
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

	fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf(`{"status": %s}`, "200"))
	w.(http.Flusher).Flush()

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

	wg.Wait()
	fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf(`{"suhu": %s}`, "-404.404"))
	w.(http.Flusher).Flush()
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
		msg := string(m.Payload())
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
		measurement := MeasurementSession{}
		err := json.Unmarshal(m.Payload(), &measurement)
		if err != nil {
			log.Panicln(err.Error())
			db.Exec("DELETE FROM active_session WHERE id = $1", *rs)
		}

		ts, _ := measurement.Timestamp.MarshalText()

		fmt.Fprintf(
			w,
			"data: %s\n\n",
			fmt.Sprintf(`{"suhu": %f, "ts": "%s"}`, measurement.Suhu, ts),
		)
		w.(http.Flusher).Flush()
		// TODO: Probably use async (goroutine)
		_, err = stmt.Exec(*rsId, measurement.Suhu, measurement.Timestamp)
		if err != nil {
			fmt.Fprintf(w, "data: %s\n\n", err.Error())
			w.(http.Flusher).Flush()
		}
	}
}

func GetRoastSessions(db *sql.DB, userId string) (string, error) {
	stmt := "SELECT rs.id, roaster_id, user_id, roast_date FROM roast_sessions rs LEFT JOIN users u ON rs.user_id = u.id WHERE u.id = $1"

	rows, err := db.Query(stmt, userId)
	if err != nil {
		log.Printf("user_id: %s", userId)
		log.Fatalf("WARN => query failed: %s", err.Error())
	}

	defer rows.Close()

	sessions := []RoastSession{}

	for rows.Next() {
		session := RoastSession{}

		if err = rows.Scan(&session.Id, &session.RoasterId, &session.UserId, &session.Timestamp); err != nil {
			return "", err
		}
		sessions = append(sessions, session)
	}

	if len(sessions) == 0 {
		return "", errors.New("No sessions found")
	}

	if err = rows.Err(); err != nil {
		return "", err
	}
	data, err := json.Marshal(sessions)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(data), nil
}

func GetMeasurements(db *sql.DB, sessionId string) (string, error) {
	stmt := "SELECT session_id, suhu, timestamp FROM session_measurements WHERE session_id = $1"

	rows, err := db.Query(stmt, sessionId)
	if err != nil {
		log.Printf("WARN => query failed: %s", err.Error())
		return "", err
	}

	defer rows.Close()

	measurements := []MeasurementSession{}

	for rows.Next() {
		measurement := MeasurementSession{}

		if err = rows.Scan(&measurement.SessionId, &measurement.Suhu, &measurement.Timestamp); err != nil {
			return "", err
		}
		measurements = append(measurements, measurement)
	}

	if len(measurements) == 0 {
		return "", errors.New("No measurement found")
	}

	if err = rows.Err(); err != nil {
		return "", err
	}
	data, err := json.Marshal(measurements)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(data), nil
}

type Session struct {
	Id    int `json:"id"`
	State int `json:"state"`
}

type SubscriberWait struct {
	SessionId int
	RoastDone bool
}

type RoastSession struct {
	Id        uint16    `json:"id"`
	RoasterId uint16    `json:"roaster_id"`
	UserId    uint16    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type MeasurementSession struct {
	SessionId int       `json:"session_id"`
	Suhu      float64   `json:"suhu"`
	Timestamp time.Time `json:"timestamp"`
}
