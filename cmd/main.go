package main

import (
	"fmt"
	"go_subtitle/pkg/video"

	// "go_subtitle/pkg/whisper"
	"net/http"
)

func main() {
	port := 30037

	// 添加你的API路由，例如：
	http.HandleFunc("/api/subtitle", processHandler)

	fmt.Printf("Server is listening on port %d...\n", port)
	fmt.Printf("API url:\n/api/subtitle\n")
	fmt.Println("==========================================")

	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}

// processHandler 函数保持不变
func processHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("接收POST")
	// 允许特定的域名进行跨域请求
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

	// 允许特定的HTTP方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

	// 允许特定的请求标头
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 允许跨域请求包含凭据（如Cookie等）
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		// 处理预检请求（OPTIONS请求），返回允许的HTTP方法和标头
		w.WriteHeader(http.StatusOK)
		return
	}

	// 影片上傳及分離
	video.UploadHandler(w, r)

	// 使用whisper 產生 srt
	// whisper.Whisper(w, r)

	// 結合影片及 srt
	video.SubtitleHandler(w, r)
}
