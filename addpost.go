package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/wuxiaobai24/gin-blog/config"
	"github.com/wuxiaobai24/gin-blog/models"
)

func main() {

	config.InitConfig()
	models.InitDB()
	for _, filename := range os.Args[1:] {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
			fmt.Println(filename)
			continue
		}

		content := string(b)
		base := filepath.Base(filename)           // ./dir/file.md -> file.md
		suffix := path.Ext(base)                  // file.md -> .md
		title := strings.TrimSuffix(base, suffix) // fild.md -> file
		fmt.Println(title)
		err = models.AddPost(map[string]interface{}{
			"Title":   title,
			"Content": content,
		})
		if err != nil {
			panic(err)
			fmt.Println(filename)
			continue
		}
	}
}
