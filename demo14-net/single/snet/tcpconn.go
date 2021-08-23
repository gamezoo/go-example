package snet

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Conn struct {
	IP       string
	Port     uint32
	TCPConn  *net.TCPConn
	MsgChan  chan []byte
	ExitChan chan bool
	Closed   bool
}

func NewConn(IP string, Port uint32) *Conn {
	s := &Conn{
		IP:       IP,
		Port:     Port,
		MsgChan:  make(chan []byte),
		ExitChan: make(chan bool),
	}
	return s
}

func (c *Conn) Start() {
	log.Printf("%s:%d start...\n", c.IP, c.Port)
	go func() {
		addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", c.IP, c.Port))
		if err != nil {
			log.Println("resolve tcp addr err ", err)
			return
		}
		listener, err := net.ListenTCP("tcp4", addr)
		if err != nil {
			log.Println("listen tcp err ", err)
			return
		}
		var connid uint32
		connid = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Println("accept tcp err ", err)
				continue
			}
			c.TCPConn = conn
			go c.StartRead()
			go c.StartWrite()
			connid++
		}
	}()
	select {}
}
func (c *Conn) StartRead() {
	log.Println("read groutine is waiting")
	defer c.Stop()
	defer log.Println("read groutine exit")
	reader := bufio.NewReader(c.TCPConn)
	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println("startread read bytes error ", err)
			break
		}
		len := len(lineBytes)
		line := lineBytes[:len-1]
		log.Println("start read from client ", string(line))
		go c.HandleMsg(line)
	}
}
func (c *Conn) StartWrite() {
	log.Println("write groutine is waiting")
	defer log.Println("write groutine exit")
	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.TCPConn.Write(data); err != nil {
				log.Println("startwrite conn write error ", err)
				return
			}
			log.Println("start write from server ", string(data))
		case <-c.ExitChan:
			return
		}
	}
}
func (c *Conn) HandleMsg(data []byte) {
	res := fmt.Sprintf("res:%s", string(data))

	c.MsgChan <- []byte(res)
}
func (c *Conn) Stop() {
	if c.Closed {
		return
	}
	c.Closed = true
	c.ExitChan <- true

	c.TCPConn.Close()
	close(c.ExitChan)
	close(c.MsgChan)
}
