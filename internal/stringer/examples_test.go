package stringer_test

import (
	"fmt"

	"github.com/VirtusLab/kubedrainer/internal/stringer"
)

func ExampleStringify() {
	var value = []string{
		"one", "two", "three",
	}
	var valueInter interface{}
	valueInter = value
	fmt.Println(stringer.Stringify(nil))
	fmt.Println(stringer.Stringify(""))
	fmt.Println(stringer.Stringify(value))
	fmt.Println(stringer.Stringify(&value))
	fmt.Println(stringer.Stringify(valueInter))
	fmt.Println(stringer.Stringify(&valueInter))

	// Output:
	// <nil>
	//
	// [one two three]
	// [one two three]
	// [one two three]
	// [one two three]
}
