package http

import "context"

type Auth interface {
	Login(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, username string, password string) (int, error)
}
