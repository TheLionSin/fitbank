package main

import (
	"context"
	activitypb "fitbank/activity-service/pkg/api"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := activitypb.NewActivityServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := c.CreateActivity(ctx, &activitypb.CreateActivityRequest{
		Type:   "Bench press(gRPC)",
		Weight: 100,
		Reps:   8,
	})
	if err != nil {
		log.Fatalf("could not create activity: %v", err)
	}

	log.Printf("Activity created via gRPC! ID: %s\n", r.GetId())
}
