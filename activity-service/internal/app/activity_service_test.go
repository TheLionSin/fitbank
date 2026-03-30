package app

import (
	"context"
	"errors"
	"fitbank/activity-service/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ActivityRepositoryMock struct {
	mock.Mock
}

func (m *ActivityRepositoryMock) Create(ctx context.Context, act domain.Activity) error {
	args := m.Called(ctx, act)
	return args.Error(0)
}

func (m *ActivityRepositoryMock) FetchAll(ctx context.Context) ([]domain.Activity, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Activity), args.Error(1)
}

func (m *ActivityRepositoryMock) GetByID(ctx context.Context, id string) (domain.Activity, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Activity), args.Error(1)
}

func (m *ActivityRepositoryMock) Update(ctx context.Context, act domain.Activity) error {
	args := m.Called(ctx, act)
	return args.Error(0)
}

func (m *ActivityRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_Create(t *testing.T) {

	ctx := context.Background()

	t.Run("success_creation", func(t *testing.T) {

		repo := new(ActivityRepositoryMock)
		service := NewService(repo)

		input := domain.Activity{
			Type:   "Bench press",
			Weight: 100,
			Reps:   10,
		}

		repo.On("Create", ctx, mock.MatchedBy(func(a domain.Activity) bool {
			return a.Type == input.Type && a.Weight == input.Weight
		})).Return(nil).Once()

		result, err := service.Create(ctx, input)

		assert.NoError(t, err)
		assert.NotEmpty(t, result.ID)
		assert.NotZero(t, result.CreatedAt)

		repo.AssertExpectations(t)
	})

	t.Run("validation_error_empty_type", func(t *testing.T) {

		repo := new(ActivityRepositoryMock)
		service := NewService(repo)

		input := domain.Activity{
			Type:   "",
			Weight: 50,
			Reps:   5,
		}

		_, err := service.Create(ctx, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "activity type cannot be empty")

		repo.AssertNotCalled(t, "Create", ctx, mock.Anything)
	})
}

func TestService_Update(t *testing.T) {
	repo := new(ActivityRepositoryMock)
	service := NewService(repo)
	ctx := context.Background()

	t.Run("update_not_found", func(t *testing.T) {
		act := domain.Activity{
			ID:     "unknown-id",
			Type:   "Squats",
			Weight: 80,
			Reps:   12,
		}

		repo.On("Update", ctx, act).Return(errors.New("activity not found")).Once()

		err := service.Update(ctx, act)

		assert.Error(t, err)
		assert.Equal(t, "activity not found", err.Error())
		repo.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {

	ctx := context.Background()

	t.Run("success_delete", func(t *testing.T) {
		repo := new(ActivityRepositoryMock)
		service := NewService(repo)

		testID := "activity-123"

		repo.On("Delete", ctx, testID).Return(nil).Once()

		err := service.Delete(ctx, testID)

		assert.NoError(t, err)
		repo.AssertExpectations(t)

	})

	t.Run("delete_not_found", func(t *testing.T) {
		repo := new(ActivityRepositoryMock)
		service := NewService(repo)

		testID := "non-existing-id"

		repo.On("Delete", ctx, testID).Return(errors.New("activity not found")).Once()

		err := service.Delete(ctx, testID)
		assert.Error(t, err)
		assert.Equal(t, "activity not found", err.Error())
		repo.AssertExpectations(t)
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success_get_by_id", func(t *testing.T) {
		repo := new(ActivityRepositoryMock)
		service := NewService(repo)
		testID := "activity-123"

		repo.On("GetByID", ctx, testID).Return(domain.Activity{ID: testID}, nil).Once()

		result, err := service.GetByID(ctx, testID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, testID, result.ID)
		repo.AssertExpectations(t)
	})

	t.Run("failure_get_by_id", func(t *testing.T) {
		repo := new(ActivityRepositoryMock)
		service := NewService(repo)
		testID := "non-existing-id"

		repo.On("GetByID", ctx, testID).Return(domain.Activity{}, errors.New("not found"))

		result, err := service.GetByID(ctx, testID)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "not found")
		repo.AssertExpectations(t)
	})
}

func TestService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success_list", func(t *testing.T) {
		repo := new(ActivityRepositoryMock)
		service := NewService(repo)

		mockActivities := []domain.Activity{
			{ID: "1", Type: "Жим лежа", Weight: 100, Reps: 10},
			{ID: "2", Type: "Приседания", Weight: 120, Reps: 8},
		}

		repo.On("FetchAll", ctx).Return(mockActivities, nil).Once()

		result, err := service.FetchAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Жим лежа", result[0].Type)
		assert.Equal(t, 120.0, result[1].Weight)

		repo.AssertExpectations(t)
	})
}

func TestService_FetchAll_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(ActivityRepositoryMock)
	service := NewService(repo)

	dbError := errors.New("db error")
	repo.On("FetchAll", ctx).Return([]domain.Activity{}, dbError).Once()

	result, err := service.FetchAll(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Empty(t, result)

	repo.AssertExpectations(t)
}
