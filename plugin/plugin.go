package plugin

import (
	"context"
)

type Call func(ctx context.Context, req interface{}) (reply interface{}, err error)

type Plugin interface {
	Lookup(name string) (call Call, err error)
	String() string
	Close() error
}
