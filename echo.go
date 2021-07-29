package main

import (
	"net/http"
	"strconv"

	"github.com/byuoitav/mayday/log"

	"github.com/labstack/echo/v4"
)

func serveStatus(port int) {
	log.P.Info("Starting status endpoint")

	router := echo.New()
	router.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "status ok")
	})

	err := router.Start(":" + strconv.Itoa(port))
	if err != nil {
		log.P.Panic(err.Error())
	}
}
