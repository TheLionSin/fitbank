package app

import (
	"context"
	"fitbank/activity-service/internal/domain"
	activitypb "fitbank/activity-service/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ActivityGRPCServer struct {
	activitypb.UnimplementedActivityServiceServer
	service ActivityUseCase
}

func NewGRPCServer(service ActivityUseCase) *ActivityGRPCServer {
	return &ActivityGRPCServer{service: service}
}

func (s *ActivityGRPCServer) CreateActivity(ctx context.Context, req *activitypb.CreateActivityRequest) (*activitypb.ActivityResponse, error) {
	act, err := s.service.Create(ctx, domain.Activity{
		Type:   req.Type,
		Weight: req.Weight,
		Reps:   int(req.Reps),
	})

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid data: %v", err)

	}

	return &activitypb.ActivityResponse{
		Id:        act.ID,
		Type:      act.Type,
		Weight:    act.Weight,
		Reps:      int32(act.Reps),
		CreatedAt: timestamppb.New(act.CreatedAt),
	}, nil
}
