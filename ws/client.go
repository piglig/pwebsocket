package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"time"
)

type Client struct {
	*websocket.Conn
	UserID uint64
	RoomID uint64
}

func (c *Client) Heartbeat(manager *Manager) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx, _ := context.WithTimeout(context.Background(), 8*time.Second)
				err := c.Ping(ctx)
				if err != nil {
					fmt.Println("remove the client", c, time.Now(), err)
					c.Close(websocket.StatusAbnormalClosure, "heartbeat fail")
					manager.RemoveConn(c)
					return
				}
			}
		}
	}()
}

func (c *Client) Do(event chan Event) {
	for {
		typ, r, err := c.Read(context.Background())
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("client received type:[%s] msg:[%s]\n", typ, string(r))
		var e Event
		err = json.Unmarshal(r, &e)
		if err != nil {
			log.Println(err)
			continue
		}
		e.Client = c
		event <- e
	}
}
