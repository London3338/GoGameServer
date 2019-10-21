package GoGameServer

import (
	"github.com/alehano/wsgame/game"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

const (
	ADDR string = ":8080"
)

func homeHandler(c http.ResponseWriter, r *http.Request) {
	var homeTempl = template.Must(template.ParseFiles("templates/home.html"))
	data := struct {
		Host       string
		RoomsCount int
	}{r.Host, roomsCount}
	homeTempl.Execute(c, data)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	playerName := "Player"
	params, _ := url.ParseQuery(r.URL.RawQuery)
	if len(params["name"]) > 0 {
		playerName = params["name"][0]
	}

	// получаем или создаем комнату
	var room *room
	if len(freeRooms) > 0 {
		for _, r := range freeRooms {
			room = r
			break
		}
	} else {
		room = NewRoom("")
	}

	// создать игрока
	player := game.NewPlayer(playerName)
	pConn := NewPlayerConn(ws, player, room)
	// присоединяем игрока к комнате
	room.join <- pConn

	log.Printf("Player: %s has joined to room: %s", pConn.Name, room.name)
}

func main() {
	// регистрируем два обработчика
	// homeHandler -выводит шаблон home.html
	// wsHandler -устанавливает WebSocket соединение и регистрирует игрока
	// WebSocket -используем пакет из набора Gorilla Toolkit «github.com/gorilla/websocket»
	http.HandleFunc("/", homeHandler)
	// создаем новое соединение ws
	http.HandleFunc("/ws", wsHandler)

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	if err := http.ListenAndServe(ADDR, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}