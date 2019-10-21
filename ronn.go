package GoGameServer

import (
	"github.com/alehano/wsgame/game"
	"github.com/alehano/wsgame/utils"
	"log"
)

// список всех созданных комнат
var allRooms = make(map[string]*room)
// комнаты с одним игроком
var freeRooms = make(map[string]*room)
// счетчик работающих комнат
var roomsCount int

// room -содержит имя комнаты
type room struct {
	name string

	//playerConns — список подключенных соединений (игроков)
	//и несколько каналов для управления
	// зарегистрированные соединения
	playerConns map[*playerConn]bool

	// обновление состояния
	// извещения нужно ли обновлять состояние игры
	updateAll chan bool

	// регистрация запросов от подключений
	// передается указатель на конкретное соединение
	join chan *playerConn

	// отмена регистрации запросов от подключений
	leave chan *playerConn
}

// запустить комнату в рутине
func (r *room) run() {
	for {
		select {
		case c := <-r.join:
			r.playerConns[c] = true
			r.updateAllPlayers()

			// если комната заполнена удаляет из freeRooms
			if len(r.playerConns) == 2 {
				delete(freeRooms, r.name)
				// пара игроков
				var p []*game.Player
				for k, _ := range r.playerConns {
					p = append(p, k.Player)
				}
				game.PairPlayers(p[0], p[1])
			}

		case c := <-r.leave:
			c.GiveUp()
			r.updateAllPlayers()
			delete(r.playerConns, c)
			if len(r.playerConns) == 0 {
				goto Exit
			}
		case <-r.updateAll:
			r.updateAllPlayers()
		}
	}

Exit:

	// удалить комнату
	delete(allRooms, r.name)
	delete(freeRooms, r.name)
	roomsCount -= 1
	log.Print("Room closed:", r.name)
}

func (r *room) updateAllPlayers() {
	for c := range r.playerConns {
		c.sendState()
	}
}

// создание новой комнаты
func NewRoom(name string) *room {
	if name == "" {
		name = utils.RandString(16)
	}

	room := &room{
		name:        name,
		playerConns: make(map[*playerConn]bool),
		updateAll:   make(chan bool),
		join:        make(chan *playerConn),
		leave:       make(chan *playerConn),
	}

	allRooms[name] = room

	freeRooms[name] = room

	go room.run()

	roomsCount += 1

	return room
}