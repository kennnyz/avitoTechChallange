package models

var (
	SucceedMessage = ResponseMessage{
		Message: "success",
	}
)

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

type AddUserToSegmentResponse struct {
	UserID           int      `json:"id"`
	AddedSegments    []string `json:"added_segments"`
	DeletedSegments  []string `json:"deleted_segments"`
	NotExistSegments []string `json:"not_exist_segments"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}
