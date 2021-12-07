package server

import (
	"JumboDB/jumboDB-core/src/config"
	"JumboDB/jumboDB-core/src/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type HttpServer struct {
	service *service.Service
	config *config.TomlConfig
}

func NewHttpServer() *HttpServer {
	server := new(HttpServer)
	server.config = config.GetConfig()
	server.service = service.NewService(server.config.Storage.Engine)
	return server
}

func (i *HttpServer) StartListening() {
	port := i.config.Connection.Port
	log.Printf("Start listening port %d ...\n", port)
	path := ":" + strconv.Itoa(port)

	server := gin.Default()
	server.Use(cors.Default())
	server.GET("/health", i.healthCheck)

	server.GET("/resources", i.getAllResources)
	server.POST("/resources", i.createResource)

	server.GET("/resources/:key", i.getResource)
	server.DELETE("/resources/:key", i.delResource)

	server.Run(path)
}

func (i *HttpServer) getAllResources(c *gin.Context) {
	var allData = i.service.GetAllElements()
	c.JSON(200, gin.H{
		"data": allData,
	})

}

func (i *HttpServer) createResource(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	i.service.PutOneElement(json["key"].(string), json["value"].(string))
	c.JSON(200, gin.H{
		"status": "success",
	})
}

func (i *HttpServer) getResource(c *gin.Context) {
	key := c.Param("key")
	value := i.service.GetOneElement(key)
	c.JSON(200, gin.H{
		"key":   key,
		"value": value,
	})
}

func (i *HttpServer) delResource(c *gin.Context) {
	key := c.Param("key")
	i.service.DelOneElement(key)
	c.JSON(200, gin.H{
		"status": "success",
	})
}

func (i *HttpServer) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "OK",
	})
}