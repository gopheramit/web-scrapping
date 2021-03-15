package service

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"

	"github.com/gopheramit/web-scrapping/cmd/web/dto"
	"github.com/gopheramit/web-scrapping/cmd/web/qutils"
	"github.com/streadway/amqp"
)

const url = "amqp://guest:guest@localhost:5672"

type QueueListener struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	sources      map[string]<-chan amqp.Delivery
	ScrapRequest *models.ScrapRequestModel
	//ea      *EventAggregatorw
}

func NewQueueListener(db *sql.DB) *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
		//ea:      NewEventAggregator(),
		ScrapRequest: &models.ScrapRequestModel{DB: db},
	}

	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

func (ql *QueueListener) DiscoverSensors() {
	ql.ch.ExchangeDeclare(
		qutils.SensorDiscoveryExchange, //name string,
		"fanout",                       //kind string,
		false,                          //durable bool,
		false,                          //autoDelete bool,
		false,                          //internal bool,
		false,                          //noWait bool,
		nil)                            //args amqp.Table)

	ql.ch.Publish(
		qutils.SensorDiscoveryExchange, //exchange string,
		"",                             //key string,
		false,                          //mandatory bool,
		false,                          //immediate bool,
		amqp.Publishing{})              //msg amqp.Publishing)
}

func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch)
	ql.ch.QueueBind(
		q.Name,       //name string,
		"",           //key string,
		"amq.fanout", //exchange string,
		false,        //noWait bool,
		nil)          //args amqp.Table)

	msgs, _ := ql.ch.Consume(
		q.Name, //queue string,
		"",     //consumer string,
		true,   //autoAck bool,
		false,  //exclusive bool,
		false,  //noLocal bool,
		false,  //noWait bool,
		nil)    //args amqp.Table)

	ql.DiscoverSensors()

	fmt.Println("listening for new sources")

	// updated the if guard below to surround all
	// of the for-loops contents to prevent
	// same sensor being registered multiple
	// times with RabbitMQ
	for msg := range msgs {
		if ql.sources[string(msg.Body)] == nil {
			fmt.Println("new source discovered")
			sourceChan, _ := ql.ch.Consume(
				string(msg.Body), //queue string,
				"",               //consumer string,
				true,             //autoAck bool,
				false,            //exclusive bool,
				false,            //noLocal bool,
				false,            //noWait bool,
				nil)              //args amqp.Table)

			ql.sources[string(msg.Body)] = sourceChan

			go ql.AddListener(sourceChan)
		}
	}
}

func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)
		fmt.Println(sd)
		fmt.Printf("Received message: %v\n", sd)
		boolean := Negation(sd.Js)
		fmt.Println("Js :")
		fmt.Println(boolean)

		ql.linkscrape(sd.Url, sd.Key)

		/*
			res, err := http.Get(sd.Url)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(doc.Html())

			key1 := genUlid()
			keystr1 := key1.String()
			fmt.Println(keystr1)
			resullt, err := doc.Html()

			stmt := `INSERT INTO ScrapRequest (uuid,guid,BLOBData)VALUES(?,?,?)`
			_, err = db.Exec(stmt, keystr1, keystr1, []byte(resullt))
		*/
		//err = app1.ScrapRequest.Insert(keystr1, keystr1, []byte(resullt))
		//if err != nil {
		//	fmt.Println("error linkscrape")
		//
		//		} else {
		//			fmt.Println("everthing ok")
		//		}
		//app1.linkscrape(sd.Url) // sd.Key)
		//ed := EventData{
		//	Name:      sd.Name,
		//		Timestamp: sd.Timestamp,
		//		Js:        sd.Js,
	}

	//	ql.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)
	//}
}

func Negation(boolean bool) bool {
	//time.Sleep(10)
	if boolean == true {
		return false
	} else {
		return true
	}
}

/*
func (app1 *application1) linkscrape(url string) { // key string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(doc.Html())

	key1 := genUlid()
	keystr1 := key1.String()
	fmt.Println(keystr1)
	resullt, err := doc.Html()
	err = app1.ScrapRequest.Insert(keystr1, keystr1, []byte(resullt))
	if err != nil {
		fmt.Println("error linkscrape")

	} else {
		fmt.Println("everthing ok")
	}
}
*/
