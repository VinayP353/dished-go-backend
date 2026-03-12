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

type CookProfileService interface {
	CreateProfile(req *models.CreateCookProfileRequest) (*models.CookProfile, error)
	GetProfile(id uint) (*models.CookProfile, error)
	GetAllProfiles() ([]models.CookProfile, error)
	UpdateProfile(id uint, req *models.UpdateCookProfileRequest) (*models.CookProfile, error)
	DeleteProfile(id uint) error
}

type cookProfileService struct {
	repo  repository.CookProfileRepository
	redis *redis.Client
}

func NewCookProfileService(repo repository.CookProfileRepository, redisClient *redis.Client) CookProfileService {
	return &cookProfileService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *cookProfileService) CreateProfile(req *models.CreateCookProfileRequest) (*models.CookProfile, error) {
	profile := &models.CookProfile{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		PreferredName: req.PreferredName,
		Address:       req.Address,
		Description:   req.Description,
		Verified:      false,
	}

	if err := s.repo.Create(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *cookProfileService) GetProfile(id uint) (*models.CookProfile, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("cook_profile:%d", id)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var profile models.CookProfile
		if err := json.Unmarshal([]byte(cached), &profile); err == nil {
			return &profile, nil
		}
	}

	profile, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	profileJSON, _ := json.Marshal(profile)
	s.redis.Set(ctx, cacheKey, profileJSON, 10*time.Minute)

	return profile, nil
}

func (s *cookProfileService) GetAllProfiles() ([]models.CookProfile, error) {
	return s.repo.GetAll()
}

func (s *cookProfileService) UpdateProfile(id uint, req *models.UpdateCookProfileRequest) (*models.CookProfile, error) {
	profile, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		profile.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		profile.LastName = *req.LastName
	}
	if req.PreferredName != nil {
		profile.PreferredName = *req.PreferredName
	}
	if req.Address != nil {
		profile.Address = *req.Address
	}
	if req.Description != nil {
		profile.Description = *req.Description
	}
	if req.Verified != nil {
		profile.Verified = *req.Verified
	}

	if err := s.repo.Update(profile); err != nil {
		return nil, err
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("cook_profile:%d", id)
	s.redis.Del(ctx, cacheKey)

	return profile, nil
}

func (s *cookProfileService) DeleteProfile(id uint) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("cook_profile:%d", id)
	s.redis.Del(ctx, cacheKey)
	return s.repo.Delete(id)
}
