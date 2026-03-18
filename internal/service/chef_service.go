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

type ChefService interface {
	Register(req *models.RegisterRequest) (*models.Chef, error)
	Login(username, password string) (*models.Chef, error)
	GetChef(id uint) (*models.Chef, error)
	GetAllChefs() ([]models.Chef, error)
	GetAllUsernames() ([]string, error)
	UpdateChef(id uint, req *models.UpdateChefRequest) (*models.Chef, error)
	UpdateProfile(chefID uint, req *models.UpdateChefProfileRequest) (*models.Chef, error)
	DeleteChef(id uint) error
}

type chefService struct {
	repo        repository.ChefRepository
	profileRepo repository.ChefProfileRepository
	redis       *redis.Client
}

func NewChefService(repo repository.ChefRepository, profileRepo repository.ChefProfileRepository, redisClient *redis.Client) ChefService {
	return &chefService{repo: repo, profileRepo: profileRepo, redis: redisClient}
}

func (s *chefService) validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return errors.New("password must contain at least 1 special character")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least 1 number")
	}
	letters := regexp.MustCompile(`[a-zA-Z]`).FindAllString(password, -1)
	if len(letters) < 6 {
		return errors.New("password must contain at least 6 letters")
	}
	return nil
}

func (s *chefService) Register(req *models.RegisterRequest) (*models.Chef, error) {
	if err := s.validatePassword(req.Password); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}
	if _, err := s.repo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	chef := &models.Chef{
		Username:         req.Username,
		Password:         string(hashedPassword),
		Email:            req.Email,
		Status:           "active",
		CookieExpiration: int64(24 * time.Hour),
		ChefProfile: &models.ChefProfile{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Verified:  false,
		},
	}

	if err := s.repo.Create(chef); err != nil {
		return nil, err
	}
	return chef, nil
}

func (s *chefService) Login(username, password string) (*models.Chef, error) {
	chef, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(chef.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	now := time.Now()
	chef.LastLogin = &now
	s.repo.Update(chef)
	return chef, nil
}

func (s *chefService) GetChef(id uint) (*models.Chef, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("chef:%d", id)

	if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
		var chef models.Chef
		if err := json.Unmarshal([]byte(cached), &chef); err == nil {
			return &chef, nil
		}
	}

	chef, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	chefJSON, _ := json.Marshal(chef)
	s.redis.Set(ctx, cacheKey, chefJSON, 10*time.Minute)
	return chef, nil
}

func (s *chefService) GetAllChefs() ([]models.Chef, error) {
	return s.repo.GetAll()
}

func (s *chefService) GetAllUsernames() ([]string, error) {
	chefs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	usernames := make([]string, len(chefs))
	for i, chef := range chefs {
		usernames[i] = chef.Username
	}
	return usernames, nil
}

func (s *chefService) UpdateChef(id uint, req *models.UpdateChefRequest) (*models.Chef, error) {
	chef, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cols := map[string]interface{}{}
	if req.Username != nil {
		cols["username"] = *req.Username
		chef.Username = *req.Username
	}
	if req.Email != nil {
		cols["email"] = *req.Email
		chef.Email = *req.Email
	}
	if req.Status != nil {
		cols["status"] = *req.Status
		chef.Status = *req.Status
	}

	if len(cols) > 0 {
		if err := s.repo.UpdateColumns(chef, cols); err != nil {
			return nil, err
		}
	}

	s.redis.Del(context.Background(), fmt.Sprintf("chef:%d", id))
	return chef, nil
}

func (s *chefService) UpdateProfile(chefID uint, req *models.UpdateChefProfileRequest) (*models.Chef, error) {
	chef, err := s.repo.GetByID(chefID)
	if err != nil {
		return nil, err
	}
	if chef.ChefProfile == nil {
		return nil, errors.New("chef profile not found")
	}
	if req.FirstName != nil {
		chef.ChefProfile.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		chef.ChefProfile.LastName = *req.LastName
	}
	if req.PreferredName != nil {
		chef.ChefProfile.PreferredName = *req.PreferredName
	}
	if req.Address != nil {
		chef.ChefProfile.Address = *req.Address
	}
	if req.Description != nil {
		chef.ChefProfile.Description = *req.Description
	}
	if req.Verified != nil {
		chef.ChefProfile.Verified = *req.Verified
	}
	if err := s.profileRepo.Update(chef.ChefProfile); err != nil {
		return nil, err
	}
	s.redis.Del(context.Background(), fmt.Sprintf("chef:%d", chefID))
	return chef, nil
}

func (s *chefService) DeleteChef(id uint) error {
	s.redis.Del(context.Background(), fmt.Sprintf("chef:%d", id))
	return s.repo.Delete(id)
}
