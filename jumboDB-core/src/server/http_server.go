package server

import (
	"JumboDB/jumboDB-core/src/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"log"
	"strconv"
)

func StartListening(port int) {
	log.Printf("Start listening port %d ...\n", port)
	path := ":" + strconv.Itoa(port)

	server := gin.Default()
	server.Use(cors.Default())
	server.GET("/health", healthCheck)

	server.GET("/resources", getAllResources)
	server.POST("/resources", createResource)

	server.GET("/resources/:key", getResource)
	server.DELETE("/resources/:key", delResource)

	server.Run(path)
}

func getAllResources(c *gin.Context) {
	var allData = service.GetAllElements()
	c.JSON(200, gin.H{
		"data": allData,
	})

}

func createResource(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	service.PutOneElement(json["key"].(string), json["value"].(string))
	c.JSON(200, gin.H{
		"status": "success",
	})
}

func getResource(c *gin.Context) {
	key := c.Param("key")
	value := service.GetOneElement(key)
	c.JSON(200, gin.H{
		"key":   key,
		"value": value,
	})
}

func delResource(c *gin.Context) {
	key := c.Param("key")
	service.DelOneElement(key)
	c.JSON(200, gin.H{
		"status": "success",
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "OK",
	})
}