package repository

import (
	"context"
	"log"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

func (u *UserSegmentRepo) GetUsersHistory(ctx context.Context, year, month string) ([]models.UserHistory, error) {
	res := []models.UserHistory{}

	getUsersQuery := "SELECT user_id, segment_name, action_type, action_timestamp::varchar FROM history WHERE date_part('year', action_timestamp) = $1 AND date_part('month', action_timestamp) = $2"
	rows, err := u.db.Query(ctx, getUsersQuery, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		history := models.UserHistory{}
		err = rows.Scan(&history.UserID, &history.Segment, &history.Action, &history.Date)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = append(res, history)
	}

	return res, nil
}
