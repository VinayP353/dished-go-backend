package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yourusername/dished-go-backend/internal/models"
	"github.com/yourusername/dished-go-backend/internal/repository"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUser(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type userService struct {
	repo  repository.UserRepository
	redis *redis.Client
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *userService) CreateUser(user *models.User) error {
	return s.repo.Create(user)
}

func (s *userService) GetUser(id uint) (*models.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", id)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	userJSON, _ := json.Marshal(user)
	s.redis.Set(ctx, cacheKey, userJSON, 10*time.Minute)

	return user, nil
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) UpdateUser(user *models.User) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	s.redis.Del(ctx, cacheKey)
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", id)
	s.redis.Del(ctx, cacheKey)
	return s.repo.Delete(id)
}
