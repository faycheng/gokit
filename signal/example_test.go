package signal

import (
	"context"
	"fmt"
)

func sayHello(c context.Context, names ...interface{}) error {
	for _, nameIface := range names {
		name := nameIface.(string)
		fmt.Println("hello,", name)
	}
	return nil
}

func ExampleSignal() {
	s := NewSignal()
	s.Connect(sayHello)
	s.Send(context.TODO(), "tony")
	// Output: hello, tony
}
