package morphinternal

import (
	"compress/gzip"
	"io"
)

// DeflateGzip декомпрессирует данные из входного потока str и записывает результат в поток res.
// Аналог C#-метода: public static void DeflateGzip(Stream str, Stream res)
func DeflateGzip(str io.Reader, res io.Writer) error {
	// создаём gzip.Reader поверх входного потока
	deflate, err := gzip.NewReader(str)
	if err != nil {
		return err
	}
	defer deflate.Close()

	buf := make([]byte, 100000)
	for {
		// очистка буфера как в C# (по умолчанию Go не гарантирует обнуление после чтения)
		for i := range buf {
			buf[i] = 0
		}

		n, err := deflate.Read(buf)
		if err != nil && err != io.EOF {
			// как в C# — ищем последний ненулевой байт
			last := -1
			for i := len(buf) - 1; i >= 0; i-- {
				if buf[i] != 0 {
					last = i
					break
				}
			}
			if last >= 0 {
				_, _ = res.Write(buf[:last+1])
			}
			break
		}
		if n < 1 {
			break
		}
		_, _ = res.Write(buf[:n])
	}
	return nil
}
