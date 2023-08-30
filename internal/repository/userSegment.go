package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

type UserSegmentRepo struct {
	db *pgxpool.Pool
}

func NewUserSegmentRepository(db *pgxpool.Pool) *UserSegmentRepo {
	return &UserSegmentRepo{
		db: db,
	}
}

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

func (u *UserSegmentRepo) AddUserToSegment(ctx context.Context, userID int, segments []string) ([]string, error) {
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
				return nil, check
			}
			return nil, err
		}
		if rw.RowsAffected() == 0 {
			err = u.CheckSegment(ctx, segmentName)
			if err != nil {
				return nil, err
			}
			continue
		}

		res = append(res, segmentName)
	}

	batch = &pgx.Batch{}
	for _, segmentName := range res {
		batch.Queue("INSERT INTO history (user_id, segment_name, action_type) VALUES ($1, $2, $3)", userID, segmentName, "add")
	}
	u.db.SendBatch(ctx, batch)

	return res, nil
}

func (u *UserSegmentRepo) DeleteUserFromSegments(ctx context.Context, userId int, segments []string) ([]string, error) {
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
			return nil, err
		}
		if rw.RowsAffected() == 0 {
			continue
		}
		history = append(history, segmentName)
	}

	batch = &pgx.Batch{}
	for _, segmentName := range history {
		batch.Queue("INSERT INTO history (user_id, segment_name, action_type) VALUES ($1, $2, $3)", userId, segmentName, "delete")
	}
	results = u.db.SendBatch(ctx, batch)
	for _, segmentName := range history {
		_, err := results.Exec()
		if err != nil {
			logrus.Println("Error adding to history:", err, " segment: ", segmentName)
			return nil, err
		}
	}

	return segments, nil
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

func (u *UserSegmentRepo) CheckSegment(ctx context.Context, segmentName string) error {
	checkSegmentQuery := "SELECT segment_name FROM segments WHERE segment_name = $1"
	err := u.db.QueryRow(ctx, checkSegmentQuery, segmentName).Scan(&segmentName)
	if err != nil {
		logrus.Println("Error checking segment:", err)
		return models.SegmentNotFoundErr(segmentName)
	}
	return nil
}
