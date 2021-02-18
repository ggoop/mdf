package utils

import (
	"bytes"
	"regexp"
	"strings"
)

/*
  first = false: "aaa_bbb_ccc" -> "aaaBbbCcc"
  first = true: "aaa_bbb_ccc" -> "AaaBbbCcc"
*/
func FirstCaseToUpper(str string, first bool) string {
	temp := strings.Split(str, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		if y == 0 && !first {
			continue
		}
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				vv[i] -= 32
				upperStr += string(vv[i])
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	if first {
		return upperStr
	} else {
		return temp[0] + upperStr
	}
}
func StringSplitFormat(str string) string {
	regExp := `[,|;|，|；\|]`
	r, _ := regexp.Compile(regExp)
	items := r.Split(str, -1)
	newItems := make([]string, 0)
	for _, item := range items {
		if item != "" {
			newItems = append(newItems, item)
		}
	}
	return strings.Join(newItems, ";")
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(name string) string {
	const (
		lower = false
		upper = true
	)
	if name == "" {
		return ""
	}
	var (
		value                                    = name
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')
		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}
	buf.WriteByte(value[len(value)-1])
	return strings.ToLower(buf.String())
}

// camel string, xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

/*
 [9 9 8 4 2 9 1 7 - a 5 4 b - 3 3 1 6 - c d f 3 - 8 7 d 9 f b 5 7] -> "99842917-a54b-3316-cdf3-87d9fb57"
*/
func ArrayToString(arrays []string) string {
	return strings.Join(arrays, "")
}
func StringIsAlphanumeric(str string) bool {
	if ok, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, str); ok {
		return ok
	}
	return false
}
func StringIsNumeric(str string) bool {
	if ok, _ := regexp.MatchString(`^[0-9]+$`, str); ok {
		return ok
	}
	return false
}
func StringIsURL(str string) bool {
	if ok, _ := regexp.MatchString(`^((https?):\/\/)+[^\s]+`, str); ok {
		return ok
	}
	return false
}
func StringIsEmail(str string) bool {
	if ok, _ := regexp.MatchString(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`, str); ok {
		return ok
	}
	return false
}
func StringIsCode(str string) bool {
	if ok, _ := regexp.MatchString(`^[a-zA-Z0-9_\/\.-]{1,50}$`, str); ok {
		return ok
	}
	return false
}
func StringIsMobile(str string) bool {
	//^(\\+\\d{2}-)?(\\d{2,3}-)?([1][3,4,5,7,8][0-9])\d{8}$
	if ok, _ := regexp.MatchString(`^(\\+\\d{2}-)?(\\d{2,3}-)?([1][3,4,5,7,8][0-9])\d{8}$`, str); ok {
		return ok
	}
	return false
}
