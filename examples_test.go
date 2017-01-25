package nlp_test

import (
	"fmt"

	"github.com/Shixzie/nlp"
)

func ExampleNL_RegisterModel() {
	type song struct {
		Name   string
		Artist string
	}
	samples := []string{
		"play {Name} from {Artist}",
		"play {Name} by {Artist}",
		"from {Artist} play {Name}",
	}
	nl := nlp.New()
	err := nl.RegisterModel(song{}, samples)
	if err != nil {
		panic(err)
	}
}

func ExampleClassifier_Classify() {
	positiveSamples := []string{
		"i love you",
		"i like chocolate",
		"you're pretty",
		"i like you",
		"good",
	}
	negativeSamples := []string{
		"you're ugly",
		"i don't like you",
		"he doesn't want you near",
		"i hate it",
		"i dislike you",
		"bad",
	}
	cls := nlp.NewClassifier()
	err := cls.NewClass("positive", positiveSamples)
	if err != nil {
		panic(err)
	}
	err = cls.NewClass("negative", negativeSamples)
	if err != nil {
		panic(err)
	}
	err = cls.Learn()
	if err != nil {
		panic(err)
	}
	fmt.Printf("got %s | want %s\n", cls.Classify("water is good"), "positive")
	fmt.Printf("got %s | want %s\n", cls.Classify("she says chocolate is bad"), "negative")
}
