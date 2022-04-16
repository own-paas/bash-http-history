package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type logger struct {
	UUID      string `json:"uuid"`
	Hostname  string `json:"hostname"`
	Type      string  `json:"type"`
	Ip        string `json:"ip"`
	Ppid      int32  `json:"ppid"`
	Pid       int32  `json:"pid"`
	Sid       int32  `json:"sid"`
	Uid       int32  `json:"uid"`
	User      string `json:"user"`
	Tty       string `json:"tty"`
	Pwd       string `json:"pwd"`
	Cmd       string `json:"cmd"`
	Timestamp int64  `json:"timestamp"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.PUT("/logger/:type", func(c *gin.Context) {
		result := logger{Type: c.Param("type")}
		err := c.ShouldBind(&result)
		fmt.Println(err,result)

		c.String(200, "ok")
	})
	return r}
func main() {
	r := setupRouter()
	r.Run(":6666")}