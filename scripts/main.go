package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	cO := mqtt.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	client := mqtt.NewClient(cO)

	if !client.Connect().WaitTimeout(time.Second * 20) {
		log.Fatal("Broker timeout connection")
		return
	}

	defer client.Disconnect(1000)

	var topics []string

	for i := 1; i <= 10; i++ {
		topic := fmt.Sprintf("tes_deh/benar/%d", i)
		topics = append(topics, topic)
	}

	wg := &sync.WaitGroup{}
	for _, topic := range topics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()

			// Define the range of values
			minValue := 80.0
			maxValue := 130.0

			// Define the number of data points to generate
			numDataPoints := 100

			// Initialize variables for the mean and the increment
			mean := (minValue + maxValue) / 2.0
			increment := (maxValue - minValue) / float64(numDataPoints)

			for i := 0; i < numDataPoints; i++ {
				f := rand.Float64()*(maxValue-minValue) + minValue
				mean += increment
				f = (f + mean) / 2.0
				smthng := client.Publish(topic, 1, false, fmt.Sprintf("%f", f))
				if err := smthng.Error(); err != nil {
					log.Fatalln(err.Error())
				}
				log.Printf("Published to topic: %s, payload: %f\n", topic, f)

				time.Sleep(time.Millisecond * 100)
			}

			client.Publish(topic, 2, false, "-1")
		}(topic)
	}
	wg.Wait()

	log.Println("TEST DONE")
}
