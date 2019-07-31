package player

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHand struct {
	mock.Mock
}

func (m *mockHand) Add(i int) {
	m.Called(i)
}
func (m *mockHand) Remove(i int) error {
	args := m.Called(i)
	return args.Error(0)
}
func (m *mockHand) Get(i int) (int, error) {
	args := m.Called(i)
	return args.Get(0).(int), args.Error(1)
}
func (m *mockHand) Show() []int {
	args := m.Called()
	return args.Get(0).([]int)
}

type mockDeck struct {
	mock.Mock
}

func (m *mockDeck) Draw() int {
	args := m.Called()
	return args.Get(0).(int)
}
func (m *mockDeck) Add(i int) {
	m.Called(i)
}

func TestPlayerImpl_GetHealth(t *testing.T) {
	type fields struct {
		name        string
		health      int
		manaCurrent int
		hand        Hand
		deck        Deck
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "get health",
			fields: fields{
				health: 5,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PlayerImpl{
				name:        tt.fields.name,
				health:      tt.fields.health,
				manaCurrent: tt.fields.manaCurrent,
				hand:        tt.fields.hand,
				deck:        tt.fields.deck,
			}
			if got := p.GetHealth(); got != tt.want {
				t.Errorf("PlayerImpl.GetHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayerImpl_IsDead(t *testing.T) {
	type fields struct {
		name        string
		health      int
		manaCurrent int
		hand        Hand
		deck        Deck
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "true",
			fields: fields{
				health: 0,
			},
			want: true,
		},
		{
			name: "false",
			fields: fields{
				health: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PlayerImpl{
				name:        tt.fields.name,
				health:      tt.fields.health,
				manaCurrent: tt.fields.manaCurrent,
				hand:        tt.fields.hand,
				deck:        tt.fields.deck,
			}
			if got := p.IsDead(); got != tt.want {
				t.Errorf("PlayerImpl.IsDead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayerImpl_ApplyDamage(t *testing.T) {
	type fields struct {
		name        string
		health      int
		manaCurrent int
		hand        Hand
		deck        Deck
	}
	type args struct {
		damage int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "apply",
			fields: fields{
				health: 1,
			},
			args: args{
				damage: 15,
			},
			want: -14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PlayerImpl{
				name:        tt.fields.name,
				health:      tt.fields.health,
				manaCurrent: tt.fields.manaCurrent,
				hand:        tt.fields.hand,
				deck:        tt.fields.deck,
			}
			p.ApplyDamage(tt.args.damage)
			assert.Equal(t, tt.want, p.health)
		})
	}
}

func TestPlayerImpl_Draw(t *testing.T) {
	type DeckArgs struct {
		ret int
	}
	type HandArgs struct {
		in int
	}
	tests := []struct {
		name     string
		DeckArgs *DeckArgs
		HandArgs *HandArgs
		wantErr  bool
	}{
		{
			name: "nonempty deck ",
			DeckArgs: &DeckArgs{
				ret: 1,
			},
			HandArgs: &HandArgs{
				in: 1,
			},
		},
		{
			name: "empty deck",
			DeckArgs: &DeckArgs{
				ret: -1,
			},
			HandArgs: nil, // inherently asserts hand is not added toname: "nonempty deck ",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := &mockDeck{}
			if tt.DeckArgs != nil {
				deck.On("Draw").Return(tt.DeckArgs.ret)
			}
			hand := &mockHand{}
			if tt.HandArgs != nil {
				hand.On("Add", tt.HandArgs.in)
			}
			p := &PlayerImpl{
				deck: deck,
				hand: hand,
			}
			if err := p.Draw(); (err != nil) != tt.wantErr {
				t.Errorf("PlayerImpl.Draw() error = %v, wantErr %v", err, tt.wantErr)
			}
			deck.AssertExpectations(t)
			hand.AssertExpectations(t)
		})
	}
}

func TestPlayerImpl_ID(t *testing.T) {
	type fields struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "show name",
			fields: fields{
				name: "helo",
			},
			want: "helo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PlayerImpl{
				name: tt.fields.name,
			}
			if got := p.ID(); got != tt.want {
				t.Errorf("PlayerImpl.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayerImpl_SetMana(t *testing.T) {
	type fields struct {
		manaCurrent int
	}
	type args struct {
		mana int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantMana int
	}{
		{
			name: "0 mana start",
			fields: fields{
				manaCurrent: 0,
			},
			args: args{
				mana: 1,
			},
			wantMana: 1,
		},
		{
			name: "overwrite",
			fields: fields{
				manaCurrent: 2,
			},
			args: args{
				mana: 1,
			},
			wantMana: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PlayerImpl{
				manaCurrent: tt.fields.manaCurrent,
			}
			p.SetMana(tt.args.mana)
			assert.Equal(t, tt.wantMana, p.manaCurrent)
		})
	}
}

func TestPlayerImpl_PrintStats(t *testing.T) {
	type mockShowArgs struct {
		ret []int
	}
	tests := []struct {
		name         string
		mockShowArgs mockShowArgs
	}{
		{
			name: "Just a util",
			mockShowArgs: mockShowArgs{
				ret: []int{1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockHand{}
			m.On("Show").Return(tt.mockShowArgs.ret)
			p := &PlayerImpl{
				hand: m,
			}
			p.PrintStats()
			m.AssertExpectations(t)
		})
	}
}

func TestPlayerImpl_PlayCard(t *testing.T) {
	type fields struct {
		manaCurrent int
	}
	type args struct {
		index int
	}
	type HandGetArgs struct {
		in  int
		out int
		err error
	}
	type HandRemoveArgs struct {
		in  int
		err error
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		HandGetArgs    HandGetArgs
		HandRemoveArgs *HandRemoveArgs
		want           int
		wantErr        bool
	}{
		{
			name: "enough mana, valid ind",
			fields: fields{
				manaCurrent: 5,
			},
			args: args{
				index: 1,
			},
			HandGetArgs: HandGetArgs{
				in:  1,
				out: 3,
				err: nil,
			},
			HandRemoveArgs: &HandRemoveArgs{
				in: 1,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "not enough mana, valid ind",
			fields: fields{
				manaCurrent: 0,
			},
			args: args{
				index: 1,
			},
			HandGetArgs: HandGetArgs{
				in:  1,
				out: 3,
				err: nil,
			},
			HandRemoveArgs: nil,
			want:           0,
			wantErr:        true,
		},
		{
			name: "invalid ind",
			fields: fields{
				manaCurrent: 0,
			},
			args: args{
				index: 1,
			},
			HandGetArgs: HandGetArgs{
				in:  1,
				out: -1,
				err: noCardsInDeck,
			},
			HandRemoveArgs: nil,
			want:           0,
			wantErr:        true,
		},
		{
			name: "enough mana, valid ind",
			fields: fields{
				manaCurrent: 5,
			},
			args: args{
				index: 1,
			},
			HandGetArgs: HandGetArgs{
				in:  1,
				out: 3,
				err: nil,
			},
			HandRemoveArgs: &HandRemoveArgs{
				in:  1,
				err: errors.New("nonsense"),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &mockHand{}
			h.On("Get", tt.HandGetArgs.in).Return(tt.HandGetArgs.out, tt.HandGetArgs.err)
			if tt.HandRemoveArgs != nil {
				h.On("Remove", tt.HandRemoveArgs.in).Return(tt.HandRemoveArgs.err)
			}
			p := &PlayerImpl{
				manaCurrent: tt.fields.manaCurrent,
				hand:        h,
			}
			got, err := p.PlayCard(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlayerImpl.PlayCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PlayerImpl.PlayCard() = %v, want %v", got, tt.want)
			}
			h.AssertExpectations(t)
		})
	}
}

func TestNewPlayer(t *testing.T) {
	type args struct {
		name        string
		health      int
		manaCurrent int
		hand        Hand
		deck        Deck
	}
	tests := []struct {
		name string
		args args
		want *PlayerImpl
	}{
		{
			name: "happy",
			args: args{
				name:        "name",
				health:      1,
				manaCurrent: 0,
				hand:        &mockHand{},
				deck:        &mockDeck{},
			},
			want: &PlayerImpl{
				name:        "name",
				health:      1,
				manaCurrent: 0,
				hand:        &mockHand{},
				deck:        &mockDeck{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlayer(tt.args.name, tt.args.health, tt.args.manaCurrent, tt.args.hand, tt.args.deck); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}
