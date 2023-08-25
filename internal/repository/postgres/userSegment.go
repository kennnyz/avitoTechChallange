package postgres

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"strings"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

type UserSegmentRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserSegmentRepo {
	return &UserSegmentRepo{
		db: db,
	}
}

func (u *UserSegmentRepo) CreateSegment(ctx context.Context, segmentName string) error {
	checkQuery := "SELECT segment_name FROM segments WHERE segment_name = $1"
	err := u.db.QueryRowContext(ctx, checkQuery, segmentName).Scan(&segmentName)
	if err == nil {
		logrus.Println("segment already exists")
		return nil
	}

	insertQuery := "INSERT INTO segments (segment_name) VALUES ($1)"
	_, err = u.db.ExecContext(ctx, insertQuery, segmentName)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserSegmentRepo) DeleteSegment(ctx context.Context, segmentName string) error {
	deleteQuery := "DELETE FROM segments WHERE segment_name = $1"
	_, err := u.db.ExecContext(ctx, deleteQuery, segmentName)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserSegmentRepo) AddUserToSegment(ctx context.Context, userID int, segmentNamesToDelete, segmentNamesToAdd []string) error {
	checkUserQuery := "SELECT id FROM users WHERE id = $1"
	err := u.db.QueryRowContext(ctx, checkUserQuery, userID).Scan(&userID)
	if err != nil {
		return models.UserNotFound
	}

	addQuery := "INSERT INTO user_segments (user_id, segment_name) VALUES ($1, $2)"
	for _, segmentName := range segmentNamesToAdd {
		_, err := u.db.ExecContext(ctx, addQuery, userID, segmentName)
		if err != nil {
			// if user already exists in segment, ignore
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				// Это ошибка, которую хотим проигнорировать
				continue
			}
			logrus.Println("error adding user to segment: ", err)
			return err
		}
	}

	deleteQuery := "DELETE FROM user_segments WHERE user_id = $1 AND segment_name = $2"
	for _, segmentName := range segmentNamesToDelete {
		_, err := u.db.ExecContext(ctx, deleteQuery, userID, segmentName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *UserSegmentRepo) GetActiveUserSegments(ctx context.Context, userID int) ([]string, error) {
	checkUserQuery := "SELECT id FROM users WHERE id = $1"
	err := u.db.QueryRowContext(ctx, checkUserQuery, userID).Scan(&userID)
	if err != nil {
		return nil, models.UserNotFound
	}

	query := "SELECT segment_name FROM user_segments WHERE user_id = $1"
	rows, err := u.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segmentNames []string
	for rows.Next() {
		var segmentName string
		err := rows.Scan(&segmentName)
		if err != nil {
			return nil, err
		}
		segmentNames = append(segmentNames, segmentName)
	}
	return segmentNames, nil
}
