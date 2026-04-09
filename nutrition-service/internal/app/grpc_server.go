package app

import (
	"context"
	"fitbank/nutrition-service/internal/domain"
	nutritionpb "fitbank/nutrition-service/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NutritionGRPCServer struct {
	// Встраиваем нереализованный сервер для обратной совместимости.
	// Если мы добавим метод в .proto, но не напишем его тут, сервер не упадет.
	nutritionpb.UnimplementedNutritionServiceServer

	// Внедряем бизнес-логику (Service Layer)
	service NutritionUseCase
}

func NewNutritionGRPCServer(service NutritionUseCase) *NutritionGRPCServer {
	return &NutritionGRPCServer{
		service: service,
	}
}

// AddMeal — обработчик gRPC запроса на добавление еды
func (s *NutritionGRPCServer) AddMeal(ctx context.Context, req *nutritionpb.AddMealRequest) (*nutritionpb.MealResponse, error) {
	// 1. Маппинг из gRPC-запроса в доменную модель нашего сервиса
	meal := domain.Meal{
		Name:          req.GetName(),
		Calories:      req.GetCalories(),
		Proteins:      req.GetProteins(),
		Fats:          req.GetFats(),
		Carbohydrates: req.GetCarbohydrates(),
		WeightGrams:   int(req.GetWeightGrams()),
	}

	result, err := s.service.AddMeal(ctx, meal)
	if err != nil {
		// В gRPC мы не возвращаем просто ошибку, мы возвращаем статус-код.
		// Если это ошибка валидации — возвращаем InvalidArgument.
		return nil, status.Errorf(codes.InvalidArgument, "could not add meal: %v", err)
	}

	return &nutritionpb.MealResponse{
		Id:            result.ID,
		Name:          result.Name,
		Calories:      result.Calories,
		Proteins:      result.Proteins,
		Fats:          result.Fats,
		Carbohydrates: result.Carbohydrates,
		WeightGrams:   int32(result.WeightGrams),
		// Превращаем стандартное время Go в формат Protobuf Timestamp
		CreatedAt: timestamppb.New(result.CreatedAt),
	}, nil
}

func (s *NutritionGRPCServer) GetDailyReport(ctx context.Context, req *nutritionpb.DailyReportRequest) (*nutritionpb.DailyReportResponse, error) {
	report, err := s.service.GetDailyReport(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate report: %v", err)
	}

	return &nutritionpb.DailyReportResponse{
		TotalCalories: report.TotalCalories,
		TotalProteins: report.TotalProteins,
		TotalFats:     report.TotalFats,
		TotalCarbs:    report.TotalCarbs,
	}, nil
}
