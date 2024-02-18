package api

import (
	//"encoding/json"
	"errors"
	db "flex_server/db/sqlc"
	"flex_server/val"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins
		return true
	},
}
var c *connection

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	//roomId
	roomId string

	//username
	username uuid.UUID

	// UserID
	userID uuid.UUID

	// from what app whether user or host
	app string

	// server
	server *Server

	// ctx
	ctx *gin.Context

	// User
	user db.User

	mu sync.Mutex // guards
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	ctx := s.conn
	defer func() {
		H.unregister <- s
		err := ctx.ws.Close()
		if err != nil {
			log.Printf("Websocket for readPump ctx.ws.Close err:%v\n", err.Error())
			return
		}
	}()
	ctx.ws.SetReadLimit(maxMessageSize)
	err := ctx.ws.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("Websocket for readPump ctx.ws.SetReadDeadline err:%v\n", err.Error())
		return
	}
	ctx.ws.SetPongHandler(func(string) error {
		err = ctx.ws.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("Websocket for readPump SetPongHandler ctx.ws.SetReadDeadline err:%v\n", err.Error())
		}
		return nil
	})
	for {
		_, msg, err := ctx.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		//log.Println("msg", msg)
		m := message{msg, s.conn.roomId, s.conn.username, c.userID}
		H.broadcast <- m
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	ctx := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := ctx.ws.Close()
		if err != nil {
			log.Printf("Websocket for writePump ctx.ws.Close err:%v\n", err.Error())
			return
		}
	}()
	for {
		select {
		case msg, ok := <-ctx.send:
			// fmt.Println("ok:", ok)
			// fmt.Println("msg:", msg)
			if !ok {
				err := ctx.write(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("Websocket for writePump ctx.write(websocket.CloseMessage err:%v\n", err.Error())
				}
				return
			}
			if err := ctx.write(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			// fmt.Println("case ticker")
			if err := ctx.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (ctx *connection) write(mt int, payload []byte) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			// Handle the panic recovery here, log it, and take appropriate actions.
		}
	}()
	err := ctx.ws.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		log.Printf("Websocket write err:%v\n", err.Error())
	}
	if mt == websocket.PingMessage {
		//log.Println("just a ping packet")
	} else if mt == websocket.TextMessage {
		if len(payload) == 0 {
			err := errors.New("no data in payload")
			log.Println(err)
		} else {
			switch ctx.roomId {
			case "search_option":
				data, hasData, err := HandleSearchTextRes(ctx, payload)
				//log.Println("got data", data)
				if err == nil && hasData {
					return ctx.ws.WriteMessage(mt, data)
				}
			case "notification_listen":
				data, hasData, err := HandleNotificationListen(ctx, payload)
				//log.Println("got data", data)
				if err == nil && hasData {
					return ctx.ws.WriteMessage(mt, data)
				}
			case "ex_search_event":
				data, hasData, err := EventExSearchText(ctx, payload)
				//log.Println("got data", data)
				if err == nil && hasData {
					return ctx.ws.WriteMessage(mt, data)
				}
			case "search_edt_option":
				data, hasData, err := HandleEDTSearchTextRes(ctx, payload)
				//log.Println("got data", data)
				if err == nil && hasData {
					return ctx.ws.WriteMessage(mt, data)
				}
			case "search_cal_option":
				data, hasData, err := HandleSearchTextCalRes(ctx, payload)
				//log.Println("got data", data)
				if err == nil && hasData {
					return ctx.ws.WriteMessage(mt, data)
				}
			case "user_search_event_name":
				data, hasData, err := HandleSearchEventByName(ctx, payload)
				//log.Println("got data", data)
				if err == nil && hasData {
					return ctx.ws.WriteMessage(mt, data)
				}
			case "message_listen":
				data, hasData, err := HandleMessageListen(ctx, payload)
				//log.Println("got data", data)
				if err == nil && hasData {
					return ctx.ws.WriteMessage(mt, data)
				}
			case "message_unread":
				data, hasData, err := HandleUnreadMessageList(ctx, payload)
				//log.Println("got data", data)
				if err == nil && hasData {
					return ctx.ws.WriteMessage(mt, data)
				}
			case "get_message":
				data, hasData, err := HandleGetMessage(ctx, payload)
				//log.Println("got data", data)
				if err == nil && hasData {
					return ctx.ws.WriteMessage(mt, data)
				}
			default:
				log.Println("at default")
				if val.ValidateMsgRoom(ctx.roomId) {
					log.Println("val.ValidateMsgRoom(ctx.roomId)", ctx.roomId, val.ValidateMsgRoom(ctx.roomId))
					data, hasData, err := HandleMessage(ctx, payload)
					if err == nil && hasData {
						log.Println("sending data")
						return ctx.ws.WriteMessage(mt, data)
					}
				}

			}
		}
	}
	return ctx.ws.WriteMessage(mt, []byte{})
}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, roomId string, username uuid.UUID, userID uuid.UUID, app string, server *Server, ctx *gin.Context, user db.User) {
	//ws, err := upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c = &connection{send: make(chan []byte, 256), ws: ws, roomId: roomId, username: username, app: app, server: server, ctx: ctx, userID: userID, user: user}
	s := subscription{c, roomId}
	H.register <- s
	go s.writePump()
	go s.readPump()
}

//var upgrader = websocket.Upgrader{
//	CheckOrigin: func(r *http.Request) bool {
//		// Allow all origins
//		return true
//	},
//}
