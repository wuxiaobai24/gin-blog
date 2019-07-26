package router

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wuxiaobai24/gin-blog/models"
	"github.com/wuxiaobai24/gin-blog/pkg"
)

var r *gin.Engine

func init() {
	r = gin.Default()

	r.SetFuncMap(template.FuncMap{
		"ToDate":   pkg.ToDate,
		"Markdown": pkg.Markdown2HTML,
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

	r.GET("/tags", getTags)
	r.GET("/tag/:id", getTag)

	r.GET("/admin", admin)

	r.POST("/upload", uploadPost)
}

func index(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{
		"Title": "wuxiaobai24's blog",
	})
}

func admin(c *gin.Context) {
	c.HTML(200, "admin.tmpl", gin.H{
		"Title": "Admin",
	})
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

//getTags return tags
func getTags(c *gin.Context) {
	tags, err := models.GetTags()
	if err != nil {
		panic(err)
		c.JSON(400, gin.H{"message": "failure"})
	}
	c.HTML(200, "tags.tmpl", gin.H{
		"Title": "tags",
		"Tags":  tags,
	})
}

func getTag(c *gin.Context) {
	id := c.Param("id")
	tag, err := models.GetTag(id)
	if err != nil {
		panic(err)
	}
	c.HTML(200, "posts.tmpl", gin.H{
		"Title": tag.Name,
		"Posts": tag.Posts,
	})
}

func uploadPost(c *gin.Context) {
	title := c.PostForm("title")
	taglist := strings.Split(c.PostForm("tags"), ",")

	file, err := c.FormFile("file")
	if err != nil {
		// c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		c.Redirect(http.StatusMovedPermanently, "/index")
		return
	}

	f, err := file.Open()
	if err != nil {
		// c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		c.Redirect(http.StatusMovedPermanently, "/index")
		return
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	var tags []*models.Tag
	for _, tagname := range taglist {
		tag, err := models.GenerateTag(strings.TrimSpace(tagname))
		if err != nil {
			panic(err)
			continue
		}
		tags = append(tags, tag)
	}

	err = models.AddPost(map[string]interface{}{
		"Title":   title,
		"Content": string(content),
		"Tags":    tags,
	})

	c.Redirect(http.StatusMovedPermanently, "/index")
}

// Run server
func Run(addr string) error {
	return r.Run(addr)
}
