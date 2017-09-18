package parser

import (
	"reflect"
	"testing"
)

func TestParseSample(t *testing.T) {
	type args struct {
		sampleID int
		sample   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []Token
		wantErr bool
	}{
		0: {
			"err: empty sample",
			args{0, nil},
			nil,
			true,
		},
		1: {
			"normal sample",
			args{1, []byte("play {Name} from {Artist}")},
			[]Token{
				{Val: []byte("play")},
				{Kw: true, Val: []byte("Name")},
				{Val: []byte("from")},
				{Kw: true, Val: []byte("Artist")},
			},
			false,
		},
		2: {
			"spacing inside keys",
			args{1, []byte("play { 	Name} from {	Artist		}")},
			[]Token{
				{Val: []byte("play")},
				{Kw: true, Val: []byte("Name")},
				{Val: []byte("from")},
				{Kw: true, Val: []byte("Artist")},
			},
			false,
		},
		3: {
			"multi word",
			args{1, []byte("I need {Name} since {Since}")},
			[]Token{
				{Val: []byte("I")},
				{Val: []byte("need")},
				{Kw: true, Val: []byte("Name")},
				{Val: []byte("since")},
				{Kw: true, Val: []byte("Since")},
			},
			false,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSample(tt.args.sampleID, tt.args.sample)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test#%d: ParseSample() error = %v, wantErr %v", i, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Test#%d: ParseSample() = %v, want %v", i, got, tt.want)
			}
		})
	}
}
