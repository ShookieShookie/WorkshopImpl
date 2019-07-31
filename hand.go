package main

import (
	"errors"
)

type HandImpl struct {
	cards []int
}

func (h *HandImpl) Add(card int) {
	if len(h.cards) < 5 {
		h.cards = append(h.cards, card)
	}
}

func (h *HandImpl) Remove(ind int) error {
	if ind >= len(h.cards) {
		return errors.New("Illegal index")
	}
	h.cards = append(h.cards[:ind], h.cards[ind+1:]...)
	return nil
}

func (h *HandImpl) Show() []int {
	return h.cards
}

func (h *HandImpl) Get(i int) (int, error) {
	if i >= len(h.cards) {
		return 0, errors.New("Illegal index")
	}
	return h.cards[i], nil
}
