package main

import (
	"net/url"
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"fmt"
	"time"
)

type PeriscopeChatListener struct {
	c *websocket.Conn
	pm PeriscopeMeta
	cm PeriscopeMetaChat
	done chan struct{}
}

func NewPeriscopeChatListener(pm PeriscopeMeta, cm PeriscopeMetaChat) *PeriscopeChatListener {
	return &PeriscopeChatListener{pm: pm, cm: cm}
}

func (l *PeriscopeChatListener) Dial() error {
	if l.c != nil {
		l.c.Close()
	}
	u, err := url.Parse(l.cm.Endpoint)
	if err != nil {
		return err
	}
	switch u.Scheme {
		case "http":
			u.Scheme = "ws"
		case "https":
			u.Scheme = "wss"
	}
	l.c, _, err = websocket.DefaultDialer.Dial(u.String() + "/chatapi/v1/chatnow", nil)
	if err != nil {
		return err
	}
	return nil
}

type frameMessage struct {
	Kind int `json:"kind"`
	Payload string `json:"payload,omitempty"`
	Body string `json:"body,omitempty"`
}

func (fm *frameMessage) String() string {
	return bytes.NewBuffer(toJsonOrPanic(fm)).String()
}

func (l *PeriscopeChatListener) initChat() error {
	dct := map[string]string{"access_token": l.cm.AccessToken}
	msg := frameMessage{Kind: 3, Payload: bytes.NewBuffer(toJsonOrPanic(dct)).String()}
	err := l.c.WriteMessage(websocket.TextMessage, toJsonOrPanic(msg))
	if err != nil {
		log.Println(toJsonOrPanic(msg))
		return err
	}
	dct = map[string]string{"room": l.cm.RoomId}
	msg = frameMessage{Kind: 1, Body: bytes.NewBuffer(toJsonOrPanic(dct)).String()}
	msg = frameMessage{Kind: 2, Payload: msg.String()}
	err = l.c.WriteMessage(websocket.TextMessage, toJsonOrPanic(msg))
	if err != nil {
		return err
	}
	return nil
}

type RetryTimeout struct{}

func (RetryTimeout) Error() string {
	return "closed connection"
}

func (l *PeriscopeChatListener) Run() error {
	var done bool
	l.done = make(chan struct{})
	go func() {
		select {
			case <-l.done:
				// TODO Check on this
				// To cleanly close a connection, a client should send a close
				// frame and wait for the server to close the connection.
				done = true
				err := l.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
				}
				l.c.Close()
		}
	}()
	defer l.Stop()

	fail_counter := 0
	for {
		if fail_counter == 3 { // Retry 3 times without waiting
			return RetryTimeout{}
		}
		if err := l.Dial(); err != nil {
			return err
		}
		if err := l.initChat(); err != nil {
			return err
		}
		for {
			_, msg, err := l.c.ReadMessage()
			if err != nil {
				if done {
					return err
				}
				//if websocket.IsCloseError()//TODO CHECK DOC
				log.Println("read:", err)
				fail_counter++
				break
			}
			fail_counter = 0
			pm, err := FrameFilter(msg)
			if err != nil {
				log.Println("FrameFilter:", err)
				continue
			}
			if pm.Type == 1 {
				if pm.Timestamp > 100000000000 { //WTF
					pm.Timestamp /= 1000
				}
				tm := time.Unix(int64(pm.Timestamp), 0)
				fmt.Printf("%s %s: %s\n",tm, pm.DisplayName, pm.Body)
			}
		}
	}
}

func (l *PeriscopeChatListener) Stop() {
	close(l.done)
}
