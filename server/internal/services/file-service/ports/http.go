package ports

import (
	"encoding/json"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
	"net/http"
)

type HttpServer struct {
	App app.Application
}

func writeProtobufError(w http.ResponseWriter, err apperr.AppError, code int) {
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(code)
}

func (h HttpServer) UploadMultipartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit memory usage for parsing multipart form.
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Retrieve the file part.
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Metadata can be extracted from the header or other form fields.
	fileName := header.Filename
	contentType := header.Header.Get("Content-Type")

	// Stream the file to your backend.
	resp, err := h.App.Commands.UploadTempFile.Handle(
		r.Context(),
		command.UploadTempFileParams{
			Reader:      file,
			FileName:    fileName,
			ContentType: contentType,
		},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to upload file: %v", err), http.StatusInternalServerError)
		return
	}

	// Return a JSON response.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
	}
}
