package service

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

var ErrTimeout = errors.New("timeout")

type InChanT struct {
	Flag bool
	In   chan int
}

type OutChanT chan int

type ServiceHandlerT func(context.Context, InChanT, OutChanT, string)

func NewHandler() ServiceHandlerT {
	return func(ctx context.Context, in InChanT, out OutChanT, name string) {
		val := rand.Intn(100) + 1
		if in.Flag {
			select {
			case <-ctx.Done():
				return
			case val = <-in.In:
				break
			}
		}
		r := rand.Intn(3) + 1
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(r) * time.Second):
			break
		}
		out <- val * r
	}
}
