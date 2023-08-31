package models

var (
	SucceedMessage = ResponseMessage{
		Message: "success",
	}
)

type AddUserToSegment struct {
	UserID           int      `json:"userid"`
	SegmentsToAdd    []string `json:"segments_to_add"`
	SegmentsToDelete []string `json:"segments_to_delete"`
}

type Segment struct {
	Name string `json:"name"`
}

type User struct {
	UserID int `json:"userid"`
}

type AddUserToSegmentResponse struct {
	UserID          int      `json:"userid"`
	AddedSegments   []string `json:"added_segments"`
	DeletedSegments []string `json:"deleted_segments"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type UserInSegment struct {
	UserID  int    `json:"userid"`
	Segment string `json:"segment"`
}

type GetHistoryRequest struct {
	Year  string `json:"year"`
	Month string `json:"month"`
}

type UserHistory struct {
	UserID  int    `json:"userid"`
	Segment string `json:"segment"`
	Action  string `json:"action"`
	Date    string `json:"date"`
}
