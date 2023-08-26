package postgres

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
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

func (u *UserSegmentRepo) AddUserToSegment(ctx context.Context, segments models.AddUserToSegment) (*models.AddUserToSegmentResponse, error) {
	response := new(models.AddUserToSegmentResponse)
	response.UserID = segments.UserID

	checkUserQuery := "SELECT id FROM users WHERE id = $1"
	err := u.db.QueryRowContext(ctx, checkUserQuery, segments.UserID).Scan(&segments.UserID)
	if err != nil {
		return nil, models.UserNotFound
	}

	addQuery := "INSERT INTO user_segments (user_id, segment_name) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	checkSegmentQuery := "SELECT segment_name FROM segments WHERE segment_name = $1"

	for _, segmentName := range segments.SegmentsToAdd {
		err := u.db.QueryRowContext(ctx, checkSegmentQuery, segmentName).Scan(&segmentName)
		if err != nil {
			if err == sql.ErrNoRows {
				response.NotExistSegments = append(response.NotExistSegments, segmentName)
				continue
			}
		}

		_, err = u.db.ExecContext(ctx, addQuery, segments.UserID, segmentName)
		if err != nil {
			logrus.Println("error adding user to segment: ", err)
			return nil, err
		}
		response.AddedSegments = append(response.AddedSegments, segmentName)
	}

	deleteQuery := "DELETE FROM user_segments WHERE user_id = $1 AND segment_name = $2"
	for _, segmentName := range segments.SegmentsToDelete {
		err = u.db.QueryRowContext(ctx, checkSegmentQuery, segmentName).Scan(&segmentName)
		if err != nil {
			if err == sql.ErrNoRows {
				response.NotExistSegments = append(response.NotExistSegments, segmentName)
				continue
			}
		}

		_, err = u.db.ExecContext(ctx, deleteQuery, segments.UserID, segmentName)
		if err != nil {
			return nil, err
		}
		response.DeletedSegments = append(response.DeletedSegments, segmentName)
	}
	return response, nil
}

func (u *UserSegmentRepo) GetActiveUserSegments(ctx context.Context, userID int) ([]string, error) {
	checkUserQuery := "SELECT id FROM users WHERE id = $1"
	err := u.db.QueryRowContext(ctx, checkUserQuery, userID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.UserNotFound
		}
		return nil, err
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
