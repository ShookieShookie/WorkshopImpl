package game

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockPlayer struct {
	mock.Mock
}

func (m *mockPlayer) ApplyDamage(i int) {
	m.Called(i)
}
func (m *mockPlayer) GetHealth() int {
	args := m.Called()
	return args.Get(0).(int)
}
func (m *mockPlayer) SetMana(i int) {
	m.Called(i)
}
func (m *mockPlayer) PlayCard(index int) (int, error) {
	args := m.Called(index)
	return args.Get(0).(int), args.Error(1)
}
func (m *mockPlayer) IsDead() bool {
	args := m.Called()
	return args.Bool(0)
}
func (m *mockPlayer) Draw() error {
	args := m.Called()
	return args.Error(0)

}
func (m *mockPlayer) ID() string {
	args := m.Called()
	return args.Get(0).(string)
}
func (m *mockPlayer) PrintStats() {
	m.Called()
}

type mockTurner struct {
	mock.Mock
}

func (m *mockTurner) turn(iter int, active, passive Player, getInput func() string) bool {
	args := m.Called(iter, active, passive, getInput)
	return args.Bool(0)
}

func TestGame_Start(t *testing.T) {
	type fields struct {
		userInput func() string
	}
	type turnArgs struct {
		iter     int
		active   Player
		passive  Player
		getInput func()
		over     bool
	}
	type mockDrawArgs struct {
		err error
	}
	tests := []struct {
		name       string
		fields     fields
		p1DrawArgs *mockDrawArgs
		p2DrawArgs *mockDrawArgs
		turnArgs   []turnArgs
	}{
		{
			name:       "multiple turns switch players",
			p1DrawArgs: &mockDrawArgs{err: nil},
			p2DrawArgs: &mockDrawArgs{err: nil},
			turnArgs:   []turnArgs{{iter: 1, getInput: nil}, {iter: 2, getInput: nil, over: true}},
			fields: fields{
				userInput: nil,
			},
		},
		{
			name:       "happy",
			p1DrawArgs: &mockDrawArgs{err: nil},
			p2DrawArgs: &mockDrawArgs{err: nil},
			turnArgs:   []turnArgs{{iter: 1, getInput: nil, over: true}},
			fields: fields{
				userInput: nil,
			},
		},
		{
			name:       "p1 draw fail",
			p1DrawArgs: &mockDrawArgs{err: errors.New("empty deck")},
			p2DrawArgs: nil,
			fields: fields{
				userInput: nil,
			},
		},
		{
			name:       "p2 draw fail",
			p1DrawArgs: &mockDrawArgs{err: nil},
			p2DrawArgs: &mockDrawArgs{err: errors.New("empty deck")},
			fields: fields{
				userInput: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := mockPlayer{}
			p2 := mockPlayer{}
			if tt.p1DrawArgs != nil {
				p1.On("Draw").Return(tt.p1DrawArgs.err)
			}
			if tt.p2DrawArgs != nil {
				p2.On("Draw").Return(tt.p2DrawArgs.err)
			}
			turner := &mockTurner{}
			if tt.turnArgs != nil {
				for _, a := range tt.turnArgs {
					var active *mockPlayer
					var passive *mockPlayer
					if a.iter%2 == 0 {
						active = &p2
						passive = &p1
					} else {
						active = &p1
						passive = &p2
					}
					turner.On("turn", a.iter, active, passive, mock.Anything).Return(a.over)
				}
			}
			g := &Game{
				p1:        &p1,
				p2:        &p2,
				userInput: tt.fields.userInput,
				turn:      turner.turn,
			}
			g.Start()
			p1.AssertExpectations(t)
			p2.AssertExpectations(t)
		})

	}
}

type mockUserInput struct {
	called int
	input  []string
}

func (u *mockUserInput) get() string {
	ret := u.input[u.called]
	u.called++
	return ret
}

