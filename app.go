package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

type croney struct {
	stopCh    chan bool
	zookeeper []string

	zkConn *zk.Conn
	zkCh   <-chan zk.Event
}

func newApp(zookeeper []string) (*croney, error) {
	stopper := make(chan bool)
	return &croney{
		zookeeper: zookeeper,
		stopCh:    stopper,
	}, nil
}

func handleZkEvent(event zk.Event) {
	switch event.State {
	case zk.StateConnecting:
		log.Println("attempting to connect...")
	case zk.StateConnected:
		log.Println("connected")
	case zk.StateHasSession:
		log.Println("connected with session")
	}
}

func (c *croney) run() {
	conn, ch, err := zk.Connect(c.zookeeper, 5*time.Second)
	if err != nil {
		panic(err)
	}
	c.zkConn = conn
	c.zkCh = ch

	log.Println("starting")
	go func() {
		for {
			m := <-c.zkCh
			handleZkEvent(m)
		}
	}()

	<-c.stopCh
}

func (c *croney) stop() {
	log.Println("closing connection")
	c.zkConn.Close()
	c.stopCh <- true
}