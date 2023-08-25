package models

type AddUserToSegment struct {
	UserID           int      `json:"id"`
	SegmentsToAdd    []string `json:"segments_to_add"`
	SegmentsToDelete []string `json:"segments_to_delete"`
}

type Segment struct {
	Name string `json:"name"`
}

type User struct {
	UserID int `json:"userid"`
}
