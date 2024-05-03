package http

import "context"

type Auth interface {
	LoginUser(ctx context.Context, username string, password string) (string, error)
	RegisterUser(ctx context.Context, username string, password string) (int, error)
	LoginAdmin(ctx context.Context, username string, password string) (string, error)
	RegisterAdmin(ctx context.Context, username, password string, registrantID, roleID int) (int, error)
}
