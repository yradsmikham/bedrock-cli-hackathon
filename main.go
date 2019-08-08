package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/yradsmikham/bedrock-cli/cmd"
)

func main() {
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "02-01-2006 15:04:05"
	formatter.FullTimestamp = true

	log.SetFormatter(formatter)

	cmd.Execute()
}
