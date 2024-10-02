package main

import (
	"io"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/yutopp/go-rtmp"

	"github.com/Brian-LEE0/go-rtmp-server/handler"
)

func main() {
	go func() {
		tcpAddr, err := net.ResolveTCPAddr("tcp", ":1935")
		if err != nil {
			log.Panicf("Failed: %+v", err)
		}

		listener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Panicf("Failed: %+v", err)
		}

		log.Printf("Start: %v", tcpAddr)
		srv := rtmp.NewServer(&rtmp.ServerConfig{
			OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
				l := log.StandardLogger()
				//l.SetLevel(logrus.DebugLevel)

				h := &handler.Handler{}

				return conn, &rtmp.ConnConfig{
					Handler: h,

					ControlState: rtmp.StreamControlStateConfig{
						DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
					},

					Logger: l,
				}
			},
		})
		if err := srv.Serve(listener); err != nil {
			log.Panicf("Failed: %+v", err)
		}
	}()

	// HTTP 서버로 FLV 스트리밍
	http.HandleFunc("/stream", handler.StreamFLVHandler)
	http.HandleFunc("/", handler.ServeHTMLHandler)
	log.Println("HTTP Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
