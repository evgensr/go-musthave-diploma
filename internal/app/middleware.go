package app

import (
	"compress/gzip"
	_ "github.com/lib/pq" // ...
	"io"
	"log"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// GzipHandleEncode gzip-сжатие ответа
func (s *APIserver) GzipHandleEncode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			s.logger.Info(r.Header.Get("Accept-Encoding"))
			s.logger.Info("Not support gzip")
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			log.Println(err)
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		s.logger.Info(r.Header.Get("Accept-Encoding"))
		s.logger.Info("Support gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// GzipHandleDecode
func (s *APIserver) GzipHandleDecode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Println(err)
			io.WriteString(w, err.Error())
			return
		}
		r.Body = reader
		next.ServeHTTP(w, r)
	})
}
