package morphinternal

import (
	"compress/gzip"
	"io"
	"log"
)

// MorphDeserializer предоставляет методы для распаковки данных
type MorphDeserializer struct{}

// DeflateGzip — аналог DeflateGzip из C# Pullenti.
// Распаковывает GZip-архив из `r` в `w`, обрабатывая случай ошибок чтения.
func (MorphDeserializer) DeflateGzip(r io.Reader, w io.Writer) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		log.Printf("gzip.NewReader error: %v", err)
		return
	}
	defer gzipReader.Close()

	buf := make([]byte, 100000)
	for {
		for i := range buf {
			buf[i] = 0
		}
		n, err := gzipReader.Read(buf)
		if err != nil && err != io.EOF {
			// Эмулируем поведение C# — ищем последний ненулевой байт
			lastNonZero := -1
			for i := len(buf) - 1; i >= 0; i-- {
				if buf[i] != 0 {
					lastNonZero = i
					break
				}
			}
			if lastNonZero >= 0 {
				_, _ = w.Write(buf[:lastNonZero+1])
			}
			break
		}
		if n <= 0 {
			break
		}
		_, _ = w.Write(buf[:n])
	}
}
