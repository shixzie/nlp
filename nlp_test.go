package nlp_test

import (
	"fmt"
	"testing"

	"github.com/Shixzie/nlp"
)

func TestP(t *testing.T) {
	type Song struct {
		Name   string
		Artist string
	}

	songSamples := []string{
		"play {Name} by {Artist}",
		"play {Name} from {Artist}",
		"play {Name}",
		"from {Artist} play {Name}",
	}

	nl := nlp.New()

	err := nl.RegisterModel(Song{}, songSamples)
	if err != nil {
		panic(err)
	}

	err = nl.Learn() // you must call Learn after all models are registered and before calling P
	if err != nil {
		panic(err)
	}
	s := nl.P("hello sir can you pleeeeeease play King by Lauren Aquilina") // after learning you can call P the times you want
	if song, ok := s.(Song); ok {
		fmt.Println("Success")
		fmt.Printf("%#v\n", song)
	} else {
		fmt.Println("Failed")
	}
	// Prints
	//
	// Success
	// nlp_test.Song{Name: "King", Artist: "Lauren Aquilina"}
}
