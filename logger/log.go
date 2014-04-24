package logger

//Verbose is set true for debug
var Verbose bool

//Debug to print some thing on screen
func Debug(s string) {
	if Verbose == true {
		println(s)
	}
}
