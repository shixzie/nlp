package parser

import (
	"reflect"
	"testing"
)

func TestParseSample(t *testing.T) {
	type args struct {
		sampleID int
		sample   string
	}
	tests := []struct {
		name    string
		args    args
		want    []Token
		wantErr bool
	}{
		{
			"err: empty sample",
			args{0, ""},
			nil,
			true,
		},
		{
			"normal sample",
			args{1, "play {Name} from {Artist}"},
			[]Token{
				Token{Val: "play"},
				Token{Kw: true, Val: "Name"},
				Token{Val: "from"},
				Token{Kw: true, Val: "Artist"},
			},
			false,
		},
		{
			"spacing inside keys",
			args{1, "play { 	Name} from {	Artist		}"},
			[]Token{
				Token{Val: "play"},
				Token{Kw: true, Val: "Name"},
				Token{Val: "from"},
				Token{Kw: true, Val: "Artist"},
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
