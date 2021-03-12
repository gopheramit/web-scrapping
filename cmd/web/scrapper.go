package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"time"

	"github.com/gopheramit/distributed-go-with-rabbitmq/src/distributed/dto"
	"github.com/gopheramit/distributed-go-with-rabbitmq/src/distributed/qutils"
	"github.com/streadway/amqp"
)

var url = "amqp://guest:guest@localhost:5672"
var name = flag.String("name", "sensor", "name of the sensor")

func main1(url1 string, js1 string) {
	flag.Parse()
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()
	dataQueue := qutils.GetQueue(*name, ch)
	publishQueueName(ch)
	discoveryQueue := qutils.GetQueue("", ch)
	ch.QueueBind(
		discoveryQueue.Name,            //name string,
		"",                             //key string,
		qutils.SensorDiscoveryExchange, //exchange string,
		false,                          //noWait bool,
		nil)                            //args amqp.Table)

	go listenForDiscoverRequests(discoveryQueue.Name, ch)
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	reading := dto.SensorMessage{
		Name:      *name,
		Url:       url1,
		Js:        js1,
		Timestamp: time.Now(),
	}
	buf.Reset()
	enc = gob.NewEncoder(buf)
	enc.Encode(reading)
	msg := amqp.Publishing{
		Body: buf.Bytes(),
	}
	ch.Publish(
		"",             //exchange string,
		dataQueue.Name, //key string,
		false,          //mandatory bool,
		false,          //immediate bool,
		msg)            //msg amqp.Publishing)

	log.Printf("Reading sent. Value: %v\n", msg)
}
func listenForDiscoverRequests(name string, ch *amqp.Channel) {
	msgs, _ := ch.Consume(
		name,  //queue string,
		"",    //consumer string,
		true,  //autoAck bool,
		false, //exclusive bool,
		false, //noLocal bool,
		false, //noWait bool,
		nil)   //args amqp.Table)

	for range msgs {
		log.Println("received discovery request")
		publishQueueName(ch)
	}
}

func publishQueueName(ch *amqp.Channel) {
	msg := amqp.Publishing{Body: []byte(*name)}
	ch.Publish(
		"amq.fanout", //exchange string,
		"",           //key string,
		false,        //mandatory bool,
		false,        //immediate bool,
		msg)          //msg amqp.Publishing)
}
