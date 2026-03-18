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

type ChefProfileService interface {
	CreateProfile(req *models.CreateChefProfileRequest) (*models.ChefProfile, error)
	GetProfile(id uint) (*models.ChefProfile, error)
	GetAllProfiles() ([]models.ChefProfile, error)
	UpdateProfile(id uint, req *models.UpdateChefProfileRequest) (*models.ChefProfile, error)
	DeleteProfile(id uint) error
}

type chefProfileService struct {
	repo  repository.ChefProfileRepository
	redis *redis.Client
}

func NewChefProfileService(repo repository.ChefProfileRepository, redisClient *redis.Client) ChefProfileService {
	return &chefProfileService{repo: repo, redis: redisClient}
}

func (s *chefProfileService) CreateProfile(req *models.CreateChefProfileRequest) (*models.ChefProfile, error) {
	profile := &models.ChefProfile{
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

func (s *chefProfileService) GetProfile(id uint) (*models.ChefProfile, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("chef_profile:%d", id)

	if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
		var profile models.ChefProfile
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

func (s *chefProfileService) GetAllProfiles() ([]models.ChefProfile, error) {
	return s.repo.GetAll()
}

func (s *chefProfileService) UpdateProfile(id uint, req *models.UpdateChefProfileRequest) (*models.ChefProfile, error) {
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
	s.redis.Del(context.Background(), fmt.Sprintf("chef_profile:%d", id))
	return profile, nil
}

func (s *chefProfileService) DeleteProfile(id uint) error {
	s.redis.Del(context.Background(), fmt.Sprintf("chef_profile:%d", id))
	return s.repo.Delete(id)
}
