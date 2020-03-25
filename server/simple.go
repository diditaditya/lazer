package server

import (
	"fmt"
	"encoding/json"
	"github.com/gin-gonic/gin"
	laze "lazer/laze"
)

type Handler struct {
	app *laze.App
}

func (handler *Handler) root(c *gin.Context) {
	data := handler.app.GetAllTables()

	resp := map[string]interface{}{
		"message": "welcome",
		"data": data,
	}

	c.JSON(200, resp)
}

func (handler *Handler) getByPk(c *gin.Context) {
	tableName := c.Param("name")
	pk := c.Param("pk")

	data, err := handler.app.FindByPk(tableName, pk)

	if err != nil {
		fmt.Println(err)
		c.JSON(404, map[string]interface{}{
			"message": "oops..",
			"error": fmt.Sprintf("%s", err),
		})
	} else {
		resp := map[string]interface{}{
			"message": "yo",
			"data": data,
		}
	
		c.JSON(200, resp)
	}
}

func (handler *Handler) getAll(c *gin.Context) {
	tableName := c.Param("name")

	fmt.Printf("%v\n", c.Query)

	data, err := handler.app.FindAll(tableName)

	if err != nil {
		fmt.Println(err)
		c.JSON(404, map[string]interface{}{
			"message": "oops..",
			"error": fmt.Sprintf("%s", err),
		})
	} else {
		resp := map[string]interface{}{
			"message": "yo",
			"data": data,
		}
	
		c.JSON(200, resp)
	}
}

func (handler *Handler) create(c *gin.Context) {
	tableName := c.Param("name")
	raw, err := c.GetRawData()

	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "your fault..",
			"error": err,
		})
	}

	mapped := make(map[string]interface{})

	err = json.Unmarshal(raw, &mapped)

	if err != nil {
		c.JSON(500, map[string]interface{}{
			"message": "error unmarshaling the body",
			"error": err,
		})
	}

	_, err = handler.app.Create(tableName, mapped)

	if err != nil {
		c.JSON(500, map[string]interface{}{
			"message": "error something idk..",
			"error": err,
		})
	}

	c.JSON(200, map[string]interface{}{
		"message": "created..",
	})
}

func Start(app *laze.App) {
	router := gin.Default()

	handler := Handler{app: app}

	router.GET("/", handler.root)
	router.GET("/:name", handler.getAll)
	router.GET("/:name/:pk", handler.getByPk)

	router.POST("/:name", handler.create)

	router.Run()
}