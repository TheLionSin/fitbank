package app

import (
	"context"
	"fitbank/activity-service/internal/domain"
	activitypb "fitbank/activity-service/pkg/api"

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
	act := domain.Activity{
		Type:   req.Type,
		Weight: req.Weight,
		Reps:   int(req.Reps),
	}

	if err := act.Validate(); err != nil {
		return nil, err
	}

	createdAct, err := s.service.Create(ctx, act)
	if err != nil {
		return nil, err
	}

	return &activitypb.ActivityResponse{
		Id:        createdAct.ID,
		Type:      createdAct.Type,
		Weight:    createdAct.Weight,
		Reps:      int32(createdAct.Reps),
		CreatedAt: timestamppb.New(createdAct.CreatedAt),
	}, nil
}
