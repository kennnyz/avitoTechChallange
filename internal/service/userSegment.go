package service

import (
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

type UserSegmentRepository interface {
	CreateSegment(ctx context.Context, segmentName string) error
	DeleteSegment(ctx context.Context, segmentName string) error
	AddUserToSegment(ctx context.Context, userID int, segments []string) ([]string, error)
	DeleteUserFromSegments(ctx context.Context, userId int, segments []string) ([]string, error)
	GetActiveUserSegments(ctx context.Context, userID int) ([]string, error)
	CheckSegment(ctx context.Context, segmentName string) error
	CheckUser(ctx context.Context, userID int) error
	GetUsersHistory(ctx context.Context, year, month string) ([]models.UserHistory, error)
}

type UserSegment struct {
	repo            UserSegmentRepository
	historyFilePath string
}

func NewUserSegment(repo UserSegmentRepository, tmpFilePath string) *UserSegment {
	return &UserSegment{
		repo:            repo,
		historyFilePath: tmpFilePath,
	}
}

func (u *UserSegment) CreateSegment(ctx context.Context, segmentName string) error {
	return u.repo.CreateSegment(ctx, segmentName)
}

func (u *UserSegment) DeleteSegment(ctx context.Context, segmentName string) error {
	err := u.repo.CheckSegment(ctx, segmentName)
	if err != nil {
		return err
	}

	return u.repo.DeleteSegment(ctx, segmentName)
}

func (u *UserSegment) AddUserToSegments(ctx context.Context, segments models.AddUserToSegment) (*models.AddUserToSegmentResponse, error) {
	err := u.repo.CheckUser(ctx, segments.UserID)
	if err != nil {
		return nil, err
	}

	res := &models.AddUserToSegmentResponse{}
	res.UserID = segments.UserID
	added, err := u.repo.AddUserToSegment(ctx, segments.UserID, segments.SegmentsToAdd)
	if err != nil {
		return nil, err
	}

	removed, err := u.repo.DeleteUserFromSegments(ctx, segments.UserID, segments.SegmentsToDelete)
	if err != nil {
		return nil, err
	}

	res.AddedSegments = added
	res.DeletedSegments = removed

	return res, nil
}

func (u *UserSegment) GetActiveUserSegments(ctx context.Context, userID int) ([]string, error) {
	err := u.repo.CheckUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return u.repo.GetActiveUserSegments(ctx, userID)
}

func (u *UserSegment) GetHistoryFile(ctx context.Context, year, month string) (string, error) {
	filePath := u.historyFilePath + year + "-" + month + ".csv"
	// check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return filepath.Abs(filePath)
	}

	history, err := u.repo.GetUsersHistory(ctx, year, month)
	if err != nil {
		return "", err
	}

	if len(history) == 0 {
		return "", models.NoDataErr
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write header
	err = writer.Write([]string{"userid", "segment", "action", "date"})
	if err != nil {
		return "", err
	}

	for _, h := range history {
		strID := strconv.Itoa(h.UserID)
		err := writer.Write([]string{strID, h.Segment, h.Action, h.Date})
		if err != nil {
			return "", err
		}
	}

	return filepath.Abs(filePath)
}
