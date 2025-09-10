package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Homework struct {
	Text      string `json:"text"`
	Files     []File `json:"files"`
	UpdatedAt string `json:"updated_at"`
}

type File struct {
	Name string `json:"name,omitempty"` // omitempty - скрываем имя файла
	Type string `json:"type"`
	URL  string `json:"url"`
}

type HomeworkRequest struct {
	Text      string `json:"text"`
	Files     []File `json:"files"`
	UpdatedAt string `json:"updated_at"`
}

const (
	dataFile = "homework.json"
	port     = ":8080"
)

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Функция для преобразования формата даты
func formatDateTime(isoTime string) string {
	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		return isoTime // возвращаем оригинальную строку при ошибке
	}

	// Форматируем в нужный формат: 2025-09-09 14:50
	return t.Format("2006-01-02 15:04")
}

func loadHomework(filename string) (*Homework, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &Homework{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var homework Homework
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&homework)
	if err != nil {
		return nil, err
	}

	// Преобразуем формат даты при загрузке
	homework.UpdatedAt = formatDateTime(homework.UpdatedAt)

	return &homework, nil
}

func saveHomework(filename string, homework *Homework) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(homework)
}

func getHomeworkHandler(w http.ResponseWriter, r *http.Request, filename string) {
	enableCORS(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	homework, err := loadHomework(filename)
	if err != nil {
		http.Error(w, "Error loading homework", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(homework)
}

func postHomeworkHandler(w http.ResponseWriter, r *http.Request, filename string) {
	enableCORS(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req HomeworkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Преобразуем дату перед сохранением
	req.UpdatedAt = formatDateTime(req.UpdatedAt)

	homework := Homework{
		Text:      req.Text,
		Files:     req.Files,
		UpdatedAt: req.UpdatedAt,
	}

	err = saveHomework(filename, &homework)
	if err != nil {
		http.Error(w, "Error saving homework", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Get the file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create uploads directory if it doesn't exist
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filepath := filepath.Join("uploads", filename)

	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Return file info (без имени файла)
	response := File{
		Type: header.Header.Get("Content-Type"),
		URL:  "/uploads/" + filename,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract filename from URL
	filename := strings.TrimPrefix(r.URL.Path, "/uploads/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	filepath := filepath.Join("uploads", filename)

	// Check if file exists
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the file
	http.ServeFile(w, r, filepath)
}

func main() {
	// Create uploads directory if it doesn't exist
	os.MkdirAll("uploads", os.ModePerm)

	// Set up routes for all subjects
	subjects := []string{
		"computer_graphics",
		"bjd",
		"com_practicum",
		"it",
		"engl113",
		"engl208",
		"math",
		"oap",
		"oss",
		"ofg",
		"op1c",
	}
	for _, subject := range subjects {
		http.HandleFunc("/api/homework/"+subject, func(subject string) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				filename := subject + ".json"
				switch r.Method {
				case "GET":
					getHomeworkHandler(w, r, filename)
				case "POST":
					postHomeworkHandler(w, r, filename)
				case "OPTIONS":
					enableCORS(&w)
					w.WriteHeader(http.StatusOK)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			}
		}(subject))
	}

	http.HandleFunc("/api/upload", uploadFileHandler)
	http.HandleFunc("/uploads/", serveFileHandler)

	// Serve static files (optional)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
