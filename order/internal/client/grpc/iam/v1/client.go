package v1

import (
	"context"

	"google.golang.org/grpc"

	authV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/auth/v1"
)

type client struct {
	iamClient authV1.AuthServiceClient
}

func NewClient(iamClient authV1.AuthServiceClient) *client {
	return &client{
		iamClient: iamClient,
	}
}

func (c *client) Login(ctx context.Context, in *authV1.LoginRequest, opts ...grpc.CallOption) (*authV1.LoginResponse, error) {
	return c.iamClient.Login(ctx, in)
}

func (c *client) Whoami(ctx context.Context, in *authV1.WhoamiRequest, opts ...grpc.CallOption) (*authV1.WhoamiResponse, error) {
	return c.iamClient.Whoami(ctx, in)
}
