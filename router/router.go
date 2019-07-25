package router

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wuxiaobai24/gin-blog/models"
)

var r *gin.Engine

func init() {
	r = gin.Default()

	r.SetFuncMap(template.FuncMap{
		"ToDate": ToDate,
	})
	r.LoadHTMLGlob("./template/*.tmpl")

	r.StaticFS("/static", http.Dir("./static"))

	// test ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("index", index)
	r.GET("/posts", getPosts)
	r.GET("/post/:id", getPost)
}

func index(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{
		"title": "wuxiaobai24's blog",
	})
}

// ToDate template func
func ToDate(t time.Time) string {
	return fmt.Sprintf("%v-%v-%v", t.Year(), t.Month(), t.Day())
}

// getPosts Get Posts
func getPosts(c *gin.Context) {
	posts, err := models.GetPosts(0, 10)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "failue",
		})
		return
	}
	c.HTML(200, "posts.tmpl", gin.H{
		"Title": "Archive",
		"Posts": posts,
	})
}

// getPost get post by id
func getPost(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPost(id)
	if err != nil {
		panic(err)
	}
	c.HTML(200, "post.tmpl", gin.H{
		"Title": post.Title,
		"Post":  post,
		"Prev":  nil,
		"Next":  nil,
	})
}

// Run server
func Run(addr string) error {
	return r.Run(addr)
}
