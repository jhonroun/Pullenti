package morphinternal

import (
	"unicode/utf8"
)

// Сделан специально для Питона – а то стандартным способом через MemoryStream
// жутко тормозит, придётся делать самим
type ByteArrayWrapper struct {
	array []byte
	len   int
}

func NewByteArrayWrapper(arr []byte) *ByteArrayWrapper {
	return &ByteArrayWrapper{
		array: arr,
		len:   len(arr),
	}
}

// Проверка конца файла
func (b *ByteArrayWrapper) IsEOF(pos int) bool {
	return pos >= b.len
}

// Десериализация байта
func (b *ByteArrayWrapper) DeserializeByte(pos *int) byte {
	if *pos >= b.len {
		return 0
	}
	val := b.array[*pos]
	*pos++
	return val
}

// Десериализация short (2 байта, little-endian)
func (b *ByteArrayWrapper) DeserializeShort(pos *int) int {
	if *pos+1 >= b.len {
		return 0
	}
	b0 := b.array[*pos]
	*pos++
	b1 := b.array[*pos]
	*pos++
	return int(b1)<<8 | int(b0)
}

// Десериализация int (4 байта, little-endian)
func (b *ByteArrayWrapper) DeserializeInt(pos *int) int {
	if *pos+3 >= b.len {
		return 0
	}
	b0 := b.array[*pos]
	*pos++
	b1 := b.array[*pos]
	*pos++
	b2 := b.array[*pos]
	*pos++
	b3 := b.array[*pos]
	*pos++
	return int(b3)<<24 | int(b2)<<16 | int(b1)<<8 | int(b0)
}

// Десериализация строки с длиной в 1 байт
func (b *ByteArrayWrapper) DeserializeString(pos *int) string {
	if *pos >= b.len {
		return ""
	}
	length := b.DeserializeByte(pos)
	if length == 0xFF {
		return ""
	}
	if length == 0 {
		return ""
	}
	if *pos+int(length) > b.len {
		return ""
	}
	res := string(b.array[*pos : *pos+int(length)])
	if !utf8.ValidString(res) {
		return ""
	}
	*pos += int(length)
	return res
}

// Десериализация строки с длиной в 2 байта (short)
func (b *ByteArrayWrapper) DeserializeStringEx(pos *int) string {
	if *pos >= b.len {
		return ""
	}
	length := b.DeserializeShort(pos)
	if length == 0xFFFF || length < 0 {
		return ""
	}
	if length == 0 {
		return ""
	}
	if *pos+length > b.len {
		return ""
	}
	res := string(b.array[*pos : *pos+length])
	if !utf8.ValidString(res) {
		return ""
	}
	*pos += length
	return res
}
