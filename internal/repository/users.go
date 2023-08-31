package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

func (u *UserSegmentRepo) AddUserToSegment(ctx context.Context, userID int, segments []string) error {
	res := []string{}

	addQuery := "INSERT INTO user_segments (user_id, segment_name) VALUES ($1, $2) ON CONFLICT DO NOTHING"

	batch := &pgx.Batch{}
	for _, segmentName := range segments {
		batch.Queue(addQuery, userID, segmentName)
	}
	results := u.db.SendBatch(ctx, batch)
	defer results.Close()

	for _, segmentName := range segments {
		rw, err := results.Exec()
		if err != nil {
			logrus.Println("Error adding user to segment:", err)
			check := u.CheckSegment(ctx, segmentName)
			if check != nil {
				return check
			}
			return err
		}
		if rw.RowsAffected() == 0 {
			err = u.CheckSegment(ctx, segmentName)
			if err != nil {
				return err
			}
			continue
		}

		res = append(res, segmentName)
	}

	batch = &pgx.Batch{}
	for _, segmentName := range res {
		batch.Queue("INSERT INTO history (user_id, segment_name, action_type) VALUES ($1, $2, $3)", userID, segmentName, "added")
	}
	u.db.SendBatch(ctx, batch)

	return nil
}

func (u *UserSegmentRepo) DeleteUserFromSegments(ctx context.Context, userId int, segments []string) error {
	history := []string{}
	deleteQuery := "DELETE FROM user_segments WHERE user_id = $1 AND segment_name = $2"

	batch := &pgx.Batch{}
	for _, segmentName := range segments {
		batch.Queue(deleteQuery, userId, segmentName)
	}
	results := u.db.SendBatch(ctx, batch)
	defer results.Close()

	for _, segmentName := range segments {
		rw, err := results.Exec()
		if err != nil {
			logrus.Println("Error removing user from segment:", err)
			return err
		}
		if rw.RowsAffected() == 0 {
			continue
		}
		history = append(history, segmentName)
	}

	batch = &pgx.Batch{}
	for _, segmentName := range history {
		batch.Queue("INSERT INTO history (user_id, segment_name, action_type) VALUES ($1, $2, $3)", userId, segmentName, "deleted")
	}
	results = u.db.SendBatch(ctx, batch)
	for _, segmentName := range history {
		_, err := results.Exec()
		if err != nil {
			logrus.Println("Error adding to history:", err, " segment: ", segmentName)
			return err
		}
	}

	return nil
}

func (u *UserSegmentRepo) GetActiveUserSegments(ctx context.Context, userID int) ([]string, error) {
	query := "SELECT segment_name FROM user_segments WHERE user_id = $1"
	rows, err := u.db.Query(ctx, query, userID)
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

func (u *UserSegmentRepo) CheckUser(ctx context.Context, userID int) error {
	checkUserQuery := "SELECT id FROM users WHERE id = $1"
	err := u.db.QueryRow(ctx, checkUserQuery, userID).Scan(&userID)
	if err != nil {
		logrus.Println("Error checking user:", err)
		return models.UserNotFoundErr
	}
	return nil
}
