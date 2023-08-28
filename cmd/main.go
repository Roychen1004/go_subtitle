package main

import (
	"fmt"
	"go_subtitle/pkg/video"
	"go_subtitle/pkg/whisper"
	"net/http"
)

func main() {
	port := 30046

	http.HandleFunc("/api/subtitle", processHandler)

	fmt.Printf("Server is listening on port %d...\n", port)
	fmt.Printf("API url:\n/api/subtitle\n")
	fmt.Println("==========================================")

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
