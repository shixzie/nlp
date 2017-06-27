// Package nlp provides general purpose Natural Language Processing.
package nlp

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Shixzie/nlp/parser"
	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
)

// NL is a Natural Language Processor
type NL struct {
	models []*model
	naive  *text.NaiveBayes
	// Output contains the training output for the
	// NaiveBayes algorithm
	Output *bytes.Buffer
}

// New returns a *NL
func New() *NL { return &NL{Output: bytes.NewBufferString("")} }

// P proccesses the expr and returns one of
// the types passed as the i parameter to the RegistryModel
// func filled with the data inside expr
func (nl *NL) P(expr string) interface{} { return nl.models[nl.naive.Predict(expr)].fit(expr) }

// Learn maps the models samples to the models themselves and
// returns an error if something occurred while learning
func (nl *NL) Learn() error {
	if len(nl.models) > 0 {
		stream := make(chan base.TextDatapoint)
		errors := make(chan error)
		nl.naive = text.NewNaiveBayes(stream, uint8(len(nl.models)), base.OnlyWordsAndNumbers)
		nl.naive.Output = nl.Output
		go nl.naive.OnlineLearn(errors)
		for i := range nl.models {
			err := nl.models[i].learn()
			if err != nil {
				return err
			}
			for _, s := range nl.models[i].samples {
				stream <- base.TextDatapoint{
					X: s,
					Y: uint8(i),
				}
			}
		}
		close(stream)
		for {
			err := <-errors
			if err != nil {
				return fmt.Errorf("error occurred while learning: %s", err)
			}
			// training is done!
			break
		}
		return nil
	}
	return fmt.Errorf("register at least one model before learning")
}

type model struct {
	tpy          reflect.Type
	fields       []field
	expected     [][]expected
	samples      []string
	timeFormat   string
	timeLocation *time.Location
}

type expected struct {
	limit bool
	value string
	field field
}

type field struct {
	index int
	name  string
	kind  interface{}
}

// ModelOption is an option for a specific model
type ModelOption func(*model)

// WithTimeFormat sets the format used in time.Parse(format, val),
// note that format can't contain any spaces, the default is 01-02-2006_3:04pm
func WithTimeFormat(format string) ModelOption {
	return func(m *model) {
		m.timeFormat = strings.Replace(format, " ", "", -1)
	}
}

// WithTimeLocation sets the location used in time.ParseInLocation(format, value, loc),
// the default is time.Local
func WithTimeLocation(loc *time.Location) ModelOption {
	return func(m *model) {
		if loc != nil {
			m.timeLocation = loc
		}
	}
}

// RegisterModel registers a model i and creates possible patterns
// from samples, the default layout when parsing time is 01-02-2006_3:04pm
// and the default location is time.Local.
// Samples must have special formatting:
//
//	"play {Name} by {Artist}"
func (nl *NL) RegisterModel(i interface{}, samples []string, ops ...ModelOption) error {
	if i == nil {
		return fmt.Errorf("can't create model from nil value")
	}
	if len(samples) == 0 {
		return fmt.Errorf("samples can't be nil or empty")
	}
	tpy, val := reflect.TypeOf(i), reflect.ValueOf(i)
	if tpy.Kind() == reflect.Struct {
		mod := &model{
			tpy:          tpy,
			samples:      samples,
			expected:     make([][]expected, len(samples)),
			timeFormat:   "01-02-2006_3:04pm",
			timeLocation: time.Local,
		}
		for _, op := range ops {
			op(mod)
		}
	NextField:
		for i := 0; i < tpy.NumField(); i++ {
			if tpy.Field(i).Anonymous || tpy.Field(i).PkgPath != "" {
				continue NextField
			}
			if v, ok := val.Field(i).Interface().(time.Time); ok {
				mod.fields = append(mod.fields, field{i, tpy.Field(i).Name, v})
				continue NextField
			}
			switch val.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.String:
				mod.fields = append(mod.fields, field{i, tpy.Field(i).Name, val.Field(i).Kind()})
			}
		}
		nl.models = append(nl.models, mod)
		return nil
	}
	return fmt.Errorf("can't create model from non-struct type")
}

