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

func (s *ActivityGRPCServer) ListActivities(ctx context.Context, req *activitypb.ListActivitiesRequest) (*activitypb.ListActivitiesResponse, error) {
	// 1. Получаем данные из бизнес-логики
	activities, err := s.service.FetchAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch activities: %v", err)
	}

	// 2. Маппим доменные модели в gRPC-ответ
	var pbActivities []*activitypb.ActivityResponse
	for _, a := range activities {
		pbActivities = append(pbActivities, &activitypb.ActivityResponse{
			Id:     a.ID,
			Type:   a.Type,
			Weight: a.Weight,
			Reps:   int32(a.Reps),
		})
	}
	return &activitypb.ListActivitiesResponse{
		Activities: pbActivities,
	}, nil
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
