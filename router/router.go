package router

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wuxiaobai24/gin-blog/models"
	"github.com/wuxiaobai24/gin-blog/pkg"
)

var r *gin.Engine
var pageSize = 10

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

	// r.GET("/admin", admin)

	// r.POST("/upload", uploadPost)
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

func Pagination(c *gin.Context) (int, int) {
	pageNumStr := c.Query("page")
	pageNum := 1

	if pageNumStr != "" {
		num, err := strconv.Atoi(pageNumStr)
		if err == nil {
			pageNum = num
		}
	}
	offset := (pageNum - 1) * pageSize
	return pageNum, offset
}

// getPosts Get Posts
func getPosts(c *gin.Context) {
	page, offset := Pagination(c)
	posts, err := models.GetPosts(offset, pageSize)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "failue",
		})
		return
	}
	var next, prev string

	if page != 1 {
		prev = fmt.Sprintf("/posts?page=%v", page-1)
	}
	count, err := models.PostCount()
	if err != nil {
		c.Abort()
		return
	}
	if offset+pageSize <= count {
		next = fmt.Sprintf("/posts?page=%v", page+1)
	}

	c.HTML(200, "posts.tmpl", gin.H{
		"Title": "Archive",
		"Posts": posts,
		"Prev":  prev,
		"Next":  next,
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
