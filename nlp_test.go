package nlp

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/cdipaolo/goml/text"
)

func TestNL_P(t *testing.T) {
	type T struct {
		String string
		Int    int
		Uint   uint
		Float  float32
		Time   time.Time
		Dur    time.Duration
	}

	tSamples := []string{
		"string {String}",
		"int {Int}",
		"uint {Uint}",
		"float {Float}",
		"time {Time}",
		"dur {Dur}",
		"string {String} int {Int}",
		"string {String} time {Time}",
	}

	nl := New()

	err := nl.RegisterModel(T{}, tSamples)
	if err != nil {
		t.Error(err)
	}

	err = nl.Learn()
	if err != nil {
		t.Error(err)
	}

	tim, err := time.ParseInLocation("01-02-2006_3:04pm", "05-18-1999_6:42pm", time.Local)
	if err != nil {
		t.Error(err)
	}

	dur, err := time.ParseDuration("4h2m")
	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		expression string
		want       interface{}
	}{
		{
			"string Hello World",
			"Hello World",
		},
		{
			"int 42",
			int(42),
		},
		{
			"uint 43",
			uint(43),
		},
		{
			"float 44",
			float32(44),
		},
		{
			"time 05-18-1999_6:42pm",
			tim,
		},
		{
			"dur 4h2m",
			dur,
		},
		{
			"string Lmao int 42",
			&T{
				String: "Lmao",
				Int:    42,
			},
		},
		{
			"string What's Up Boy time 05-18-1999_6:42pm",
			&T{
				String: "What's Up Boy",
				Time:   tim,
			},
		},
	}
	for i, test := range cases {
		res := nl.P(test.expression)
		if c, ok := res.(*T); ok {
			switch v := test.want.(type) {
			case string:
				if c.String != v {
					t.Errorf("test#%d: got %v want %v", i, c.String, v)
				}
			case int:
				if c.Int != v {
					t.Errorf("test#%d: got %v want %v", i, c.Int, v)
				}
			case uint:
				if c.Uint != v {
					t.Errorf("test#%d: got %v want %v", i, c.Uint, v)
				}
			case float32:
				if c.Float != v {
					t.Errorf("test#%d: got %v want %v", i, c.Float, v)
				}
			case time.Time:
				if !c.Time.Equal(v) {
					t.Errorf("test#%d: got %v want %v", i, c.Time, v)
				}
			case time.Duration:
				if c.Dur != v {
					t.Errorf("test#%d: got %v want %v", i, c.Dur, v)
				}
			case *T:
				if !reflect.DeepEqual(c, v) {
					t.Errorf("test#%d: got %v want %v", i, c, v)
				}
			}
		} else {
			t.Errorf("test#%d: got %T want %T", i, res, &T{})
		}
	}
}

func TestNL_RegisterModel(t *testing.T) {
	type fields struct {
		models []*model
		naive  *text.NaiveBayes
		Output *bytes.Buffer
	}
	type args struct {
		i       interface{}
		samples []string
		ops     []ModelOption
	}
	type T struct {
		unexported int
		Time       time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"nil struct",
			fields{},
			args{nil, nil, nil},
			true,
		},
		{
			"nil samples",
			fields{},
			args{args{}, nil, nil},
			true,
		},
		{
			"non-struct",
			fields{},
			args{[]int{}, []string{""}, nil},
			true,
		},
		{
			"unexported & time.Time",
			fields{},
			args{T{}, []string{""}, nil},
			false,
		},
		{
			"options",
			fields{},
			args{T{}, []string{""}, []ModelOption{
				WithTimeFormat("02-01-2006"),
				WithTimeLocation(time.Local),
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := &NL{
				models: tt.fields.models,
				naive:  tt.fields.naive,
				Output: tt.fields.Output,
			}
			if err := nl.RegisterModel(tt.args.i, tt.args.samples, tt.args.ops...); (err != nil) != tt.wantErr {
				t.Errorf("NL.RegisterModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNL_Learn(t *testing.T) {
	type fields struct {
		models []*model
		naive  *text.NaiveBayes
		Output *bytes.Buffer
	}
	type T struct {
		Name string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"no models",
			fields{},
			true,
		},
		{
			"empty model sample",
			fields{
				models: []*model{
					&model{
						samples: []string{""},
					},
				},
				Output: bytes.NewBufferString(""),
			},
			true,
		},
		{
			"mistyped field",
			fields{
				models: []*model{
					&model{
						samples: []string{"Hello {Namee}"},
					},
				},
				Output: bytes.NewBufferString(""),
			},
			true,
		},
		{
			"sample with no keys",
			fields{
				models: []*model{
					&model{
						samples: []string{"Hello"},
					},
				},
				Output: bytes.NewBufferString(""),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := &NL{
				models: tt.fields.models,
				naive:  tt.fields.naive,
				Output: tt.fields.Output,
			}
			if err := nl.Learn(); (err != nil) != tt.wantErr {
				t.Errorf("NL.Learn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_model_learn(t *testing.T) {
	type fields struct {
		tpy      reflect.Type
		fields   []field
		expected [][]expected
		samples  []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &model{
				tpy:      tt.fields.tpy,
				fields:   tt.fields.fields,
				expected: tt.fields.expected,
				samples:  tt.fields.samples,
			}
			if err := m.learn(); (err != nil) != tt.wantErr {
				t.Errorf("model.learn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
