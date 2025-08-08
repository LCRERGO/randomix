package randomix

import "strings"

func (r *randomix) readBytes(buf []byte) {
	for i := range len(buf) {
		buf[i] = byte(r.Uint32())
	}
}

func randomString(r *randomix, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[r.IntN(len(charset))])
	}

	return sb.String()
}
