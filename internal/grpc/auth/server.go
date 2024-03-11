package auth

import (
	"context"
	"errors"

	ssov1 "github.com/LDmitryLD/protos_less/gen/go/sso"
	"github.com/LDmitryLD/sso/sso/internal/services/auth"
	"github.com/LDmitryLD/sso/sso/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type Auth interface {
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	atuh Auth
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{atuh: auth})
}

func (s *serverAPI) Register(ctx context.Context, in *ssov1.RegsiterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(in); err != nil {
		return nil, err
	}

	userID, err := s.atuh.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(in); err != nil {
		return nil, err
	}

	token, err := s.atuh.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentails) {
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(in); err != nil {
		return nil, err
	}

	isAdmin, err := s.atuh.IsAdmin(ctx, in.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(in *ssov1.LoginRequest) error {
	if len(in.GetEmail()) < 5 {
		return status.Error(codes.InvalidArgument, "email is invalid")
	}

	if in.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if in.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func validateRegister(in *ssov1.RegsiterRequest) error {
	if len(in.GetEmail()) < 5 {
		return status.Error(codes.InvalidArgument, "email is invalid")
	}

	if in.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateIsAdmin(in *ssov1.IsAdminRequest) error {
	if in.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
