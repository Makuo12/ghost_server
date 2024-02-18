package api

import (
	//"encoding/json"
	//db "flex_server/db/sqlc"
	"flex_server/val"
	"log"
	"sync"

	//"log"

	"github.com/google/uuid"
)

type message struct {
	Data     []byte    `json:"data"`
	Room     string    `json:"Room"`
	Username uuid.UUID `json:"Username"`
	UserID   uuid.UUID `json:"user_id"`
	// Type string `json:"type"`
}

type subscription struct {
	conn *connection
	Room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription

	//server *Server

	broadcastPersonal chan message
}

var H = hub{
	broadcast:         make(chan message),
	broadcastPersonal: make(chan message),
	register:          make(chan subscription),
	unregister:        make(chan subscription),
	rooms:             make(map[string]map[*connection]bool),
}

func (h *hub) Run() {
	for {
		select {
		//first we listen for a register message using channel
		case s := <-h.register:
			//if register message is received,
			connections := h.rooms[s.Room]
			if connections == nil {
				//remember that a if a Room exist then it should have a connection/
				//because this Room has no connections we then just add the connection to the Room

				conns := make(map[*connection]bool)
				h.rooms[s.Room] = conns
			}
			func() {
				s.conn.mu.Lock()
				defer s.conn.mu.Unlock()
				h.rooms[s.Room][s.conn] = true
			}()

			log.Println(s.Room, connections)
		//case s := <-h.unregister:
		//	connections := h.rooms[s.Room]
		//	if connections != nil {
		//		if _, ok := connections[s.conn]; ok {
		//			delete(connections, s.conn)
		//			close(s.conn.send)
		//		}
		//		if len(connections) == 0 {
		//			delete(h.rooms, s.Room)
		//		}
		//	}
		case s := <-h.unregister:
			connections := h.rooms[s.Room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					// Use sync.Once to ensure the channel is closed only once
					var closeOnce sync.Once

					// Safely close the channel using sync.Once
					closeFunc := func() {
						closeOnce.Do(func() {
							defer func() {
								if r := recover(); r != nil {
									// Handle the panic, if any
									log.Println("Panic occurred during channel closing:", r)
								}
							}()

							s.conn.mu.Lock()
							defer s.conn.mu.Unlock()
							close(s.conn.send)
						})
					}
					// Call closeFunc to close the channel
					closeFunc()

					// Delete the connection
					delete(connections, s.conn)

					// Check if all connections in the room are closed and delete the room
					if len(connections) == 0 {
						delete(h.rooms, s.Room)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.Room]
			for c := range connections {
				if c.roomId == "search_option" {
					if c.username == m.Username {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}
				}
				if c.roomId == "search_edt_option" {
					if c.username == m.Username {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}
				}
				// cal is for the CalenderView
				if c.roomId == "search_cal_option" {
					if c.username == m.Username {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}
				}
				if c.roomId == "user_search_event_name" {
					if c.username == m.Username {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}
				}
				if c.roomId == "map_experience_location" {
					if c.username == m.Username {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}
				}
				if c.roomId == "ex_search_event" {
					if c.username == m.Username {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}
				}
				if c.roomId == "message_listen" {
					if c.userID == m.UserID {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}

				}
				if c.roomId == "notification_listen" {
					if c.userID == m.UserID {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}

				}
				if c.roomId == "message_unread" {
					if c.userID == m.UserID {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}

				}

				if c.roomId == "reserve" {
					if c.userID == m.UserID {
						select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}
					}
				}
				if val.ValidateMsgRoom(c.roomId) {
					log.Println("saw message")
					select {
						case c.send <- m.Data:
						default:
							func() {
								defer func() {
									if r := recover(); r != nil {
										// Handle the panic, if any
										log.Println("Panic occurred during send operation:", r)
									}
								}()

								close(c.send)
								delete(connections, c)
								if len(connections) == 0 {
									delete(h.rooms, m.Room)
								}
							}()
						}

				}
			}
		}
	}
}

//ws://0.0.0.0:8080/