func (m *model) learn() error {
	for sid, s := range m.samples {
		tokens, err := parser.ParseSample(sid, s)
		if err != nil {
			return err
		}
		var exps []expected
		var hasAtLeastOneKey bool
		for _, tk := range tokens {
			if tk.Kw {
				hasAtLeastOneKey = true
				mistypedField := true
				for _, f := range m.fields {
					if tk.Val == f.name {
						mistypedField = false
						exps = append(exps, expected{field: f, value: tk.Val})
					}
				}
				if mistypedField {
					return fmt.Errorf("sample#%d: mistyped field %q", sid, tk.Val)
				}
			} else {
				exps = append(exps, expected{limit: true, value: tk.Val})
			}
		}
		if !hasAtLeastOneKey {
			return fmt.Errorf("sample#%d: need at least one keyword", sid)
		}
		m.expected[sid] = exps
	}
	return nil
}

func (m *model) selectBestSample(expr string) []expected {
	// slice [sample_id]score
	scores := make([]int, len(m.samples))

	tokens, err := parser.ParseSample(0, expr)
	if err != nil {
		return nil
	}

	// fmt.Printf("tokens: %v\n", tokens)

	mapping := make([][]expected, len(m.samples))
	limitsOrder := make([][]string, len(m.samples)+1)

	for sid, exps := range m.expected {
		var currentVal []string
		var reading bool
		var lastToken int
		isLimit := func(s string) bool {
			for _, e := range exps {
				if e.value == s {
					return true
				}
			}
			return false
		}
	expecteds:
		for _, e := range exps {
			// fmt.Printf("expecting: %v - limit: %v\n", e.value, e.limit)
			if e.limit {
				reading = false
				limitsOrder[sid+1] = append(limitsOrder[sid+1], e.value)
			} else {
				reading = true
			}
			// fmt.Printf("reading: %v\n", reading)
			for i := lastToken; i < len(tokens); i++ {
				t := tokens[i]
				// fmt.Printf("token: %v - isLimit: %v\n", t.Val, isLimit(t.Val))
				if isLimit(t.Val) {
					if sid == 0 {
						limitsOrder[0] = append(limitsOrder[0], t.Val)
					}
					scores[sid] = scores[sid] + 1
					if len(currentVal) > 0 {
						// fmt.Printf("appending: %v {%v}\n", strings.Join(currentVal, " "), e.field.n)
						mapping[sid] = append(mapping[sid], expected{field: e.field, value: strings.Join(currentVal, " ")})
						currentVal = currentVal[:0]
						lastToken = i
						continue expecteds
					}
					lastToken = i + 1
					continue expecteds
				} else {
					if reading {
						// fmt.Printf("adding: %v\n", t.Val)
						currentVal = append(currentVal, t.Val)
					}
				}
			}
			if len(currentVal) > 0 {
				// fmt.Printf("appending: %v {%v}\n", strings.Join(currentVal, " "), e.field.n)
				mapping[sid] = append(mapping[sid], expected{field: e.field, value: strings.Join(currentVal, " ")})
			}
		}
		// fmt.Printf("\n\n")
	}
order:
	for i := 1; i < len(limitsOrder); i++ {
		if len(limitsOrder[0]) < len(limitsOrder[i]) {
			continue order
		}
		for j := range limitsOrder[i] {
			if limitsOrder[i][j] != limitsOrder[0][j] {
				continue order
			}
		}
		scores[i-1] = scores[i-1] + 1
	}

	// fmt.Printf("orders: %v\n\n", limitsOrder)
	// fmt.Printf("scores: %v\n", scores)

	bestScore, bestMapping := -1, -1
	for id, score := range scores {
		if score > bestScore {
			bestScore = score
			bestMapping = id
		}
	}
	if bestScore == -1 {
		return nil
	}
	return mapping[bestMapping]
}

func (m *model) fit(expr string) interface{} {
	val := reflect.New(m.tpy)
	exps := m.selectBestSample(expr)
	if len(exps) > 0 {
		for _, e := range exps {
			switch t := e.field.kind.(type) {
			case reflect.Kind:
				switch t {
				case reflect.String:
					val.Elem().Field(e.field.index).SetString(e.value)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					v, _ := strconv.ParseUint(e.value, 10, 0)
					val.Elem().Field(e.field.index).SetUint(v)
				case reflect.Float32, reflect.Float64:
					v, _ := strconv.ParseFloat(e.value, 64)
					val.Elem().Field(e.field.index).SetFloat(v)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v, _ := strconv.ParseInt(e.value, 10, 0)
					val.Elem().Field(e.field.index).SetInt(v)
				}
			case time.Time:
				v, _ := time.ParseInLocation(m.timeFormat, e.value, m.timeLocation)
				val.Elem().Field(e.field.index).Set(reflect.ValueOf(v))
			}
		}
	}
	return val.Interface()
}
