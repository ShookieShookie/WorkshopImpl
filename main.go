package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var originalDeck = []int{0, 0, 1, 1, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 5, 5, 6, 6, 7, 8}

type Player interface {
	ApplyDamage(int)
	GetHealth() int
	SetMana(int)
	PlayCard(index int) (int, error)
	IsDead() bool
	Draw() error
	ID() string
	PrintStats()
}

func main() {
	rand.Seed(int64(time.Now().Second()))
	h1 := HandImpl{cards: []int{}}
	d1 := DeckImpl{cards: []int{}}
	d1.getIndex = rand.Intn
	for i := 0; i < 20; i++ {
		d1.Add(originalDeck[i])
	}
	p1 := PlayerImpl{health: 30, name: "player 1", hand: &h1, deck: &d1}
	h2 := HandImpl{cards: []int{}}
	d2 := DeckImpl{cards: []int{}}
	d2.getIndex = rand.Intn
	for i := 0; i < 20; i++ {
		d2.Add(originalDeck[i])
	}
	p2 := PlayerImpl{health: 30, name: "player 2", hand: &h2, deck: &d2}

	g := Game{p1: &p1, p2: &p2, userInput: getCard, turn: turn}
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