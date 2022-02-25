package datastore

import "time"

const (
	TableNameWriteLock = "write_locks"
)

// WriteLock WriteMarker lock
type WriteLock struct {
	AllocationID string    `gorm:"primaryKey, column:allocation_id"`
	SessionID    string    `gorm:"column:session_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

// TableName get table name of migrate
func (WriteLock) TableName() string {
	return TableNameWriteLock
}