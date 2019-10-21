package GoGameServer

import (
	"github.com/alehano/wsgame/game"
	"github.com/gorilla/websocket"
)

type playerConn struct {
	ws *websocket.Conn
	*game.Player
	room *room
}

// получать сообщения от ws в goroutine
func (pc *playerConn) receiver() {
	for {
		_, command, err := pc.ws.ReadMessage()
		if err != nil {
			break
		}
		// выполнить команду
		pc.Command(string(command))
		// обновить все conn
		pc.room.updateAll <- true
	}
	pc.room.leave <- pc
	pc.ws.Close()
}

func (pc *playerConn) sendState() {
	go func() {
		msg := pc.GetState()
		err := pc.ws.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			pc.room.leave <- pc
			pc.ws.Close()
		}
	}()
}

func NewPlayerConn(ws *websocket.Conn, player *game.Player, room *room) *playerConn {
	pc := &playerConn{ws, player, room}
	go pc.receiver()
	return pc
}