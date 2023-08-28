// whisper/whisper.go

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
)

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

func Whisper(w http.ResponseWriter, r *http.Request) {
	url := "https://transcribe.whisperapi.com"
	headers := map[string]string{
		"Authorization": "Bearer WDLSUBC9KD6T6RJDXTSKRWNQHPJL36SZ",
	}

	audioPath := "/home/roy/go_subtitle/uploads/audio.mp3"
	audioFile, err := os.Open(audioPath)
	if err != nil {
		log.Println("Failed to open audio file", err)
		http.Error(w, "Failed to open audio file", http.StatusInternalServerError)
		return
	}
	defer audioFile.Close()

	audioData, err := io.ReadAll(audioFile)
	if err != nil {
		log.Println("Failed to read audio data", err)
		http.Error(w, "Failed to read audio data", http.StatusInternalServerError)
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
		log.Println("Failed to send request", err)
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}

	var whisperResponse WhisperResponse
	err = json.Unmarshal(response.Body(), &whisperResponse)
	if err != nil {
		log.Println("Failed to parse Whisper response", err)
		http.Error(w, "Failed to parse Whisper response", http.StatusInternalServerError)
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
		log.Println("Failed to generate .srt file", err)
		http.Error(w, "Failed to generate .srt file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(whisperResponse)
}
