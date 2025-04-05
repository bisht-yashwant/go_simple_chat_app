package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var roomManager *Manager

func main() {
	// Initialize the Manager once
	roomManager = NewRoomManager()

	router := gin.Default()
	router.SetHTMLTemplate(html)

	// Join page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "join_page", nil)
	})
	// Handle join form
	router.POST("/join", func(c *gin.Context) {
		roomid := c.PostForm("roomid")
		username := c.PostForm("username")
		c.Redirect(http.StatusFound, fmt.Sprintf("/room/%s?user=%s", roomid, username))
	})
	// Chat room UI
	router.GET("/room/:roomid", roomGET)
	// Receive posted messages
	router.POST("/room/:roomid", roomPOST)
	// Server-Sent Events stream
	router.GET("/stream/:roomid", stream)

	fmt.Println("Server running at http://localhost:8080")
	router.Run(":8080")
}

func roomGET(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := c.Query("user")
	if userid == "" {
		userid = "Guest"
	}
	c.HTML(http.StatusOK, "chat_room", gin.H{
		"roomid": roomid,
		"userid": userid,
	})
}

func roomPOST(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := c.PostForm("user")
	message := c.PostForm("message")
	roomManager.Submit(userid, roomid, message)
	c.Status(http.StatusOK)
}

func stream(c *gin.Context) {
	roomid := c.Param("roomid")
	listener := roomManager.OpenListener(roomid)
	defer roomManager.CloseListener(roomid, listener)

	clientGone := c.Request.Context().Done()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case msg := <-listener:
			c.SSEvent("message", msg)
			return true
		}
	})
}
