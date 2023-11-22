package util

func CopySliceString(vals []string) []string {
	ret := make([]string, len(vals))
	for idx, val := range vals {
		ret[idx] = val
	}
	return ret
}
