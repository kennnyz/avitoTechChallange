package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

func (u *UserSegmentRepo) CreateSegment(ctx context.Context, segmentName string) error {
	insertQuery := "INSERT INTO segments (segment_name) VALUES ($1) ON CONFLICT DO NOTHING"
	_, err := u.db.Exec(ctx, insertQuery, segmentName)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserSegmentRepo) DeleteSegment(ctx context.Context, segmentName string) error {
	userInSegment := []models.UserInSegment{}

	getUsers := "SELECT user_id, segment_name FROM user_segments WHERE segment_name = $1"
	rows, err := u.db.Query(ctx, getUsers, segmentName)
	if err != nil {
		logrus.Println("Error executing get query:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.UserInSegment{}
		err := rows.Scan(&user.UserID, &user.Segment)
		if err != nil {
			logrus.Println("Error scanning row:", err)
			return err
		}
		userInSegment = append(userInSegment, user)
	}

	deleteQuery := "DELETE FROM segments WHERE segment_name = $1"
	_, err = u.db.Exec(ctx, deleteQuery, segmentName)
	if err != nil {
		logrus.Println("Error executing delete query:", err)
		return err
	}

	batch := &pgx.Batch{}
	for _, user := range userInSegment {
		batch.Queue("INSERT INTO history (user_id, segment_name, action_type) VALUES ($1, $2, $3)", user.UserID, user.Segment, "deleted")
	}
	results := u.db.SendBatch(ctx, batch)
	results.Close()

	return nil
}

func (u *UserSegmentRepo) CheckSegment(ctx context.Context, segmentName string) error {
	checkSegmentQuery := "SELECT segment_name FROM segments WHERE segment_name = $1"
	err := u.db.QueryRow(ctx, checkSegmentQuery, segmentName).Scan(&segmentName)
	if err != nil {
		logrus.Println("Error checking segment:", err)
		return models.SegmentNotFoundErr(segmentName)
	}
	return nil
}
