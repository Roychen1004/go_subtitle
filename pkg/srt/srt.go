package srt

import (
	"fmt"
	"os"
	"strconv"
)

type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}
type WhisperResponse struct {
	Text     string    `json:"text"`
	Language string    `json:"language"`
	Segments []Segment `json:"segments"`
}

// 生成 .srt 檔案格式並輸出到檔案
func GenerateSRTFile(data WhisperResponse, subtitlePath string) error {
	fileContent := ""

	for i, segment := range data.Segments {
		// 轉換時間戳記為 .srt 檔案格式的時間格式
		startTime := formatTime(segment.Start)
		endTime := formatTime(segment.End)

		// 組成 .srt 檔案格式的內容
		fileContent += strconv.Itoa(i+1) + "\n"
		fileContent += startTime + " --> " + endTime + "\n"
		fileContent += segment.Text + "\n\n"
	}

	// 寫入 .srt 檔案
	err := os.WriteFile(subtitlePath, []byte(fileContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// 轉換時間戳記為 .srt 檔案格式的時間格式
func formatTime(timestamp float64) string {
	seconds := int(timestamp)
	milliseconds := int((timestamp - float64(seconds)) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d,%03d", seconds/3600, (seconds/60)%60, seconds%60, milliseconds)
}
