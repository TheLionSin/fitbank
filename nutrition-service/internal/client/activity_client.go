package client

import (
	"context"
	activitypb "fitbank/activity-service/pkg/api"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ActivityClientInterface interface {
	GetTotalCaloriesBurned(ctx context.Context) (float64, error)
}

type ActivityClient struct {
	client activitypb.ActivityServiceClient
}

func NewActivityClient(target string) (*ActivityClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("could not connect to activity-service: %v", err)
	}

	return &ActivityClient{
		client: activitypb.NewActivityServiceClient(conn),
	}, nil
}

func (c *ActivityClient) GetTotalCaloriesBurned(ctx context.Context) (float64, error) {
	resp, err := c.client.ListActivities(ctx, &activitypb.ListActivitiesRequest{})
	if err != nil {
		return 0, fmt.Errorf("error calling activity-service: %v", err)
	}

	var total float64
	for _, act := range resp.Activities {
		// Допустим, каждая активность — это какой-то расход.
		// Пока для примера считаем по простой формуле: вес * повторения / 10
		total += (act.Weight * float64(act.Reps)) / 10.0
	}

	return total, nil
}
