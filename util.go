package main

func Clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func IndexOf(a []string, x string) int {
	for i, n := range a {
		if n == x {
			return i
		}
	}
	return -1
}

func Insert(a []string, x string) []string {
	if Contains(a, x) {
		return a
	}
	return append(a, x)
}

func Remove(a []string, x string) []string {
	i := IndexOf(a, x)
	if i == -1 {
		return a
	}
	return append(a[:i], a[i+1:]...)
}
