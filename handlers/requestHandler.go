package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"project0/repository"
)

func RequestHandler() {
	mysqlHost, _ := os.LookupEnv("MYSQL_HOST")

	r := gin.Default()

	r.GET("/get_full_map/:id", func(c *gin.Context) {
		id := repository.ToInt(c.Param("id"))
		resp := repository.GetFullMap(id)

		c.JSON(200, gin.H{
			"message": resp,
		})
	})

	r.GET("/get_maps", func(c *gin.Context) {
		resp := repository.GetMaps()

		c.JSON(200, gin.H{
			"message": resp,
		})
	})

	err := r.Run(fmt.Sprintf("%s:33061", mysqlHost))

	if err != nil {
		panic(err)
	}

}
