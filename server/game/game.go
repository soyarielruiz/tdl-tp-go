package game

import (
	"log"

	"github.com/soyarielruiz/tdl-borbotones-go/server/action"
	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

type Game struct {
	User_chan chan user.User
	Users     []user.User
}

func Start(uc chan user.User, game_number int) {
	log.Printf("Initializing game number: %d\n", game_number)
	g := Game{User_chan: uc, Users: make([]user.User, 0)}
	RecvUsers(&g, game_number)
	i := 0
	for {
		log.Printf("Waiting for action from %s in game %d\n", g.Users[i].PlayerId, game_number)
		action := RecvMsg(&g, &g.Users[i])
		SendMsg(&g, &action)
		if action.Command == "exit" {
			log.Printf("Exit command received in game %d\n", game_number)
			return
		}
		i++
		if i == 3 {
			i = 0
		}
	}
}

func RecvUsers(g *Game, number int) {
	for {
		u := <-g.User_chan
		log.Printf("New usr received in game %d. %s", number, u)
		g.Users = append(g.Users, u)
		if len(g.Users) == 3 {
			log.Printf("3 users connect to game %d. Starting game", number)
			return
		} else {
			log.Printf("No enough users connected to game %d for start the game", number)
		}
	}
}

func SendMsg(g *Game, a *action.Action) {
	for _, u := range g.Users {
		u.SendChannel <- *a
	}
}

func RecvMsg(g *Game, u *user.User) action.Action {
	return <-u.ReceiveChannel
}
