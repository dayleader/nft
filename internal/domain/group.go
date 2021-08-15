package domain

// GroupID - group id.
type GroupID string

// GroupWrite struct.
type GroupWrite struct {
	Name     string `json:"name"`
	Priotiry int    `json:"priotiry"`
	XPos     int    `json:"xpos"`
	YPos     int    `json:"ypos"`
}

// GroupRead struct.
type GroupRead struct {
	ID GroupID `json:"id"`
	GroupWrite
}

// GroupRepository - provides access to the storage.
type GroupRepository interface {
	Create(group *GroupWrite) (GroupID, error)
	GetByID(groupID GroupID) (*GroupRead, error)
	GetByName(name string) (*GroupRead, error)
	GetAll() ([]*GroupRead, error)
}
