package service

import (
	"context"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

type UserSegmentRepository interface {
	CreateSegment(ctx context.Context, segmentName string) error
	DeleteSegment(ctx context.Context, segmentName string) error
	AddUserToSegment(ctx context.Context, segments models.AddUserToSegment) (*models.AddUserToSegmentResponse, error)
	GetActiveUserSegments(ctx context.Context, userID int) ([]string, error)
}

type UserSegment struct {
	repo UserSegmentRepository
}

func NewUserSegment(repo UserSegmentRepository) *UserSegment {
	return &UserSegment{
		repo: repo,
	}
}

func (u *UserSegment) CreateSegment(ctx context.Context, segmentName string) error {
	return u.repo.CreateSegment(ctx, segmentName)
}

func (u *UserSegment) DeleteSegment(ctx context.Context, segmentName string) error {
	return u.repo.DeleteSegment(ctx, segmentName)
}

func (u *UserSegment) AddUserToSegments(ctx context.Context, segments models.AddUserToSegment) (models.AddUserToSegmentResponse, error) {
	res, err := u.repo.AddUserToSegment(ctx, segments)
	return *res, err
}

func (u *UserSegment) GetActiveUserSegments(ctx context.Context, userID int) ([]string, error) {
	return u.repo.GetActiveUserSegments(ctx, userID)
}
