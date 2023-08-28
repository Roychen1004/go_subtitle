package main

import (
	"fmt"
	"go_subtitle/pkg/video"
	"go_subtitle/pkg/whisper"

	// "go_subtitle/pkg/whisper"
	"net/http"
)

func main() {
	port := 30046

	http.HandleFunc("/api/subtitle", processHandler)
	// http.HandleFunc("/api/whisper", whisper.Whisper)
	// http.HandleFunc("/api/combine", video.SubtitleHandler)

	fmt.Printf("Server is listening on port %d...\n", port)
	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	// 影片上傳及分離
	video.UploadHandler(w, r)

	// 使用whisper 產生 srt
	whisper.Whisper(w, r)

	// 結合影片及 srt
	video.SubtitleHandler(w, r)

}
