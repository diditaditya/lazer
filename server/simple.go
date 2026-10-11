// Simple delivery mechanism of the lazer through http
package server

import (
	"github.com/gin-gonic/gin"

	laze "lazer/laze"

	"lazer/server/handler"
	"lazer/server/table"
)

// start the server
func Start(app *laze.App) {
	router := gin.Default()

	dataHandle := handler.New(app)
	tableHandle := table.New(app)

	// table CRUD
	tableRouter := router.Group("/_table")
	{
		tableRouter.GET("/", tableHandle.GetAll)
		tableRouter.POST("/", tableHandle.Create)
	}

	// data CRUD
	dataRouter := router.Group("/_data")
	{
		dataRouter.GET("/:name", dataHandle.GetAll)
		dataRouter.GET("/:name/:pk", dataHandle.GetByPk)
	
		dataRouter.POST("/:name", dataHandle.Create)
	
		dataRouter.PUT("/:name/:pk", dataHandle.UpdateByPk)
		dataRouter.PUT("/:name", dataHandle.Update)
	
		dataRouter.DELETE("/:name/:pk", dataHandle.DeleteByPk)
		dataRouter.DELETE("/:name", dataHandle.Delete)
	}

	router.Run()
}
