package main

const (
	black  = "0;30"
	red    = "0;31"
	green  = "0;32"
	yellow = "0;33"
	blue   = "0;34"
	purple = "0;35"
	cyan   = "0;36"
	white  = "0;37"
)

func ANSIColor(color string) string {
	return "\033[" + color + "m"
}

func ANSINormal() string {
	return "\033[0m"
}

func ColorFromSize(size int64) string {

	var pow = func(b, e int64) int64 {
		var out int64 = 1
		for e != 0 {
			out *= b
			e -= 1
		}
		return out
	}

	var color string
	if size < 1024 {
		color = black // B
	} else if size < pow(1024, 2) {
		color = white // KB
	} else if size < pow(1024, 3) {
		color = yellow // MB
	} else if size < pow(1024, 4) {
		color = red // GB
	} else {
		color = purple
	}
	return ANSIColor(color)
}
