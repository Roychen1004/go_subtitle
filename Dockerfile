# 使用官方的 Go 映像作為基礎映像
FROM golang:1.21

# 設定環境變數
ENV PORT=30046

# 創建工作目錄
WORKDIR /app

# 複製專案代碼到容器內的工作目錄
COPY . .

# 安裝依賴並編譯應用程式
RUN go mod download
# RUN go build -o subtitle-app cmd/main.go

# 暴露指定的端口
EXPOSE $PORT

# 運行主要應用程式
CMD ["go","run","cmd/main.go"]
