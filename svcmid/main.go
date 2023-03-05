package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type WeatherResponse struct {
	Lat     float64         `json:"lat,omitempty"`
	Lon     float64         `json:"lon,omitempty"`
	Current *CurrentWeather `json:"current,omitempty"`
}

type CurrentWeather struct {
	FeelsLike float64 `json:"feels_like,omitempty"`
	TempF     float64 `json:"temp,omitempty"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", pong)
	r.GET("/call", handler2)
	r.Run()
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func handler(c *gin.Context) {
	resp, err := http.Get("http://svcbot.default:8080/call")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var botResponse ServiceResponse
	if err := json.Unmarshal(body, &botResponse); err != nil {
		botResponse.Error = err.Error()
		c.JSON(http.StatusBadRequest, botResponse)
	} else {
		c.JSON(http.StatusBadRequest, botResponse)
	}
}

func handler2(c *gin.Context) {
	url := "https://scottsbaldwin.github.io/weatherapi/weather/austin"
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var weatherResponse WeatherResponse
	var svcResponse ServiceResponse
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		svcResponse.Error = err.Error()
		c.JSON(http.StatusBadRequest, svcResponse)
	} else {
		svcResponse.Message = fmt.Sprintf("Current temperature is %.2f°F, but it feels like %.2f°F.", weatherResponse.Current.TempF, weatherResponse.Current.FeelsLike)
		c.JSON(http.StatusBadRequest, svcResponse)
	}
}
