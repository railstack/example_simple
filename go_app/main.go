package main

import (
	"flag"

	c "./controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	// The app will run on port 4000 by default, you can custom it with the flag -port
	servePort := flag.String("port", "4000", "Http Server Port")
	flag.Parse()

	// Here we are instantiating the router
	r := gin.Default()
	// Switch to "release" mode in production
	// gin.SetMode(gin.ReleaseMode)
	r.LoadHTMLGlob("views/*")
	// Create a static assets router
	// r.Static("/assets", "./public/assets")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	// Then we bind some route to some handler(controller action)
	// for the articles
	r.GET("/", c.HomeHandler)
	r.GET("/articles", c.ArticlesIndex)
	r.POST("/articles", c.ArticlesCreate)
	r.GET("/articles/:id", c.ArticlesShow)
	r.DELETE("/articles/:id", c.ArticlesDestroy)
	r.PUT("/articles/:id", c.ArticlesUpdate)
	// for the comments
	r.GET("/articles/:id/comments", c.CommentsIndex)
	r.POST("/comments", c.CommentsCreate)
	r.GET("/comments/:id", c.CommentsShow)
	r.DELETE("/comments/:id", c.CommentsDestroy)
	r.PUT("/comments/:id", c.CommentsUpdate)
	// Let's start the server
	r.Run(":" + *servePort)
}
