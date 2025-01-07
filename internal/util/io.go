package util

type SizeWriter struct {
	n int64
}

func (s *SizeWriter) Write(p []byte) (int, error) {
	s.n += int64(len(p))
	return len(p), nil
}

func (s *SizeWriter) Size() int64 {
	return s.n
}
