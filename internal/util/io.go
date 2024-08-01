package util

type SizeWriter struct {
	n uint64
}

func (s *SizeWriter) Write(p []byte) (int, error) {
	s.n += uint64(len(p))
	return len(p), nil
}

func (s *SizeWriter) Size() uint64 {
	return s.n
}
