package deck

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeckImpl_Draw(t *testing.T) {
	type fields struct {
		getIndex func(int) int
		cards    []int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "call to nonempty deck",
			fields: fields{
				cards:    []int{1, 2, 3, 4, 5},
				getIndex: func(i int) int { return 2 },
			},
			want: 3,
		},
		{
			name: "call to empty deck",
			fields: fields{
				cards:    []int{},
				getIndex: func(i int) int { return 2 },
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DeckImpl{
				getIndex: tt.fields.getIndex,
				cards:    tt.fields.cards,
			}
			if got := d.Draw(); got != tt.want {
				t.Errorf("DeckImpl.Draw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeckImpl_Add(t *testing.T) {
	type fields struct {
		getIndex func(i int) int
		cards    []int
	}
	type args struct {
		c int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			name: "add card",
			fields: fields{
				getIndex: nil,
				cards:    []int{},
			},
			args: args{
				c: 1,
			},
			want: []int{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DeckImpl{
				getIndex: tt.fields.getIndex,
				cards:    tt.fields.cards,
			}
			d.Add(tt.args.c)
			assert.Equal(t, tt.want, d.cards)
		})
	}
}

func TestNewDeck(t *testing.T) {
	type args struct {
		getIndex func(int) int
	}
	tests := []struct {
		name string
		args args
		want *DeckImpl
	}{
		{
			name: "happy",
			args: args{
				getIndex: nil,
			},
			want: &DeckImpl{
				cards:    []int{},
				getIndex: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeck(tt.args.getIndex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeck() = %v, want %v", got, tt.want)
			}
		})
	}
}
