# 使用官方的 Go 映像作為基礎映像
FROM golang:1.21

# 設定環境變數
ENV PORT=30046

# 暴露指定的端口
EXPOSE $PORT

# 設定工作目錄為應用程式根目錄
WORKDIR /app

# 複製go.mod和go.sum以便進行依賴管理
COPY go.mod go.sum ./

# 下載並安裝依賴
RUN go mod download

# 複製整個應用程式到容器內
COPY . .

# 編譯應用程式
RUN go build -o main ./cmd

# 指定執行時的命令
CMD ["./main"]