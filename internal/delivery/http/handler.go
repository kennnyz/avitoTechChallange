package http_delivery

import (
	"context"
	"net/http"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

//1. Create Segment
//2. Update segment -
//3. Delete Segment
//4. Add user to segment
//5. Delete user from segment
//6. Get users by id from segment
//7. get all users

type UserSegmentService interface {
	CreateSegment(ctx context.Context, segmentName string) error
	DeleteSegment(ctx context.Context, segmentName string) error
	AddUserToSegments(ctx context.Context, segments models.AddUserToSegment) (models.AddUserToSegmentResponse, error)
	GetActiveUserSegments(ctx context.Context, userID int) ([]string, error)
}

type Handler struct {
	// access to business logic of our services
	userSegmentService UserSegmentService
}

func NewHandler(userSegmentService UserSegmentService) *Handler {
	return &Handler{
		userSegmentService: userSegmentService,
	}
}

func (h *Handler) Init() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/create-segment", h.createSegment)
	mux.HandleFunc("/delete-segment", h.deleteSegment)
	mux.HandleFunc("/add-user-to-segment", h.addUserToSegment)
	mux.HandleFunc("/active-user-segments", h.getActiveUserSegments)

	return mux
}
