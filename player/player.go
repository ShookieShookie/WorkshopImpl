package player

import (
	"errors"
	"fmt"
)

type Hand interface {
	Add(int)
	Remove(int) error
	Get(int) (int, error)
	Show() []int
}
type Deck interface {
	Draw() int
	Add(int)
}

type PlayerImpl struct {
	name        string
	health      int
	manaCurrent int
	hand        Hand
	deck        Deck
}

var notEnoughMana = errors.New("not enough mana")
var noCardsInDeck = errors.New("no cards in deck")

func NewPlayer(name string, health, manaCurrent int, hand Hand, deck Deck) *PlayerImpl {
	return &PlayerImpl{
		name:        name,
		health:      health,
		manaCurrent: manaCurrent,
		hand:        hand,
		deck:        deck,
	}
}

func (p *PlayerImpl) GetHealth() int {
	return p.health
}

func (p *PlayerImpl) IsDead() bool {
	return p.health <= 0
}

func (p *PlayerImpl) ApplyDamage(damage int) {
	p.health -= damage
}
func (p *PlayerImpl) Draw() error {
	v := p.deck.Draw()
	if v == -1 {
		return noCardsInDeck
	}
	p.hand.Add(v)
	return nil
}

func (p *PlayerImpl) ID() string {
	return p.name
}

func (p *PlayerImpl) SetMana(mana int) {
	p.manaCurrent = mana
}

func (p *PlayerImpl) PrintStats() {
	fmt.Printf("Current Health %d \n", p.health)
	fmt.Printf("Current Mana %d \n", p.manaCurrent)
	fmt.Printf("Current Hand %+v \n", p.hand.Show())
}

// returns damage done
func (p *PlayerImpl) PlayCard(index int) (int, error) {

	v, err := p.hand.Get(index)
	if err != nil {
		return 0, err
	}
	if v > p.manaCurrent {
		return 0, notEnoughMana
	}
	p.manaCurrent -= v
	err = p.hand.Remove(index)
	if err != nil {
		return 0, err
	}
	return v, nil
}
