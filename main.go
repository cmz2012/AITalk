// Code generated by hertz generator.

package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cmz2012/AITalk/dal"
	"github.com/hertz-contrib/cors"
)

func main() {
	// init
	dal.InitClient()
	dal.InitDB()

	h := server.Default()
	h.NoHijackConnPool = true
	h.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	register(h)
	h.Spin()
}
