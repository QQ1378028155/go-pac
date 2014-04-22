package logger

var Verbose bool

func Debug(s string) {
	if Verbose == true {
		println(s)
	}
}
