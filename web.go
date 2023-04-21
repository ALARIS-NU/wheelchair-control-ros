package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func init_gin() {
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
		if command.Value > 238 {
			command.Value = 238
		}
		if command.Value < 56 {
			command.Value = 56
		}
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

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "status" {
		fmt.Fprintf(w, "<h1>The connection is %v</h1>", Arduino.isConnected)
	} else {
		fmt.Fprintf(w, `<h1>Error 404</h1><p><a href="/status">/status</a></p>`)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p.Body)
}

func loadPage(title string) (*Page, error) {
	color.Yellow("loadpage, got: %s", title)
	if title == "" {
		title = "index.html"
	}
	filename := title
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
