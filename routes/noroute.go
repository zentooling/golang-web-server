package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NoRoute handles rendering of the 404 page
func NoRoute(c *gin.Context) {
	pd := DefaultPageData(c)
	pd.Title = pd.Trans("404 Not Found")
	c.HTML(http.StatusOK, "404.html", pd)
}
