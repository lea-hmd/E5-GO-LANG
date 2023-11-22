package middleware

import (
	"bytes"
	"github.com/sirupsen/logrus" // Permet de faire du logging structuré
	"net/http"
	"os"
	"time"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("logs/dictionary.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var buffer bytes.Buffer
		recorder := &responseRecorder{ResponseWriter: w, Body: &buffer}

		next.ServeHTTP(recorder, r)

		end := time.Now()

		responseBody := buffer.String()

		queryParams := r.URL.Query()

		logFields := logrus.Fields{
			"method":      r.Method,
			"path":        r.RequestURI,
			"proto":       r.Proto,
			"duration":    end.Sub(start).String(),
			"userAgent":   r.UserAgent(),
			"queryParams": queryParams,
			"statusCode":  recorder.statusCode,
			"response":    responseBody,
		}

		if recorder.statusCode >= 400 {
			log.WithFields(logFields).Error("Request failed")
		} else {
			log.WithFields(logFields).Info("Request handled successfully")
		}
	})
}
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	Body       *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}