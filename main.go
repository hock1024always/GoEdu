package main

import "github.com/hock1024always/GoEdu/router"

// 入口文件
func main() {
	r := router.Router()
	r.Run(":9090")
}
