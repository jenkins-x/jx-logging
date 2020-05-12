package main

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
	"io/ioutil"
	"os"
)

func main() {
	file, err := ioutil.TempFile("/tmp", "test3.*.log")
	if err != nil {
		panic(err)
	}


	os.Setenv(log.JxLogLevel, "debug")
	os.Setenv(log.JxLogFile, file.Name())

	log.Logger().Debug("Debug")
	log.Logger().Info("Info")
	log.Logger().Warn("Warn")
	log.Logger().Error("Error")
	log.Logger().Fatal("Fatal")

}