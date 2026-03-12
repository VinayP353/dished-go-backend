package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yourusername/dished-go-backend/internal/models"
	"github.com/yourusername/dished-go-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type CookService interface {
	CreateCook(req *models.CreateCookRequest) (*models.Cook, error)
	GetCook(id uint) (*models.Cook, error)
	GetAllCooks() ([]models.Cook, error)
	UpdateCook(id uint, req *models.UpdateCookRequest) (*models.Cook, error)
	DeleteCook(id uint) error
	Login(username, password string) (*models.Cook, error)
}

type cookService struct {
	repo  repository.CookRepository
	redis *redis.Client
}

func NewCookService(repo repository.CookRepository, redisClient *redis.Client) CookService {
	return &cookService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *cookService) validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) < 6 {
		return errors.New("password must contain at least 6 letters")
	}

	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)

	if !hasSpecial {
		return errors.New("password must contain at least 1 special character")
	}
	if !hasNumber {
		return errors.New("password must contain at least 1 number")
	}
	if !hasLetter {
		return errors.New("password must contain at least 6 letters")
	}

	return nil
}

func (s *cookService) CreateCook(req *models.CreateCookRequest) (*models.Cook, error) {
	if err := s.validatePassword(req.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	cook := &models.Cook{
		Username:         req.Username,
		Password:         string(hashedPassword),
		Email:            req.Email,
		Status:           "active",
		CookieExpiration: 24 * time.Hour,
	}

	if err := s.repo.Create(cook); err != nil {
		return nil, err
	}

	return cook, nil
}

func (s *cookService) GetCook(id uint) (*models.Cook, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("cook:%d", id)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var cook models.Cook
		if err := json.Unmarshal([]byte(cached), &cook); err == nil {
			return &cook, nil
		}
	}

	cook, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cookJSON, _ := json.Marshal(cook)
	s.redis.Set(ctx, cacheKey, cookJSON, 10*time.Minute)

	return cook, nil
}

func (s *cookService) GetAllCooks() ([]models.Cook, error) {
	return s.repo.GetAll()
}

func (s *cookService) UpdateCook(id uint, req *models.UpdateCookRequest) (*models.Cook, error) {
	cook, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Username != nil {
		cook.Username = *req.Username
	}
	if req.Email != nil {
		cook.Email = *req.Email
	}
	if req.Status != nil {
		cook.Status = *req.Status
	}

	if err := s.repo.Update(cook); err != nil {
		return nil, err
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("cook:%d", id)
	s.redis.Del(ctx, cacheKey)

	return cook, nil
}

func (s *cookService) DeleteCook(id uint) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("cook:%d", id)
	s.redis.Del(ctx, cacheKey)
	return s.repo.Delete(id)
}

func (s *cookService) Login(username, password string) (*models.Cook, error) {
	cook, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cook.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	cook.LastLogin = &now
	s.repo.Update(cook)

	return cook, nil
}
