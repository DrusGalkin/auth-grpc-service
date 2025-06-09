package grpc_client

import (
	"fmt"
	pk "github.com/DrusGalkin/auth-protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pk.AuthClient
}

func NewClient(address string) (*Client, error) {
	const op = "grpc_client.NewClient"

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		pk.NewAuthClient(conn),
	}, nil
}
