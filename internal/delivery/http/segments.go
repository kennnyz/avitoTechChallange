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
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: models.MethodNotProvideErr.Error()})
		return
	}
	logrus.Println("creating segment")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var segmentData models.Segment
	err = json.Unmarshal(body, &segmentData)
	if err != nil {
		logrus.Println("failed to decode JSON data")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: "invalid json"})
		return
	}

	err = h.validateSegment(&segmentData)
	if err != nil {
		logrus.Println("failed to validate segment ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: err.Error()})
		return
	}

	err = h.userSegmentService.CreateSegment(r.Context(), segmentData.Name)
	if err != nil {
		m := models.ResponseMessage{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(m)
		return
	}
	logrus.Println("segment created: ", segmentData.Name)

	_ = json.NewEncoder(w).Encode(models.SucceedMessage)
}

func (h *Handler) deleteSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: models.MethodNotProvideErr.Error()})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var segmentData models.Segment
	err = json.Unmarshal(body, &segmentData)
	if err != nil {
		logrus.Println("failed to decode JSON data")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: "invalid json"})
		return
	}

	err = h.validateSegment(&segmentData)
	if err != nil {
		logrus.Println("failed to validate segment ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: err.Error()})
		return
	}

	logrus.Println("deleting segment" + segmentData.Name)

	err = h.userSegmentService.DeleteSegment(r.Context(), segmentData.Name)
	if err != nil {
		m := models.ResponseMessage{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(m)
		return
	}

	_ = json.NewEncoder(w).Encode(models.SucceedMessage)
}

func (h *Handler) validateSegment(in *models.Segment) error {
	if in.Name == "" {
		return models.SegmentNameEmptyErr
	}
	return nil
}
