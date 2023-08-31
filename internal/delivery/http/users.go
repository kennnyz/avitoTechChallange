package http_delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

// @Summary Add To User Segments And Delete From User Segments
// @Description Add To User Segments And Delete From User Segments
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.AddUserToSegment true "user info"
// @Success 200 {object} models.ResponseMessage
// @Failure 400 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Failure default {object} models.ResponseMessage
// @Router /add-user-to-segment [post]
func (h *Handler) addUserToSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: models.MethodNotProvideErr.Error()})
		return
	}
	logrus.Println("adding user to segment")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var userData models.AddUserToSegment
	err = json.Unmarshal(body, &userData)
	if err != nil {
		logrus.Println("failed to decode JSON data")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: "invalid json"})
		return
	}
	if userData.UserID <= 0 {
		logrus.Println("bad user id")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: "bad user id"})
		return
	}

	err = h.userSegmentService.AddUserToSegments(r.Context(), userData)
	if err != nil {
		m := models.ResponseMessage{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(m)
		return
	}
	_ = json.NewEncoder(w).Encode(models.SucceedMessage)
}

// @Summary Get Active User Segments
// @Description Get Active User Segments
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.User true "user info"
// @Success 200 {object} []string
// @Failure 400 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Failure default {object} models.ResponseMessage
// @Router /active-user-segments [get]
func (h *Handler) getActiveUserSegments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: models.MethodNotProvideErr.Error()})
		return
	}
	logrus.Println("getting user segments")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		logrus.Println("failed to decode JSON data")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: "invalid json"})
		return
	}
	if user.UserID < 0 {
		logrus.Println("bad user id")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: "bad user id"})
		return
	}

	segments, err := h.userSegmentService.GetActiveUserSegments(r.Context(), user.UserID)
	if err != nil {
		m := models.ResponseMessage{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(m)
		return
	}
	_ = json.NewEncoder(w).Encode(segments)
}
