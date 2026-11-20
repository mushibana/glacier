package domain

import "strings"

func StringOrNil(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}
func Int64OrNil(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

func BoolOrNil(b bool) *bool {
	if b == false {
		return nil
	}
	return &b
}
