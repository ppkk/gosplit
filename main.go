package main

import (
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Project struct {
	ProjectID int
}

var (
	projects = []*Project{}
	mux      sync.Mutex
)

func createProject(c *gin.Context) {
	mux.Lock()
	defer mux.Unlock()

	projectID := len(projects) + 1
	project := &Project{ProjectID: projectID}
	projects = append(projects, project)

	c.JSON(200, gin.H{
		"projectID": projectID,
	})
}

func readProject(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))

	if projectID <= 0 || projectID > len(projects) {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	project := projects[projectID-1]

	c.JSON(200, project)
}

func main() {
	r := gin.Default()

	isMaster := os.Getenv("IS_MASTER") == "true"

	if isMaster {
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		apiURL, _ := url.Parse("http://api:8080")
		proxy := httputil.NewSingleHostReverseProxy(apiURL)
		r.Any("/eqh/*action", gin.WrapH(proxy))
	} else {
		eqh := r.Group("/eqh")
		{
			eqh.POST("/createProject", createProject)
			eqh.GET("/:id/readProject", readProject)
		}
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}
