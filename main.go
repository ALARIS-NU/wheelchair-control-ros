package main

import (
	"flag"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/jacobsa/go-serial/serial"
)

type Device struct {
	isConnected     bool
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

func main() {
	Arduino.isConnected = false
	Arduino.mode = "slider"

	isUARTlogsNeeded := flag.Bool("uart", false, "do you need uart logs?")
	isROSneeded := flag.Bool("ros", true, "do you need a ros node?")
	rosMasterAdress := flag.String("rosMaster", "127.0.0.1:11311", "ros master adress")
	isGUIneeded := flag.Bool("gui", true, "do you need GUI?")
	// wordPtr := flag.String("port", "/dev/tty.usbmodem21201", "serial device abs path")
	wordPtr := flag.String("port", "/dev/tty.usbserial-1120", "serial device abs path")
	boudRate := flag.Int("rate", 115200, "serial boudrate uint (9600,115200,?)")
	flag.Parse()
	Arduino.logs = *isUARTlogsNeeded
	Arduino.rosMasterAdress = *rosMasterAdress

	// Set up options.
	options := serial.OpenOptions{
		PortName: *wordPtr,
		// PortName:        "/dev/tty.ACM0",
		BaudRate:        uint(*boudRate),
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {

	} else {
		Arduino.isConnected = true
		defer port.Close()
		go rx(port)
		if *isROSneeded {
			go initROS()
		}
	}

	if *isGUIneeded {
		color.Yellow("GUI is needed")
		r := gin.Default()
		r.LoadHTMLGlob("templates/*")

		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.tmpl", gin.H{
				"isConnected": Arduino.isConnected,
				"mode":        Arduino.mode,
			})
		})
		// if Arduino.isConnected {
		r.Static("/static", "./static")
		// }
		r.GET("/changeModeToSlider", func(c *gin.Context) {
			Arduino.mode = "slider"
			c.HTML(200, "index.tmpl", gin.H{
				"isConnected": Arduino.isConnected,
				"mode":        Arduino.mode,
			})
		})

		r.GET("/changeModeToJoystick", func(c *gin.Context) {
			Arduino.mode = "joystick"
			c.HTML(200, "index.tmpl", gin.H{
				"isConnected": Arduino.isConnected,
				"mode":        Arduino.mode,
			})
		})

		r.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"isConnected": Arduino.isConnected,
				"mode":        Arduino.mode,
			})
		})
		r.GET("/action/on", func(ctx *gin.Context) {
			EasyTransferSend(port, commandPack{Action: byte(CTurnOn)})
			if Arduino.isConnected {
				ctx.JSON(200, gin.H{
					"status": "ok",
					"action": "on",
				})
			} else {
				ctx.JSON(500, gin.H{"error": "not connected"})
			}
		})
		r.GET("/action/off", func(ctx *gin.Context) {
			EasyTransferSend(port, commandPack{Action: byte(CTurnOff)})
			if Arduino.isConnected {
				ctx.JSON(200, gin.H{
					"status": "ok",
					"action": "off",
				})
			} else {
				ctx.JSON(500, gin.H{"error": "not connected"})
			}
		})
		r.GET("/action/horn", func(ctx *gin.Context) {
			EasyTransferSend(port, commandPack{Action: byte(CHorn)})
			if Arduino.isConnected {
				ctx.JSON(200, gin.H{
					"status": "ok",
					"action": "horn",
				})
			} else {
				ctx.JSON(500, gin.H{"error": "not connected"})
			}
		})
		r.GET("/action/speedDown", func(ctx *gin.Context) {
			EasyTransferSend(port, commandPack{Action: byte(CSmaller)})
			if Arduino.isConnected {
				ctx.JSON(200, gin.H{
					"status": "ok",
					"action": "speed Down",
				})
			} else {
				ctx.JSON(500, gin.H{"error": "not connected"})
			}
		})
		r.GET("/action/speedUp", func(ctx *gin.Context) {
			EasyTransferSend(port, commandPack{Action: byte(CBigger)})
			if Arduino.isConnected {
				ctx.JSON(200, gin.H{
					"status": "ok",
					"action": "speed Up",
				})
			} else {
				ctx.JSON(500, gin.H{"error": "not connected"})
			}
		})
		r.GET("/action/:ch_name/:value", func(ctx *gin.Context) {
			var command commandPack
			if err := ctx.ShouldBindUri(&command); err != nil {
				ctx.JSON(400, gin.H{"error": "could not bind command", "msg": err})
				return
			}
			if command.Ch_name > 4 {
				ctx.JSON(400, gin.H{"msg": "channel name out of bound [1..4]"})
				return
			}
			ctx.JSON(200, gin.H{
				"status":  "ok",
				"command": CSetCh,
				"ch_name": command.Ch_name,
				"value":   command.Value,
			})
			switch command.Ch_name {
			case 0:
				Arduino.forward = command.Value
			case 1:
				Arduino.right = command.Value
			}
			command.Action = byte(CSetCh)
			EasyTransferSend(port, command)
		})
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		openbrowser("http://localhost:8080")
		r.Run()

	}

	for {
	}

}
