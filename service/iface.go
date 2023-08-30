package service

type Manager interface {
	PostPaS(PersonID int, Segments []string) error
	PostS(segments []string) error
	GetSegments(personID int) ([]string, error)
	GetIDs(segment string) ([]int64, error)
	DeleteSegment(segment string) error
	DeleteSegments(personID int, segments []string) error
}
