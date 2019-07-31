package hand

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestHandImpl_Add(t *testing.T) {
	type fields struct {
		cards []int
	}
	type args struct {
		card int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			name: "empty",
			fields: fields{
				cards: []int{},
			},
			args: args{
				card: 1,
			},
			want: []int{1},
		},
		{
			name: "full discards",
			fields: fields{
				cards: []int{1, 2, 3, 4, 5},
			},
			args: args{
				card: 66,
			},
			want: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandImpl{
				cards: tt.fields.cards,
			}
			h.Add(tt.args.card)
			assert.Equal(t, tt.want, h.cards)
		})
	}
}

func TestHandImpl_Remove(t *testing.T) {
	type fields struct {
		cards []int
	}
	type args struct {
		ind int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    []int
	}{
		{
			name: "remove legal",
			fields: fields{
				cards: []int{1, 2, 3, 4, 5},
			},
			args: args{
				ind: 1,
			},
			wantErr: false,
			want:    []int{1, 3, 4, 5},
		},
		{
			name: "remove illegal",
			fields: fields{
				cards: []int{1, 2, 3, 4, 5},
			},
			args: args{
				ind: 6,
			},
			wantErr: true,
			want:    []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandImpl{
				cards: tt.fields.cards,
			}
			if err := h.Remove(tt.args.ind); (err != nil) != tt.wantErr {
				t.Errorf("HandImpl.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, h.cards)
		})
	}
}

func TestHandImpl_Show(t *testing.T) {
	type fields struct {
		cards []int
	}
	tests := []struct {
		name   string
		fields fields
		want   []int
	}{
		{
			name: "show",
			fields: fields{
				cards: []int{1, 2, 3, 4},
			},
			want: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandImpl{
				cards: tt.fields.cards,
			}
			if got := h.Show(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandImpl.Show() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandImpl_Get(t *testing.T) {
	type fields struct {
		cards []int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "get legal",
			fields: fields{
				cards: []int{1, 2, 3, 4, 5},
			},
			args: args{
				i: 1,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "get illegal",
			fields: fields{
				cards: []int{1, 2, 3, 4, 5},
			},
			args: args{
				i: 10,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandImpl{
				cards: tt.fields.cards,
			}
			got, err := h.Get(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandImpl.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HandImpl.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
