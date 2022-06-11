package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Ping(c *gin.Context) {
	ws, err := UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		log.Println("webasocket test", string(message))

		err = ws.WriteMessage(mt, message)

		if err != nil {
			break
		}
	}
}
