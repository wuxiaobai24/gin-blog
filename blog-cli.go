package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/wuxiaobai24/gin-blog/config"
	"github.com/wuxiaobai24/gin-blog/models"
	"gopkg.in/yaml.v2"
)

var source = flag.String("path", "./source", "markdown source path")

func main() {
	config.InitConfig()
	models.InitDB()
	defer models.CloseDB()

	flag.Parse()
	var dirName = *source
	fmt.Printf("Source path is %s\n", *source)
	filenames, err := ioutil.ReadDir(*source)
	if err != nil {
		panic(err)
	}

	errch := make(chan error)
	okch := make(chan bool)

	for _, filename := range filenames {
		go func(filename string) {
			data, err := ioutil.ReadFile(filename)

			post, err := parseData(data)
			if err != nil {
				errch <- errors.New(fmt.Sprintf("%s:%v", filename, err))
				return
			}
			// save to models
			var tags []*models.Tag
			for _, tagname := range post["tags"].([]string) {
				tag, err := models.GenerateTag(tagname)
				if err != nil {
					errch <- errors.New(fmt.Sprintf("%s:%v", filename, err))
					return
				}
				tags = append(tags, tag)
			}
			models.AddPost(map[string]interface{}{
				"Title":   post["title"],
				"Content": post["content"],
				"Tags":    tags,
			})
			okch <- true
		}(filepath.Join(dirName, filename.Name()))
	}

	var count int
	var e error
	for _ = range filenames {
		select {
		case e = <-errch:
			fmt.Println(e)
		case <-okch:
			count++
		}
	}
	fmt.Println("count", count)
}

type YamlData struct {
	Title      string    `yaml:"title"`
	Date       time.Time `yaml:"date"`
	Categories string    `yaml:"categories"`
	Tags       []string  `yaml:"tags"`
}

func parseData(data []byte) (map[string]interface{}, error) {
	strs := strings.Split(string(data), "---")
	fmt.Println(strs[0])
	if len(strs) < 3 || strs[0] != "" {
		return nil, errors.New("Not Yaml")
	}
	var c YamlData
	err := yaml.Unmarshal([]byte(strs[1]), &c)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"title": c.Title, "date": c.Date,
		"categories": c.Categories,
		"content":    strings.Join(strs[2:], "---"),
		"tags":       c.Tags,
	}, nil
}
