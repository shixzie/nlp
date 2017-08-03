package nlp

import (
	"testing"
	"time"
)

func BenchmarkNL_P(b *testing.B) {
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

	nl.RegisterModel(T{}, tSamples)

	err := nl.RegisterModel(T{}, tSamples)
	if err != nil {
		b.Error(err)
	}

	err = nl.Learn()
	if err != nil {
		b.Error(err)
	}

	tim, err := time.ParseInLocation("01-02-2006_3:04pm", "05-18-1999_6:42pm", time.Local)
	if err != nil {
		b.Error(err)
	}

	dur, err := time.ParseDuration("4h2m")
	if err != nil {
		b.Error(err)
	}

	cases := []struct {
		name       string
		expression string
		want       interface{}
	}{
		{
			"string",
			"string Hello World",
			"Hello World",
		},
		{
			"int",
			"int 42",
			int(42),
		},
		{
			"uint",
			"uint 43",
			uint(43),
		},
		{
			"float",
			"float 44",
			float32(44),
		},
		{
			"time",
			"time 05-18-1999_6:42pm",
			tim,
		},
		{
			"duration",
			"dur 4h2m",
			dur,
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

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			nl.P(c.expression)
		})
	}
}
