package service

import (
	"context"
)

type Service struct {
	name           string
	serviceHandler ServiceHandlerT
}

func New(name string, handler ServiceHandlerT) *Service {
	return &Service{
		name:           name,
		serviceHandler: handler,
	}
}

func (s *Service) CallHandler(ctx context.Context, in InChanT) (OutChanT, error) {
	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	default:
		break
	}
	out := make(chan int)
	go s.serviceHandler(ctx, in, out, s.name)
	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	default:
		break
	}

	return out, nil
}

func (s Service) GetName() string {
	return s.name
}

