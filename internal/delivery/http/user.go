package http_delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

// adding user to segment
func (h *Handler) addUserToSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, models.MethodNotProvideErr.Error(), 405)
		return
	}
	logrus.Println("adding user to segment")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	var userData models.AddUserToSegment
	err = json.Unmarshal(body, &userData)
	if err != nil {
		logrus.Error("failed to decode JSON data")
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	res, err := h.userSegmentService.AddUserToSegments(r.Context(), userData)
	if err != nil {
		http.Error(w, "Failed to add user to segment: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) getActiveUserSegments(w http.ResponseWriter, r *http.Request) {
	// Getting all segments for one user
	if r.Method != http.MethodGet {
		http.Error(w, models.MethodNotProvideErr.Error(), 405)
		return
	}
	logrus.Println("getting user segments")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	var user models.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		logrus.Error("failed to decode JSON data")
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	segments, err := h.userSegmentService.GetActiveUserSegments(r.Context(), user.UserID)
	if err != nil {
		http.Error(w, "Failed to get user segments: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// write segments to response in body as json
	_ = json.NewEncoder(w).Encode(segments)

}
