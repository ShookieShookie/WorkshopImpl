package main

import (
	"bufio"
	"fmt"
	"github.com/ShookieShookie/WorkshopImpl/deck"
	"github.com/ShookieShookie/WorkshopImpl/game"
	"github.com/ShookieShookie/WorkshopImpl/hand"
	"github.com/ShookieShookie/WorkshopImpl/player"
	"math/rand"
	"os"
	"time"
)

var originalDeck = []int{0, 0, 1, 1, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 5, 5, 6, 6, 7, 8}

func main() {
	rand.Seed(int64(time.Now().Second()))
	h1 := hand.NewHand()
	d1 := deck.NewDeck(rand.Intn)
	for i := 0; i < 20; i++ {
		d1.Add(originalDeck[i])
	}
	p1 := player.NewPlayer("player1", 30, 0, h1, d1)
	h2 := hand.NewHand()
	d2 := deck.NewDeck(rand.Intn)
	for i := 0; i < 20; i++ {
		d2.Add(originalDeck[i])
	}
	p2 := player.NewPlayer("player2", 30, 0, h2, d2)

	g := game.NewGame(p1, p2, getCard, game.Turn)
	g.Start()
}

// requires an integration test
func getCard() string {
	fmt.Printf("Enter card index to play (-1 to end turn): ")
	defer fmt.Printf("\n")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
