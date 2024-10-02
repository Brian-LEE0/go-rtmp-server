package handler

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// FLV 파일을 청크 단위로 스트리밍하는 HTTP 핸들러
func StreamFLVHandler(w http.ResponseWriter, r *http.Request) {

	// FLV 스트림을 보내기 위한 헤더 설정
	w.Header().Set("Content-Type", "video/x-flv")
	w.Header().Set("Transfer-Encoding", "chunked")
	flvFile, err := os.Open(flvFilePath) // FLV 파일 경로
	if err != nil {
		http.Error(w, "Failed to open FLV file", http.StatusInternalServerError)
		return
	}

	log.Printf("Send FLV file: %v", flvFilePath)
	// flvFile.Seek(-chunkSize*60, 2) // 파일 포인터를 처음으로 이동

	buffer := make([]byte, chunkSize)
	for {
		n, err := flvFile.Read(buffer)
		if n > 0 {
			log.Printf("Send chunk: %v", n)
			if _, err := w.Write(buffer[:n]); err != nil {
				log.Printf("Failed to write chunk: %v", err)
				break
			}
			w.(http.Flusher).Flush()          // 클라이언트로 청크를 즉시 전송
			time.Sleep(10 * time.Millisecond) // 속도 조절 (옵션)
		}
		if err == io.EOF {
			log.Printf("Finish to send FLV file: %v", flvFilePath)
			break
		}
		if err != nil {
			log.Printf("Failed to read FLV file: %v", err)
			break
		}
	}
	flvFile.Close()
}

func ServeHTMLHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/index.html") // HTML 파일 경로
}
