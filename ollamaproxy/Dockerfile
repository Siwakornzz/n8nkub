# ขั้นตอน Builder
FROM golang:1.23.1 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o ollamaproxy -ldflags "-w -s" -a -tags netgo -installsuffix netgo
RUN ls -lh /app/ollamaproxy  # Debug: ตรวจสอบไฟล์

# ขั้นตอนสุดท้าย
FROM alpine:latest
RUN apk add --no-cache curl
COPY --from=builder /app/ollamaproxy /ollamaproxy
RUN ls -lh /ollamaproxy  # Debug: ตรวจสอบไฟล์
RUN chmod +x /ollamaproxy
EXPOSE 5001
CMD ["/ollamaproxy"]