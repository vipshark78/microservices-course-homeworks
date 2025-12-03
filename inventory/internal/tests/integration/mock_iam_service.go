package integration

import (
	"context"
	"net"

	"google.golang.org/grpc"

	authV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/auth/v1"
	commonV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/common/v1"
)

type mockIAMService struct {
	authV1.UnimplementedAuthServiceServer
}

func (m *mockIAMService) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	return &authV1.LoginResponse{
		SessionUuid: "mock-session-uuid",
	}, nil
}

func (m *mockIAMService) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	return &authV1.WhoamiResponse{
		User: &commonV1.User{
			Uuid: "mock-user-id",
			Info: &commonV1.UserInfo{
				Login: "test-user",
				Email: "test@example.com",
			},
		},
	}, nil
}

type MockIAMContainer struct {
	server   *grpc.Server
	listener net.Listener
	address  string
}

func NewMockIAMContainer(address string) *MockIAMContainer {
	return &MockIAMContainer{
		address: address,
	}
}

func (m *MockIAMContainer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", m.address)
	if err != nil {
		return err
	}
	m.listener = listener

	m.server = grpc.NewServer()
	authV1.RegisterAuthServiceServer(m.server, &mockIAMService{})

	go func() {
		if err := m.server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (m *MockIAMContainer) Stop() error {
	if m.server != nil {
		m.server.GracefulStop()
	}
	if m.listener != nil {
		return m.listener.Close()
	}
	return nil
}

func (m *MockIAMContainer) Address() string {
	return m.address
}
