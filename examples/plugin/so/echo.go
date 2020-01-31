package main

import (
	"context"
	"fmt"
)

func Echo(ctx context.Context, req interface{}) (reply interface{}, err error) {
	fmt.Println("hello world")
	return nil, nil
}
