package nlp

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/cdipaolo/goml/text"
)

func TestNL_P(t *testing.T) {
	type Order struct {
		Product  string
		Quantity int
	}

	orderSamples := []string{
		"dame {Quantity}, {Product}",
		"ordena {Quantity}, {Product}",
		"compra un {Product}",
		"compra {Quantity}, {Product}",
	}

	nl := New()

	err := nl.RegisterModel(Order{}, orderSamples)
	if err != nil {
		panic(err)
	}

	err = nl.Learn() // you must call Learn after all models are registered and before calling P
	if err != nil {
		panic(err)
	}
	o := nl.P("compra 250 , cajas vacías") // after learning you can call P the times you want
	if order, ok := o.(*Order); ok {
		if order.Product != "cajas vacías" || order.Quantity != 250 {
			t.Error("wrong values")
		}
	} else {
		t.Error("not an order")
	}
	// Prints
	//
	// Success
	// &nlp_test.Order{Product: "cajas vacías", Quantity: 250}
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
			args{nil, nil},
			true,
		},
		{
			"nil samples",
			fields{},
			args{args{}, nil},
			true,
		},
		{
			"non-struct",
			fields{},
			args{[]int{}, []string{""}},
			true,
		},
		{
			"unexported & time.Time",
			fields{},
			args{T{}, []string{""}},
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
			if err := nl.RegisterModel(tt.args.i, tt.args.samples); (err != nil) != tt.wantErr {
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
