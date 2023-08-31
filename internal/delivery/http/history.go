package http_delivery

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

func (h *Handler) getHistoryFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: models.MethodNotProvideErr.Error()})
		return
	}
	logrus.Println("getting history")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var historyDate models.GetHistoryRequest
	err = json.Unmarshal(body, &historyDate)
	if err != nil {
		logrus.Println("failed to decode JSON data")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: "invalid json, please provide date in format YYYY-MM"})
		return
	}

	_, err = h.userSegmentService.GetHistoryFile(r.Context(), historyDate.Year, historyDate.Month)
	if err != nil {
		m := models.ResponseMessage{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(m)
		return
	}

	fileURL := fmt.Sprintf("http://%s/tmp/%s-%s.csv", r.Host, historyDate.Year, historyDate.Month)
	_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: fileURL})

}

func (h *Handler) getFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(models.ResponseMessage{Message: models.MethodNotProvideErr.Error()})
		return
	}

	filename := path.Base(r.URL.Path)
	filePath := filepath.Join(os.Getenv("HISTORY_FILES_PATH"), filename)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// getting file in bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// setting headers
	w.Header().Set("Content-Type", "application/csv")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	w.Write(fileBytes)
}
