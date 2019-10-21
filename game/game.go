package game

import (
	"log"
)

// игроки
type Player struct {
	Name  string
	Enemy *Player
}

func NewPlayer(name string) *Player {
	player := &Player{Name: name}
	return player
}

// соединение игроков
func PairPlayers(p1 *Player, p2 *Player) {
	p1.Enemy, p2.Enemy = p2, p1
}

func (p *Player) Command(command string) {

	log.Print("Command: '", command, "' received by player: ", p.Name)
}

// получить текущее состояние игры для данного игрока
func (p *Player) GetState() string {
	return "Game state for Player: " + p.Name
}

// сдаться и присвоить победу противнику
func (p *Player) GiveUp() {
	log.Print("Player gave up: ", p.Name)
}
