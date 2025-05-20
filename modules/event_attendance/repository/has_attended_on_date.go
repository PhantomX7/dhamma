package repository

import (
	"context"
	"time"

	"github.com/PhantomX7/dhamma/entity"
)

// HasAttendedOnDate checks if a follower attended a specific event on a given date.
// It queries the database for an EventAttendance record matching the followerID, eventID,
// and where the attended_at timestamp falls within the specified date.
func (r *repository) HasAttendedOnDate(ctx context.Context, followerID uint64, eventID uint64, date time.Time) (bool, error) {
	var count int64

	location, err := time.LoadLocation("Asia/Jakarta") // GMT+7
	if err != nil {
		return false, err
	}

	// Define the start and end of the given date
	year, month, day := date.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, location)
	endOfDay := time.Date(year, month, day, 23, 59, 59, 999999999, location)

	// Query the database
	err = r.db.WithContext(ctx).
		Model(&entity.EventAttendance{}).
		Where("follower_id = ?", followerID).
		Where("event_id = ?", eventID).
		Where("attended_at >= ?", startOfDay).
		Where("attended_at <= ?", endOfDay).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
