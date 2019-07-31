package main

import (
	"fmt"
	"strconv"
)

type Game struct {
	p1        Player
	p2        Player
	userInput func() string
	turn      func(iter int, active, passive Player, getInput func() string) bool
}

func (g *Game) Start() {
	fmt.Println("GAME START")
	for i := 0; i < 3; i++ {
		if err := g.p1.Draw(); err != nil {
			fmt.Println("Cannot start game with inadequate sized deck")
			return
		}

		if err := g.p2.Draw(); err != nil {
			fmt.Println("Cannot start game with inadequate sized deck")
			return
		}
	}
	count := 0
	active := g.p1
	passive := g.p2
	for {
		count++ // does this increase once both players have gone?
		over := g.turn(count, active, passive, g.userInput)
		if over {
			break
		}
		t := active
		active = passive
		passive = t
	}
	fmt.Println("GAME OVER")
}

func turn(iter int, active, passive Player, getInput func() string) bool {
	fmt.Printf("%s's turn!\n", active.ID())
	active.SetMana(min(iter, 10))
	err := active.Draw()
	if err != nil {
		active.ApplyDamage(1) // no deck
	}
	for {
		fmt.Println(active.ID(), "health:", active.GetHealth(), passive.ID(), "health:", passive.GetHealth())
		active.PrintStats()
		s := getInput()
		i, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Invalid integer")
			continue
		}
		if i == -1 {
			break
		}
		damage, err := active.PlayCard(i)
		if err != nil {
			fmt.Println(err)
			continue
		}
		passive.ApplyDamage(damage)
		if passive.IsDead() {
			fmt.Println(passive.ID(), "Is Dead!")
			fmt.Println(active.ID(), "WINS!")
			return true
		}
	}
	return false
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
