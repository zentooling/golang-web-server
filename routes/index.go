package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index renders the HTML of the index page
func Index(c *gin.Context) {
	pd := DefaultPageData(c)
	pd.Title = pd.Trans("Home")
	c.HTML(http.StatusOK, "index.html", pd)
}
