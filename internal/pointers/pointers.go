package pointers

func ToInt64(ip *int64) (i int64) {
	if ip != nil {
		i = *ip
	}
	return i
}

func ToString(sp *string) (s string) {
	if sp != nil {
		s = *sp
	}
	return s
}

func FromString(s string) (sp *string) {
	if s != "" {
		sp = &s
	}
	return sp
}