// The handler package takes all the http requests through the server
package table

import (
	laze "lazer/laze"
)

// The Handler struct
type Handler struct {
	app *laze.App
}

// The function will return the instance of Handler struct
func New(app *laze.App) Handler {
	return Handler{app: app}
}
