package main

import "flag"
import "github.com/QQ1378028155/go-pac/logger"
import "github.com/QQ1378028155/go-pac/worker"

var (
	file    = flag.String("f", "", "configuration file to input.")
	verbose = flag.Bool("v", true, "give out some output on the screen")
)

func main() {
	flag.Parse()
	logger.Verbose = *verbose
	err := worker.Run(*file)
	if err != nil {
		logger.Debug(err.Error())
	}

}
