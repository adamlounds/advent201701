package main

import (
	"errors"
	"fmt"
	"github.com/firstrow/tcp_server"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"strconv"
	"strings"
)

const (
	defaultTelnetPort = ":3333"
)

var serverCtx context.Context

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger
	serverCtx = context.Background()
	serverCtx = lg.WithLoggerContext(serverCtx, logger)
	lg.Log(serverCtx).Infof("Starting Innocuous server %s", "v1.0")

	linesChan := make(chan string)
	go generateResponses(linesChan)

	startTelnetServer(defaultTelnetPort, linesChan)

	defer close(linesChan)

	select {}
}

func generateResponses(linesChan chan string) {
	lg.Log(serverCtx).Infof("parser starting")
	for {
		line := <-linesChan
		lg.Log(serverCtx).Debugf("received line %s", line)
		line = strings.TrimSuffix(line, "\n")

		nums, err := strToArray(line)
		if err != nil {
			continue
		}

		var tot uint64 = 0
		var prevNum uint8 = 0
		for _, num := range nums {
			if num == prevNum {
				tot += uint64(num)
			}
			prevNum = num
		}

		if nums[0] == nums[len(nums)-1] {
			lg.Log(serverCtx).Debugf("wrap match! %d", nums[0])

			tot += uint64(nums[0])
		}

		lg.Log(serverCtx).Infof("total %d", tot)
	}
}

func strToArray(line string) ([]uint8, error) {
	digits := strings.Split(line, "")
	nums := make([]uint8, len(digits))

	for _, digit := range digits {
		num, err := strconv.ParseUint(digit, 10, 8)
		if err != nil {
			lg.Log(serverCtx).Debugf("bad non-digit %s", digit)
			return nil, errors.New("not a number")
		}
		nums = append(nums, uint8(num))
	}
	return nums, nil
}

func startTelnetServer(TelnetPort string, linesChan chan string) {
	server := tcp_server.New(TelnetPort)

	server.OnNewClient(func(c *tcp_server.Client) {
		c.Send(fmt.Sprintf("Hello\n"))
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		c.Send(fmt.Sprintf("Thank you! [%s]\n", message))
		linesChan <- message
		c.Send(fmt.Sprintf("parsing [%s]\n", message))

	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
	})

	go server.Listen()
	lg.Log(serverCtx).Infof("Started telnet server on port %s", TelnetPort)
}
