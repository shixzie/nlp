package nlp_test

import (
	"fmt"
	"testing"

	"github.com/Shixzie/nlp"
)

func TestP(t *testing.T) {
	type Order struct {
		Product  string
		Quantity int
	}

	orderSamples := []string{
		"dame {Quantity} , {Product}",
		"ordena {Quantity} , {Product}",
		"compra un {Product}",
		"compra {Quantity} , {Product}",
	}

	nl := nlp.New()

	err := nl.RegisterModel(Order{}, orderSamples)
	if err != nil {
		panic(err)
	}

	err = nl.Learn() // you must call Learn after all models are registered and before calling P
	if err != nil {
		panic(err)
	}
	o := nl.P("compra 250 , cajas vacías") // after learning you can call P the times you want
	if order, ok := o.(Order); ok {
		fmt.Println("Success")
		fmt.Printf("%#v\n", order)
	} else {
		fmt.Println("Failed")
	}
	// Prints
	//
	// Success
	// nlp_test.Order{Product: "King", Quantity: 25}
}
