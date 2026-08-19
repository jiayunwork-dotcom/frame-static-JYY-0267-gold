package format

import (
	"strconv"
	"strings"
)

// Fixed formats a float with fixed precision, trimming trailing zeros.
func Fixed(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

// Signed formats a float with an explicit leading sign when non-negative.
func Signed(f float64, prec int) string {
	if f >= 0 {
		return "+" + Fixed(f, prec)
	}
	return Fixed(f, prec)
}

// Pad right-pads s with spaces to width w.
func Pad(s string, w int) string {
	if len(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-len(s))
}

// Percent formats f (a fraction in [0,1]) as a percentage with prec decimals.
func Percent(f float64, prec int) string {
	return Fixed(f*100, prec) + "%"
}

// Eng renders a float in engineering notation: a mantissa scaled to a
// power-of-1000 exponent (… µ, m, "", k, M, G …).
func Eng(f float64) string {
	if f == 0 {
		return "0"
	}
	neg := false
	if f < 0 {
		neg = true
		f = -f
	}
	exp := 0
	for f >= 1000 {
		f /= 1000
		exp += 3
	}
	for f < 1 {
		f *= 1000
		exp -= 3
	}
	units := map[int]string{
		-9: "n", -6: "µ", -3: "m",
		0:  "", 3: "k", 6: "M", 9: "G",
	}
	u, ok := units[exp]
	if !ok {
		u = "e" + strconv.Itoa(exp)
	}
	mant := strconv.FormatFloat(f, 'f', 3, 64)
	if neg {
		return "-" + mant + u
	}
	return mant + u
}
