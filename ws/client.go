package ws

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"time"
)

type Client struct {
	*websocket.Conn
	RoomID string
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
					fmt.Println("remove the client", c, time.Now())
					c.Close(websocket.StatusAbnormalClosure, "heartbeat fail")
					manager.RemoveConn(c)
					return
				}
			}
		}
	}()
}
