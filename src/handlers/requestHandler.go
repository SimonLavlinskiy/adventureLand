package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"project0/src/controllers/itemController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
)

func RequestHandler() {
	mysqlHost, _ := os.LookupEnv("MYSQL_HOST")

	r := gin.Default()

	r.GET("/get_full_map/:id", func(c *gin.Context) {
		id := helpers.ToInt(c.Param("id"))
		resp := repositories.GetFullMap(id)

		c.JSON(200, gin.H{
			"message": resp,
		})
	})

	r.GET("/get_maps", func(c *gin.Context) {
		resp := repositories.GetMaps()

		c.JSON(200, gin.H{
			"message": resp,
		})
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := uint(helpers.ToInt(c.Param("id")))
		resp := repositories.GetUserId(id)

		c.JSON(200, resp)
	})

	r.GET("/users", func(c *gin.Context) {
		resp := repositories.GetUsers()

		c.Header("X-Total-Count", fmt.Sprintf("%d", len(resp)))
		c.JSON(200, resp)
	})

	r.GET("/items", func(c *gin.Context) {
		resp := itemController.GetItems()

		c.Header("X-Total-Count", fmt.Sprintf("%d", len(resp)))
		c.JSON(200, resp)
	})

	r.GET("/items/:id", func(c *gin.Context) {
		id := uint(helpers.ToInt(c.Param("id")))
		resp := itemController.GetItemId(id)

		c.JSON(200, resp)
	})

	r.GET("/instruments", func(c *gin.Context) {
		itemId := helpers.ToInt(c.Request.URL.Query().Get("item_id"))

		resp := models.GetInstrumentsByItemId(itemId)

		c.Header("X-Total-Count", fmt.Sprintf("%d", len(resp)))
		c.JSON(200, resp)
	})

	err := r.Run(fmt.Sprintf("%s:33061", mysqlHost))

	if err != nil {
		panic(err)
	}

}
