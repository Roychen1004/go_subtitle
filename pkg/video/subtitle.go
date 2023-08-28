package video

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"go_subtitle/pkg/srt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	// "os/exec"
	// "strconv"
)

// 定義 JSON 資料對應的結構體
type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
	// WholeWordTimestamps []struct {
	// 	Word        string  `json:"word"`
	// 	Start       float64 `json:"start"`
	// 	End         float64 `json:"end"`
	// 	Probability float64 `json:"probability"`
	// 	Timestamp   float64 `json:"timestamp"`
	// } `json:"whole_word_timestamps"`
}
type WhisperResponse struct {
	Text     string        `json:"text"`
	Language string        `json:"language"`
	Segments []srt.Segment `json:"segments"`
}

func SubtitleHandler(w http.ResponseWriter, r *http.Request) {

	videoPath := "/home/roy/go_subtitle/uploads/video.mp4"     // 影片檔案路徑
	subtitlePath := "/home/roy/go_subtitle/uploads/output.srt" // 字幕檔案路徑
	outputPath := "/home/roy/go_subtitle/uploads/output.mp4"   // 輸出影片檔案路徑

	// 使用ffmpeg添加字幕到影片
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", subtitlePath, "-c:v", "copy", "-c:s", "mov_text", outputPath)

	// 執行ffmpeg命令
	err_cmd := cmd.Run()
	if err_cmd != nil {
		log.Println("執行ffmpeg命令時出現錯誤：", err_cmd)
		// fmt.Println("執行ffmpeg命令時出現錯誤：", err_cmd)
		return
	}
	// 此種方式為軟字幕
	fmt.Println("字幕已成功添加到影片中，輸出影片：", outputPath)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse uploaded file
	file, _, err := r.FormFile("video")
	if err != nil {
		log.Println("Failed to read file", err)
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Save uploaded file
	uploadDir := "/home/roy/go_subtitle/uploads"
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		log.Println("Failed to create upload directory", err)
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	videoPath := filepath.Join(uploadDir, "video.mp4")
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
	audioPath := filepath.Join(uploadDir, "audio.mp3")
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vn", "-acodec", "mp3", audioPath)
	err = cmd.Run()
	if err != nil {
		log.Println("Error extracting audio", err)
		http.Error(w, "Error extracting audio", http.StatusInternalServerError)
		return
	}

}
