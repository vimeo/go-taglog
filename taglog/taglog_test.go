package taglog

import "os"

func ExampleCopy() {
	orig := New(os.Stdout, "A ", 0)
	orig.Println("a")
	copy := orig.Copy()
	orig.Println("b")
	copy.SetPrefix("B ")
	orig.Println("c")
	copy.Println("c")
	// Output:
	// A a
	// A b
	// A c
	// B c
}
