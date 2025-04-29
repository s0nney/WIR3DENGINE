package main

import (
	"WIR3DENGINE/controller"
	"html/template"

	"WIR3DENGINE/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	router := gin.Default()

	store := cookie.NewStore([]byte("changeme"))
	router.Use(sessions.Sessions("tempsession", store))
	router.Use(csrf.Middleware(csrf.Options{
		Secret: "changeme",
		ErrorFunc: func(c *gin.Context) {
			c.HTML(400, "err403.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"error":      "Forbidden.",
			})
			c.Abort()
		},
	}))

	router.MaxMultipartMemory = 2 << 20

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static/")
	router.Static("/tmp", "./tmp/")
	router.Static("/images", "./public/images/")
	router.StaticFile("/favicon.ico", "./favicon.ico")

	controller.SpinUpRoutes(router)
	router.Run()
}
