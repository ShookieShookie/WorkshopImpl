package deck

type DeckImpl struct {
	getIndex func(int) int
	cards    []int
}

func (d *DeckImpl) Draw() int {
	if len(d.cards) == 0 {
		return -1
	}
	ind := d.getIndex(len(d.cards))
	val := d.cards[ind]
	d.cards = append(d.cards[:ind], d.cards[ind+1:]...)
	return val
}

func (d *DeckImpl) Add(c int) {
	d.cards = append(d.cards, c)
}