func TestTurn(t *testing.T) {
	type ActiveSetManaArgs struct {
		in int
	}
	type ActiveDrawArgs struct {
		ret error
	}
	type ActiveApplyDamageArgs struct {
		damage int
	}
	type ActivePrintStatsArgs struct {
	}
	type ActivePlayCardArgs struct {
		index  []int
		damage []int
		err    []error
	}
	type ActiveIDArgs struct {
		ret string
	}
	type ActiveGetHealthArgs struct {
		health int
	}
	type PassiveGetHealthArgs struct {
		health int
	}
	type PassiveIDArgs struct {
		ret string
	}
	type PassiveApplyDamageArgs struct {
		damage int
	}
	type PassiveIsDeadArgs struct {
		ret bool
	}

	type args struct {
		iter      int
		userInput []string
	}
	tests := []struct {
		name                   string
		args                   args
		ActiveSetManaArgs      *ActiveSetManaArgs
		ActiveDrawArgs         *ActiveDrawArgs
		ActiveApplyDamageArgs  *ActiveApplyDamageArgs
		ActivePrintStatsArgs   *ActivePrintStatsArgs
		ActivePlayCardArgs     *ActivePlayCardArgs
		ActiveIDArgs           *ActiveIDArgs
		ActiveGetHealthArgs    *ActiveGetHealthArgs
		PassiveGetHealthArgs   *PassiveGetHealthArgs
		PassiveIDArgs          *PassiveIDArgs
		PassiveApplyDamageArgs *PassiveApplyDamageArgs
		PassiveIsDeadArgs      *PassiveIsDeadArgs
		want                   bool
	}{
		{
			name:                   "round 1 game over",
			args:                   args{iter: 1, userInput: []string{"0"}},
			ActiveSetManaArgs:      &ActiveSetManaArgs{in: 1},
			ActiveDrawArgs:         &ActiveDrawArgs{ret: nil},
			ActiveApplyDamageArgs:  nil,
			ActivePrintStatsArgs:   &ActivePrintStatsArgs{},
			ActivePlayCardArgs:     &ActivePlayCardArgs{index: []int{0}, damage: []int{5}, err: []error{nil}},
			ActiveIDArgs:           &ActiveIDArgs{"name1"},
			ActiveGetHealthArgs:    &ActiveGetHealthArgs{health: 1},
			PassiveGetHealthArgs:   &PassiveGetHealthArgs{health: -1},
			PassiveIDArgs:          &PassiveIDArgs{ret: "name2"},
			PassiveApplyDamageArgs: &PassiveApplyDamageArgs{5},
			PassiveIsDeadArgs:      &PassiveIsDeadArgs{true},
			want:                   true,
		},
		{
			name:                   "player 1 is out of deck and receives damage",
			args:                   args{iter: 1, userInput: []string{"0"}},
			ActiveSetManaArgs:      &ActiveSetManaArgs{in: 1},
			ActiveDrawArgs:         &ActiveDrawArgs{ret: errors.New("out of deck")},
			ActiveApplyDamageArgs:  &ActiveApplyDamageArgs{damage: 1},
			ActivePrintStatsArgs:   &ActivePrintStatsArgs{},
			ActivePlayCardArgs:     &ActivePlayCardArgs{index: []int{0}, damage: []int{5}, err: []error{nil}},
			ActiveIDArgs:           &ActiveIDArgs{"name1"},
			ActiveGetHealthArgs:    &ActiveGetHealthArgs{health: 1},
			PassiveGetHealthArgs:   &PassiveGetHealthArgs{health: -1},
			PassiveIDArgs:          &PassiveIDArgs{ret: "name2"},
			PassiveApplyDamageArgs: &PassiveApplyDamageArgs{5},
			PassiveIsDeadArgs:      &PassiveIsDeadArgs{true},
			want:                   true,
		},
		{
			name:                   "bad user input first try",
			args:                   args{iter: 1, userInput: []string{"asdf", "0"}},
			ActiveSetManaArgs:      &ActiveSetManaArgs{in: 1},
			ActiveDrawArgs:         &ActiveDrawArgs{ret: errors.New("out of deck")},
			ActiveApplyDamageArgs:  &ActiveApplyDamageArgs{damage: 1},
			ActivePrintStatsArgs:   &ActivePrintStatsArgs{},
			ActivePlayCardArgs:     &ActivePlayCardArgs{index: []int{0}, damage: []int{5}, err: []error{nil}},
			ActiveIDArgs:           &ActiveIDArgs{"name1"},
			ActiveGetHealthArgs:    &ActiveGetHealthArgs{health: 1},
			PassiveGetHealthArgs:   &PassiveGetHealthArgs{health: -1},
			PassiveIDArgs:          &PassiveIDArgs{ret: "name2"},
			PassiveApplyDamageArgs: &PassiveApplyDamageArgs{5},
			PassiveIsDeadArgs:      &PassiveIsDeadArgs{true},
			want:                   true,
		},
		{
			name:                   "user doesn't want to play any more cards",
			args:                   args{iter: 1, userInput: []string{"-1"}},
			ActiveSetManaArgs:      &ActiveSetManaArgs{in: 1},
			ActiveDrawArgs:         &ActiveDrawArgs{ret: nil},
			ActiveApplyDamageArgs:  nil,
			ActivePrintStatsArgs:   &ActivePrintStatsArgs{},
			ActivePlayCardArgs:     nil,
			ActiveIDArgs:           &ActiveIDArgs{"name1"},
			ActiveGetHealthArgs:    &ActiveGetHealthArgs{health: 1},
			PassiveGetHealthArgs:   &PassiveGetHealthArgs{health: -1},
			PassiveIDArgs:          &PassiveIDArgs{ret: "name2"},
			PassiveApplyDamageArgs: nil,
			PassiveIsDeadArgs:      nil,
			want:                   false,
		},
		{
			name:                   "play card illegal index causes a second turn",
			args:                   args{iter: 1, userInput: []string{"12435", "1"}},
			ActiveSetManaArgs:      &ActiveSetManaArgs{in: 1},
			ActiveDrawArgs:         &ActiveDrawArgs{ret: nil},
			ActiveApplyDamageArgs:  nil,
			ActivePrintStatsArgs:   &ActivePrintStatsArgs{},
			ActivePlayCardArgs:     &ActivePlayCardArgs{index: []int{12435, 1}, damage: []int{0, 5}, err: []error{errors.New("illegal index"), nil}},
			ActiveIDArgs:           &ActiveIDArgs{"name1"},
			ActiveGetHealthArgs:    &ActiveGetHealthArgs{health: 1},
			PassiveGetHealthArgs:   &PassiveGetHealthArgs{health: -1},
			PassiveIDArgs:          &PassiveIDArgs{ret: "name2"},
			PassiveApplyDamageArgs: &PassiveApplyDamageArgs{5},
			PassiveIsDeadArgs:      &PassiveIsDeadArgs{true},
			want:                   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			active := mockPlayer{}
			if tt.ActiveSetManaArgs != nil {
				active.On("SetMana", tt.ActiveSetManaArgs.in)
			}
			if tt.ActiveDrawArgs != nil {
				active.On("Draw").Return(tt.ActiveDrawArgs.ret)
			}
			if tt.ActiveApplyDamageArgs != nil {
				active.On("ApplyDamage", tt.ActiveApplyDamageArgs.damage)
			}
			if tt.ActivePrintStatsArgs != nil {
				active.On("PrintStats")
			}
			if tt.ActivePlayCardArgs != nil {
				for ind := range tt.ActivePlayCardArgs.index {
					active.On("PlayCard", tt.ActivePlayCardArgs.index[ind]).Return(tt.ActivePlayCardArgs.damage[ind], tt.ActivePlayCardArgs.err[ind])
				}
			}
			if tt.ActiveIDArgs != nil {
				active.On("ID").Return(tt.ActiveIDArgs.ret)
			}
			if tt.ActiveGetHealthArgs != nil {
				active.On("GetHealth").Return(tt.ActiveGetHealthArgs.health)
			}
			passive := mockPlayer{}
			if tt.PassiveGetHealthArgs != nil {
				passive.On("GetHealth").Return(tt.PassiveGetHealthArgs.health)
			}
			if tt.PassiveIDArgs != nil {
				passive.On("ID").Return(tt.PassiveIDArgs.ret)
			}
			if tt.PassiveApplyDamageArgs != nil {
				passive.On("ApplyDamage", tt.PassiveApplyDamageArgs.damage)
			}
			if tt.PassiveIsDeadArgs != nil {
				passive.On("IsDead").Return(tt.PassiveIsDeadArgs.ret)
			}
			mockUserInputInst := mockUserInput{input: tt.args.userInput}
			// g := &Game{userInput: mockUserInputInst.get}
			if got := turn(tt.args.iter, &active, &passive, mockUserInputInst.get); got != tt.want {
				t.Errorf("Turn() = %v, want %v", got, tt.want)
				assert.Equal(t, 1, 1)
			}
			active.AssertExpectations(t)
			passive.AssertExpectations(t)
		})
	}
}

func Test_min(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "first is smaller",
			args: args{
				i: 1,
				j: 2,
			},
			want: 1,
		},
		{
			name: "second is smaller",
			args: args{
				i: 2,
				j: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}
