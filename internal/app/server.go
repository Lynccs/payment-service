package app

import "context"

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}
