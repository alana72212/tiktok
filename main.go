package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tiktok/src/models"
	"tiktok/src/signers/bogus"
	"tiktok/src/signers/edata"
	"tiktok/src/signers/gnarly"
	"tiktok/src/signers/mssdk"
	"tiktok/src/signers/strdata"
)

func writeErr(w http.ResponseWriter, msg string) {
	json.NewEncoder(w).Encode(models.Response{Error: msg})
}

func writeOK(w http.ResponseWriter, result string) {
	json.NewEncoder(w).Encode(models.Response{Result: result})
}

func handler[T any](sign func(*T) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req T
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, "invalid JSON")
			return
		}
		writeOK(w, sign(&req))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bogus", handler(
		func(r *models.BogusRequest) string {
			return bogus.BogusEnc(r.Params, r.Body, r.UserAgent, r.CFP)
		},
	))
	mux.HandleFunc("/mssdk", handler(
		func(r *models.DataRequest) string {
			return mssdk.MssdkEnc(r.Data)
		},
	))
	mux.HandleFunc("/gnarly", handler(
		func(r *models.GnarlyRequest) string {
			return gnarly.GnarlyEnc(r.Params, r.Body, r.UserAgent, r.Version, r.CFP)
		},
	))
	mux.HandleFunc("/strdata", handler(
		func(r *models.DataRequest) string {
			return strdata.StrDataEnc(r.Data)
		},
	))
	mux.HandleFunc("/edata", handler(
		func(r *models.DataRequest) string {
			return edata.EdataEnc(r.Data)
		},
	))

	host := "0.0.0.0:8080" // 0.0.0.0 defaults to your IP address (open) | 127.0.0.1 (local)
	srv := &http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("API : %s", host)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to boot : %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown : %v", err)
	}
}
