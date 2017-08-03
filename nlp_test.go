package nlp

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/cdipaolo/goml/text"
)

func failTest(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

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
	failTest(t, err)

	err = nl.Learn()
	failTest(t, err)

	tim, err := time.ParseInLocation("01-02-2006_3:04pm", "05-18-1999_6:42pm", time.Local)
	failTest(t, err)

	dur, err := time.ParseDuration("4h2m")
	failTest(t, err)

	cases := []struct {
		name       string
		expression string
		want       *T
	}{
		{
			"string",
			"string Hello World",
			&T{String: "Hello World"},
		},
		{
			"int",
			"int 42",
			&T{Int: 42},
		},
		{
			"uint",
			"uint 43",
			&T{Uint: 43},
		},
		{
			"float",
			"float 44",
			&T{Float: 44},
		},
		{
			"time",
			"time 05-18-1999_6:42pm",
			&T{Time: tim},
		},
		{
			"duration",
			"dur 4h2m",
			&T{Dur: dur},
		},
		{
			"string int",
			"string Lmao int 42",
			&T{
				String: "Lmao",
				Int:    42,
			},
		},
		{
			"string time",
			"string What's Up Boy time 05-18-1999_6:42pm",
			&T{
				String: "What's Up Boy",
				Time:   tim,
			},
		},
	}
	for i, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if res := nl.P(tt.expression); !reflect.DeepEqual(res, tt.want) {
				t.Errorf("test#%d: got %v want %v", i, res, tt.want)
			}
		})
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
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := &NL{
				models: tt.fields.models,
				naive:  tt.fields.naive,
				Output: tt.fields.Output,
			}
			if err := nl.RegisterModel(tt.args.i, tt.args.samples, tt.args.ops...); (err != nil) != tt.wantErr {
				t.Errorf("[%d] NL.RegisterModel() error = %v, wantErr %v", i, err, tt.wantErr)
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
					{
						samples: [][]byte{{}},
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
					{
						samples: [][]byte{[]byte("Hello {Namee}")},
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
					{
						samples: [][]byte{[]byte("Hello")},
					},
				},
				Output: bytes.NewBufferString(""),
			},
			true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := &NL{
				models: tt.fields.models,
				naive:  tt.fields.naive,
				Output: tt.fields.Output,
			}
			if err := nl.Learn(); (err != nil) != tt.wantErr {
				t.Errorf("[%d] NL.Learn() error = %v, wantErr %v", i, err, tt.wantErr)
			}
		})
	}
}

func TestWithTimeFormat(t *testing.T) {
	type args struct {
		format string
		m      *model
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid format",
			args{"2006 01 02", &model{}},
			true,
		},
		{
			"valid format",
			args{"2006", &model{}},
			false,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := WithTimeFormat(tt.args.format)
			if err := op(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("[%d] WithTimeFormat() error = %v, wantErr %v", i, err, tt.wantErr)
			}
		})
	}
}

func TestWithTimeLocation(t *testing.T) {
	type args struct {
		loc *time.Location
		m   *model
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid location",
			args{nil, &model{}},
			true,
		},
		{
			"valid format",
			args{time.Local, &model{}},
			false,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := WithTimeLocation(tt.args.loc)
			if err := op(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("[%d] WithTimeFormat() error = %v, wantErr %v", i, err, tt.wantErr)
			}
		})
	}
}
