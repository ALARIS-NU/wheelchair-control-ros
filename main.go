package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/jacobsa/go-serial/serial"
)

type Device struct {
	isConnected     bool
	encoder_con     bool
	log_speed       bool
	mode            string
	logs            bool
	rosMasterAdress string

	forward byte
	right   byte
}

type command byte

const (
	// since iota starts with 0, the first value
	// defined here will be the default
	CUndefined command = iota
	CTurnOn
	CTurnOff
	CHorn
	CSmaller
	CBigger
	CSetCh
)

var Arduino Device
var port chan commandPack
var wordPtr *string
var boudRate *int
var isROSneeded *bool

func main() {
	port = make(chan commandPack)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	Arduino.isConnected = false
	Arduino.mode = "slider"

	isUARTlogsNeeded := flag.Bool("uart", false, "do you need uart logs?")
	isROSneeded = flag.Bool("ros", true, "do you need a ros node?")
	rosMasterAdress := flag.String("rosMaster", "127.0.0.1:11311", "ros master adress")
	isGUIneeded := flag.Bool("gui", true, "do you need GUI?")
	isSpeedLogNeeded := flag.Bool("speed_log", false, "do you need speed logs?")
	// wordPtr := flag.String("port", "/dev/tty.usbmodem21201", "serial device abs path")
	wordPtr = flag.String("port", "/dev/tty.usbserial-1120", "serial device abs path for motor")
	encoder_port := flag.String("port_encoder", "/dev/tty.usbmodem11101", "serial device abs path for encoders")
	boudRate := flag.Int("rate", 115200, "serial boudrate uint (9600,115200,?)")
	flag.Parse()
	Arduino.logs = *isUARTlogsNeeded
	Arduino.rosMasterAdress = *rosMasterAdress
	Arduino.log_speed = *isSpeedLogNeeded

	options_encs := serial.OpenOptions{
		PortName: *encoder_port,
		// PortName:        "/dev/tty.ACM0",
		BaudRate:        uint(*boudRate),
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port_enc, err := serial.Open(options_encs)
	if err != nil {
		color.Red("encoder is not connected")
	} else {
		Arduino.encoder_con = true
		defer port_enc.Close()
		go read_encoder(port_enc)
		if *isROSneeded {
			go init_encoder_ROS()
			go init_twistListener()
		}
	}

	if *isGUIneeded {
		go init_gin()
	}

	// for {
	sig := <-sigCh
	fmt.Println("Received signal:", sig)
	// }

	var command commandPack
	command.Action = byte(CSetCh)
	command.Ch_name = 0
	command.Value = 174
	EasyTransferSend(port, command)

	command.Action = byte(CSetCh)
	command.Ch_name = 1
	command.Value = 174
	EasyTransferSend(port, command)

	color.Green("Application terminated gracefully")
}
