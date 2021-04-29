package bot

func breakupString(s string, size int) []string {
	if size >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, size)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == size {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}
