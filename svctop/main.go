package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", pong)
	r.GET("/call", handler)
	r.Run()
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func handler(c *gin.Context) {
	resp, err := http.Get("http://svcmid:8080/call")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var midResponse ServiceResponse
	if err := json.Unmarshal(body, &midResponse); err != nil {
		midResponse.Error = err.Error()
		c.JSON(http.StatusBadRequest, midResponse)
	} else {
		c.JSON(http.StatusBadRequest, midResponse)
	}
}
