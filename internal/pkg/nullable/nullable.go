package nullable

func StrOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Int16OrZero(i *int16) int16 {
	if i == nil {
		return 0
	}
	return *i
}
