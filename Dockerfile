# 使用官方的 Go 映像作為基礎映像
FROM golang:1.21

# # 設定環境變數
# ENV PORT=30037

# # 暴露指定的端口
# EXPOSE $PORT

# 設定工作目錄為應用程式根目錄
WORKDIR /app

# 複製整個專案到容器中
COPY . .

RUN ls -a

# 下載專案的依賴
RUN go mod download

RUN ls -a

# 設定工作目錄
WORKDIR /app/cmd

RUN ls -a

# 編譯專案
RUN go build -o my-golang

COPY .env /app/cmd

# 指定容器執行時運行的命令
CMD ["./my-golang"]
