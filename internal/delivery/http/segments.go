package http_delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

func (h *Handler) createSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, models.MethodNotProvideErr.Error(), 405)
		return
	}
	logrus.Println("creating segment")

	body, err := io.ReadAll(r.Body)

	var segmentData models.Segment
	err = json.Unmarshal(body, &segmentData)
	if err != nil {
		logrus.Println("failed to decode JSON data")
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	err = h.userSegmentService.CreateSegment(r.Context(), segmentData.Name)
	if err != nil {
		http.Error(w, "Failed to create segment "+err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Println("segment created: ", segmentData.Name)

	_ = json.NewEncoder(w).Encode(segmentData)
}

func (h *Handler) deleteSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, models.MethodNotProvideErr.Error(), 405)
		return
	}
	logrus.Println("deleting segment")

	body, err := io.ReadAll(r.Body)

	var segmentData models.Segment
	err = json.Unmarshal(body, &segmentData)
	if err != nil {
		logrus.Println("failed to decode JSON data")
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	err = h.userSegmentService.DeleteSegment(r.Context(), segmentData.Name)
	if err != nil {
		http.Error(w, "Failed to delete segment "+err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(segmentData)
}
