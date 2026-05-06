package handler

import (
	"context"
	"user-service/internal/service"
	pb "user-service/proto"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	err := h.service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Message: "register success",
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}