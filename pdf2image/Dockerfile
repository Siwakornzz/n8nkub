# ขั้นตอน Builder
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o pdf2image -ldflags "-w -s" -a -tags netgo -installsuffix netgo
RUN ls -lh /app/pdf2image  # Debug: ตรวจสอบไฟล์

# ขั้นตอนสุดท้าย
FROM alpine:latest
RUN apk add --no-cache mupdf-tools curl
COPY --from=builder /app/pdf2image /pdf2image
RUN ls -lh /pdf2image  # Debug: ตรวจสอบไฟล์
RUN chmod +x /pdf2image
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
CMD ["/pdf2image"]