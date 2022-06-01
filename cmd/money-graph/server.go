package main

import (
	"github.com/sbreitf1/money-graph/internal/moneydb"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func runServer(listenAddress string, db *moneydb.Database) error {
	gin.SetMode(gin.ReleaseMode)
	e := gin.Default()

	e.GET("/api/dbinfo", func(c *gin.Context) {
		c.JSON(200, struct {
			Name string `json:"name"`
		}{db.Name})
	})

	e.GET("/api/groups", func(c *gin.Context) {
		c.JSON(200, db.Groups)
	})

	logrus.Infof("run server on %q", listenAddress)
	return e.Run(listenAddress)
}
