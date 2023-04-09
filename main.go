package main

import "gin-gorm-oj/router"

func main() {
	r := router.Router()
	_ = r.Run()
}
