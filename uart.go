package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/jacobsa/go-serial/serial"
)

func rx(f io.ReadWriteCloser) {
	// buff := make([]byte, 100)
	for {
		for {
			buf := make([]byte, 32)
			n, err := f.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading from serial port: ", err)
					f.Close()
					Arduino.isConnected = false
					break
				}
			} else {
				buf = buf[:n]
				if Arduino.logs {
					color.Yellow("%s", string(buf))
					color.Cyan("%s", hex.Dump(buf))
				}
			}
		}
		for !Arduino.isConnected {
			time.Sleep(time.Second * 1)
			color.Yellow("retrying Arduino connection")
		}
	}
}

// func sendSingleCommand(port chan commandPack, command byte) {
// 	b := commandPack{command, 0, 0}
// 	EasyTransferSend(port, b)
// }

// func sendCommand(port chan commandPack, action byte, x_axis float32, y_axis float32) {
// 	x_byte := byte(x_axis + 100)
// 	y_byte := byte(y_axis + 100)
// 	b := commandPack{action, x_byte, y_byte}
// 	EasyTransferSend(port, b)
// 	Arduino.forward = y_byte
// 	Arduino.right = x_byte
// }

func stop_wheelchair(port chan commandPack) {
	send_joy(port, 174, 174)
}

func send_joy(port chan commandPack, f byte, r byte) {
	var command commandPack
	command.Action = byte(CSetCh)
	command.Ch_name = 0
	command.Value = f
	EasyTransferSend(port, command)

	command.Action = byte(CSetCh)
	command.Ch_name = 1
	command.Value = r
	EasyTransferSend(port, command)
}

type commandPack struct {
	Action  byte
	Ch_name byte `uri:"ch_name"  binding:"number"`
	Value   byte `uri:"value" binding:"number"`
}

func EasyTransferSend(ch chan commandPack, in commandPack) {
	ch <- in
}

func EasyTransferInit(ch chan commandPack) {
	// Set up options.
	options := serial.OpenOptions{
		PortName: *wordPtr,
		// PortName:        "/dev/tty.ACM0",
		BaudRate:        uint(Arduino.boudRate),
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	port, err := serial.Open(options)
	if err != nil {
		color.Red("EasyTransfer: error opening port: %v", err)
		for {
			in := <-ch
			color.Red("EasyTransfer: dropped %v", in)
		}
	} else {
		Arduino.isConnected = true
		color.Green("EasyTransfer: connected to %s", *wordPtr)
		defer port.Close()
		go rx(port)
		if *isROSneeded {
			go initROS()
		}

		for {
			color.Cyan("EasyTransfer: waiting for data to send\n")
			in := <-ch
			color.Yellow("EasyTransfer: sending %v\n", in)
			size := reflect.TypeOf(in).Size()
			CS := byte(size)
			toOut := []byte{0x06, 0x85}
			toOut = append(toOut, byte(size))

			toOut = append(toOut, in.Action)
			CS ^= in.Action
			toOut = append(toOut, in.Ch_name)
			CS ^= in.Ch_name
			toOut = append(toOut, in.Value)
			CS ^= in.Value

			toOut = append(toOut, CS)
			if Arduino.logs {
				color.Cyan("Writing %v, as %v bytes using EasyTransfer\n", in, toOut)
			}
			if Arduino.isConnected {
				_, err := port.Write(toOut)
				if err != nil {
					log.Fatalf("port.Write: %v", err)
				}
			} else {
				color.Red("EasyTranfer: no device connected")
			}

		}
	}
}
