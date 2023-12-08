package event

import (
	"L0/models"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

type IEvent interface {
	Connect(url string) error
	Publish(subj string, msg *models.EventOrder) error
	Subscribe(subj string, ch chan *models.EventOrder) error
	Close()
}

type NatsEvent struct {
	conn *nats.Conn
}

func NewEvent() IEvent {
	return &NatsEvent{}
}

func (e *NatsEvent) Connect(url string) error {
	conn, err := nats.Connect(url) //  nats.UserInfo("myname", "password")
	//conn, err := nats.Connect("demo.nats.io") //  nats.UserInfo("myname", "password")
	if err != nil {
		return err
	}
	e.conn = conn
	return nil
}

func (e *NatsEvent) Subscribe(subj string, ch chan *models.EventOrder) error {
	_, err := e.conn.Subscribe(subj, func(msg *nats.Msg) {
		log.Println(string(msg.Data))
		var eventMsg models.EventOrder
		if err := json.Unmarshal(msg.Data, &eventMsg); err != nil {
			log.Println(err)
		}
		ch <- &eventMsg
	})
	return err
}

func (e *NatsEvent) Publish(subj string, msg *models.EventOrder) error {
	mBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return e.conn.Publish(subj, mBytes)
}

func (e *NatsEvent) Close() {
	e.conn.Close()
}
