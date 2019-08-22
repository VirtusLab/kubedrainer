package stringer_test

import (
	"fmt"

	"github.com/VirtusLab/kubedrainer/internal/stringer"
)

type tester struct {
	one    string
	two    *string
	nested *tester
}

func ExampleStringify() {
	var value = []string{
		"one", "two", "three",
	}
	valuePtr := &value
	//lint:ignore S1021 needed for the test
	var valueInter interface{}
	valueInter = value
	valueInterPtr := &valueInter
	fmt.Println(stringer.Stringify(nil))
	fmt.Println(stringer.Stringify(""))
	fmt.Println(stringer.Stringify(value))
	fmt.Println(stringer.Stringify(valuePtr))
	fmt.Println(stringer.Stringify(&valuePtr))
	fmt.Println(stringer.Stringify(valueInter))
	fmt.Println(stringer.Stringify(valueInterPtr))
	fmt.Println(stringer.Stringify(&valueInterPtr))

	// Output:
	// <nil>
	//
	// [one two three]
	// [one two three]
	// [one two three]
	// [one two three]
	// [one two three]
	// [one two three]
}

func ExampleStringify_nested() {
	second := "second"
	t := &tester{
		one: "first",
		two: &second,
		nested: &tester{
			one:    "",
			two:    &second,
			nested: nil,
		},
	}

	fmt.Println(stringer.Stringify(t))

	// Output:
	// {one:first two:second nested:{one: two:second nested:<nil>}}
}
