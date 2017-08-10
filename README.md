[![GoDoc](https://godoc.org/github.com/shixzie/nlp?status.svg)](https://godoc.org/github.com/shixzie/nlp) 
[![Go Report Card](https://goreportcard.com/badge/github.com/shixzie/nlp)](https://goreportcard.com/report/github.com/shixzie/nlp)
[![Build Status](https://travis-ci.org/shixzie/nlp.svg?branch=master)](https://travis-ci.org/shixzie/nlp)
[![codecov](https://codecov.io/gh/shixzie/nlp/branch/master/graph/badge.svg)](https://codecov.io/gh/shixzie/nlp)


# nlp

> `nlp` is a general purpose any-lang Natural Language Processor that parses the data inside a text and returns a filled model

## Supported types
```go
int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64
float32 float64
string
time.Time
time.Duration
```

## Installation
```
go get github.com/Shixzie/nlp
```

**Go 1.7+ is required**


**Feel free to create PR's and open Issues :)**

## How it works

You will always begin by creating a NL type calling nlp.New(), the NL type is a 
Natural Language Processor that owns 3 funcs, RegisterModel(), Learn() and P().

### RegisterModel(i interface{}, samples []string, ops ...ModelOption) error

RegisterModel takes 3 parameters, an empty struct, a set of samples and some options for the model.

The empty struct lets nlp know all possible values inside the text, for example:
```go
type Song struct {
	Name        string // fields must be exported
	Artist      string
	ReleasedAt  time.Time
}
err := nl.RegisterModel(Song{}, someSamples, nlp.WithTimeFormat("2006"))
if err != nil {
	panic(err)
}
// ...
```

tells nlp that inside the text may be a Song.Name, a Song.Artist and a Song.ReleasedAt.

The samples are the key part about nlp, not just because they set the *limits*
between *keywords* but also because they will be used to choose which model 
use to handle an expression.

Samples must have a special syntax to set those *limits* and *keywords*.
```go
songSamples := []string{
	"play {Name} by {Artist}",
	"play {Name} from {Artist}",
	"play {Name}",
	"from {Artist} play {Name}",
	"play something from {ReleasedAt}",
}
```

In the example below, you can see we're reffering to the Name and Artist fields
of the `Song` type declared above, both `{Name}` and `{Artist}` are our *keywords* 
and yes! you guessed it! Everything between `play` and `by` will be treated as a
`{Name}`, and everything that's after `by` will be treated as an `{Artist}` meaning 
that `play` and `by` are our *limits*.
```
     limits
 ┌─────┴─────┐
┌┴─┐        ┌┴┐
play {Name} by  {Artist}
     └─┬──┘     └───┬──┘
       └──────┬─────┘
           keywords
```

Any character can be a *limit*, a `,` for example can be used as a limit.

*keywords* as well as *limits* are `CaseSensitive` so be sure to type them right.

**Note that putting 2 *keywords* together will cause that only 1 or none of them will be detected**

> *limits are important* - Me :3


### Learn() error

Learn maps all models samples to their respective models using the NaiveBayes 
algorithm based on those samples. `Learn()` also trains all registered models
so they're able to fit expressions in the future.

```go
// must call after all models are registrated and before calling nl.P()
err := nl.Learn() 
if err != nil {
	panic(err)
}
// ...
```

Once the algorithm has finished learning, we're now ready to start Processing 
those texts.

**Note that you must call NL.Learn() after all models are registrated and before calling NL.P()**

### P(expr string) interface{}

P first asks the trained algorithm which model should be used, once we get
the right *and already trained* model, we just make it fit the expression.

**Note that everything in the expression must be separated by a _space_ or _tab_**

When processing an expression, nlp searches for the *limits* inside that 
expression and evaluates which sample fits better the expression, it doesn't
matter if the text has `trash`. In this example:
```
     limits
 ┌─────┴─────┐
┌┴─┐        ┌┴┐
play {Name} by  {Artist}
     └─┬──┘     └───┬──┘
       └──────┬─────┘
           keywords
```

we have 2 *limits*, `play` and `by`, it doesn't matter if we had an expression 
*hello sir can you pleeeeeease play King by Lauren Aquilina*, since:
```
                                  limits
            trash              ┌────┴────┐
┌─────────────┴─────────────┐ ┌┴─┐      ┌┴┐
hello sir can you pleeeeeease play King by  Lauren Aquilina
                                   └┬─┘     └─────┬───────┘
                                 {Name}       {Artist}
                                 └─┬──┘       └───┬──┘
                                   └──────┬───────┘
                                       keywords
```

`{Name}` would be replaced with `King`, 
`{Artist}` would be replaced with `Lauren Aquilina`, 
`trash` would be ignored as well as the *limits* `play` and `by`, 
and then **a pointer to a filled struct with the type used to register the model** (`Song`) 
( `Song.Name` being `{Name}` and `Song.Artist` beign `{Artist}` ) 
**will be returned**.

## Usage

```go
type Song struct {
	Name       string
	Artist     string
	ReleasedAt time.Time
}

songSamples := []string{
	"play {Name} by {Artist}",
	"play {Name} from {Artist}",
	"play {Name}",
	"from {Artist} play {Name}",
	"play something from {ReleasedAt}",
}

nl := nlp.New()
err := nl.RegisterModel(Song{}, songSamples, nlp.WithTimeFormat("2006"))
if err != nil {
	panic(err)
}

err = nl.Learn() // you must call Learn after all models are registered and before calling P
if err != nil {
	panic(err)
}

// after learning you can call P the times you want
s := nl.P("hello sir can you pleeeeeease play King by Lauren Aquilina") 
if song, ok := s.(*Song); ok {
	fmt.Println("Success")
	fmt.Printf("%#v\n", song)
} else {
	fmt.Println("Failed")
}

// Prints
//
// Success
// &main.Song{Name: "King", Artist: "Lauren Aquilina"}
```
