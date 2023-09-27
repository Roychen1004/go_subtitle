package video

import (
	// "bytes"
	// "encoding/json"
	"encoding/json"
	"fmt"
	"go_subtitle/pkg/srt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
	// "os/exec"
	// "strconv"
)

// 定義 JSON 資料對應的結構體
type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}
type WhisperResponse struct {
	Text     string        `json:"text"`
	Language string        `json:"language"`
	Segments []srt.Segment `json:"segments"`
}

var VIDEO_OUTPUT_PATH string
var VIDEO_UPLOAD_PATH string

func SubtitleHandler(w http.ResponseWriter, r *http.Request) {

	err_env := godotenv.Load()
	if err_env != nil {
		log.Fatal("Error loading .env file")
	}

	videoPath := VIDEO_UPLOAD_PATH + "video.mp4"       // 影片檔案路徑
	subtitlePath := VIDEO_UPLOAD_PATH + "output.srt"   // 字幕檔案路徑
	VIDEO_OUTPUT_PATH = os.Getenv("VIDEO_OUTPUT_PATH") // 輸出影片檔案路徑

	err := os.MkdirAll(VIDEO_OUTPUT_PATH, 0755)
	if err != nil {
		log.Println("Failed to create output directory:", err)
		http.Error(w, "Failed to create output directory", http.StatusInternalServerError)
		return
	}
	outputPath := VIDEO_OUTPUT_PATH + "output.mp4"

	// 使用ffmpeg添加字幕到影片
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", subtitlePath, "-c:v", "copy", "-c:s", "mov_text", outputPath)

	// 執行ffmpeg命令
	err_cmd := cmd.Run()
	if err_cmd != nil {
		log.Println("執行ffmpeg命令時出現錯誤:", err_cmd)
		// fmt.Println("執行ffmpeg命令時出現錯誤：", err_cmd)
		return
	}
	// 此種方式為軟字幕
	fmt.Println("字幕已成功添加到影片中，輸出影片：", outputPath)
	// uploadDir := "/home/roy/go_subtitle/uploads/"
	// 刪除 uploadDir 底下所有檔案
	// err_remove := os.RemoveAll(VIDEO_UPLOAD_PATH)
	// if err != nil {
	// 	log.Println("Failed to delete uploaded files:", err_remove)
	// }

	// 創建包含下載連結的JSON回應
	jsonResponse := map[string]string{"download_url": outputPath}
	responseJSON, err := json.Marshal(jsonResponse)
	if err != nil {
		log.Println("無法生成JSON回應:", err)
		http.Error(w, "無法生成JSON回應", http.StatusInternalServerError)
		return
	}
	// log.Println("responseJSON:", responseJSON)

	// 設置HTTP標頭
	w.Header().Set("Content-Type", "application/json")

	// 寫入JSON回應到HTTP回應
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Println("無法寫入JSON回應:", err)
		http.Error(w, "無法寫入JSON回應", http.StatusInternalServerError)
		return
	}

}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	err_env := godotenv.Load()
	if err_env != nil {
		log.Println("Error loading .env file:", err_env)
		log.Fatal("Error loading .env file")
	}

	VIDEO_UPLOAD_PATH = os.Getenv("VIDEO_UPLOAD_PATH")

	if r.Method != http.MethodPost {
		log.Println("Method not allowed:")
		http.Error(w, "Method not allowed:", http.StatusMethodNotAllowed)
		return
	}

	// Parse uploaded file
	file, _, err := r.FormFile("video")
	if err != nil {
		log.Println("Failed to read file:", err)
		http.Error(w, "Failed to read file:", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Save uploaded file
	// uploadDir := "/home/roy/go_subtitle/uploads"
	err = os.MkdirAll(VIDEO_UPLOAD_PATH, 0755)
	if err != nil {
		log.Println("Failed to create upload directory:", err)
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	videoPath := filepath.Join(VIDEO_UPLOAD_PATH, "video.mp4")
	outputFile, err := os.Create(videoPath)
	if err != nil {
		log.Println("Failed to create file", err)
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, file)
	if err != nil {
		log.Println("Failed to save file", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Extract audio using FFmpeg
	audioPath := filepath.Join(VIDEO_UPLOAD_PATH, "audio.mp3")
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vn", "-acodec", "mp3", audioPath)
	err = cmd.Run()
	if err != nil {
		log.Println("Error extracting audio", err)
		http.Error(w, "Error extracting audio", http.StatusInternalServerError)
		return
	}

}
