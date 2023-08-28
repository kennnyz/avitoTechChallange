package http_delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

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
	if userData.UserID < 0 {
		logrus.Println("bad user id")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: "bad user id"})
		return
	}

	res, err := h.userSegmentService.AddUserToSegments(r.Context(), userData)
	if err != nil {
		m := models.ResponseMessage{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(m)
		return
	}
	_ = json.NewEncoder(w).Encode(res)
}

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
