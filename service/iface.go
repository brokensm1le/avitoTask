package service

import (
	"github.com/lib/pq"
	"time"
)

type Story struct {
	PersonId  int64          `json:"personID"`
	Operation string         `json:"operation"`
	Segments  pq.StringArray `gorm:"type:text[]" json:"segments"`
	Time      time.Time      `time_format:"2006-01-02 15:04:05" json:"time"`
}

type Manager interface {
	PostPaS(PersonID int, Segments []string) error
	PostS(segments []string) error
	GetSegments(personID int) ([]string, error)
	GetIDs(segment string) ([]int64, error)
	DeleteSegment(segment string) error
	DeleteSegments(personID int, segments []string) error
	PostWithPer(segment string, IDs []int64, per int) ([]int64, error)
	GetHistory(personID int, timeFrom time.Time, timeTo time.Time) ([]Story, error)
}
