package whisper

import (
	"bytes"
	"encoding/json"
	"go_subtitle/pkg/srt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

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

var VIDEO_UPLOAD_PATH string

func Whisper(w http.ResponseWriter, r *http.Request) {
	// 載入 .env 檔案中的環境變數
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	VIDEO_UPLOAD_PATH = os.Getenv("VIDEO_UPLOAD_PATH")

	WHISPER_API_KEY := os.Getenv("WHISPER_API_KEY")
	url := "https://transcribe.whisperapi.com"
	headers := map[string]string{
		"Authorization": "Bearer " + WHISPER_API_KEY,
	}

	audioPath := VIDEO_UPLOAD_PATH + "audio.mp3"
	audioFile, err := os.Open(audioPath)
	if err != nil {
		log.Println("Failed to open audio file:", err)
		http.Error(w, "Failed to open audio file:", http.StatusInternalServerError)
		return
	}
	defer audioFile.Close()

	audioData, err := io.ReadAll(audioFile)
	if err != nil {
		log.Println("Failed to read audio data:", err)
		http.Error(w, "Failed to read audio data:", http.StatusInternalServerError)
		return
	}

	client := resty.New()

	req := client.R().
		SetHeaders(headers).
		SetFileReader("file", "audio.wav", bytes.NewReader(audioData)).
		SetFormData(map[string]string{
			"model":       "whisper-1",
			"fileType":    "wav",
			"diarization": "false",
			"numSpeakers": "1",
			"language":    "en",
			"task":        "transcribe",
		})

	response, err := req.Post(url)
	if err != nil {
		log.Println("Failed to send request:", err)
		http.Error(w, "Failed to send request:", http.StatusInternalServerError)
		return
	}

	var whisperResponse WhisperResponse
	err = json.Unmarshal(response.Body(), &whisperResponse)
	if err != nil {
		log.Println("Failed to parse Whisper response:", err)
		http.Error(w, "Failed to parse Whisper response:", http.StatusInternalServerError)
		return
	}

	// 使用 log 打印 Whisper API 响应的部分内容
	log.Println("Whisper API Response Text:", whisperResponse.Text)
	log.Println("Whisper API Response Language:", whisperResponse.Language)

	// 处理 Whisper API 响应
	subtitlePath := "/home/roy/go_subtitle/uploads/output.srt" // 替换为实际的输出路径
	err = srt.GenerateSRTFile(srt.WhisperResponse{
		Text:     whisperResponse.Text,
		Language: whisperResponse.Language,
		Segments: whisperResponse.Segments,
	}, subtitlePath)

	if err != nil {
		log.Println("Failed to generate .srt file:", err)
		http.Error(w, "Failed to generate .srt file:", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(whisperResponse)
}
