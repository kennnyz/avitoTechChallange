package http_delivery

import (
	"context"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "www.github.com/kennnyz/avitochallenge/docs"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_hash.go

type UserSegmentService interface {
	CreateSegment(ctx context.Context, segmentName string) error
	DeleteSegment(ctx context.Context, segmentName string) error
	AddUserToSegments(ctx context.Context, segments models.AddUserToSegment) error
	GetActiveUserSegments(ctx context.Context, userID int) ([]string, error)
	GetHistoryFile(ctx context.Context, year, month string) (string, error)
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

func (h *Handler) Init(swaggerUrl string) http.Handler {
	r := chi.NewRouter()

	r.Post("/create-segment", h.createSegment)
	r.Delete("/delete-segment", h.deleteSegment)
	r.Post("/add-user-to-segment", h.addUserToSegment)
	r.Get("/active-user-segments", h.getActiveUserSegments)
	r.Get("/get-history", h.getHistoryFile)
	r.Get("/tmp/*", h.getFile)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerUrl)))

	return r
}
