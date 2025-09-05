package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/config"
)

type ConfigPageData struct {
	PageData
	Config *config.Config
}

func ConfigRouteHandler(c *gin.Context) {
	pd := DefaultPageData(c)
	pd.Title = pd.Trans("Admin")

	cd := ConfigPageData{
		PageData: pd,
		Config:   config.LairInstance().GetConfig(),
	}

	c.HTML(http.StatusOK, "config.html", cd)
}

func ConfigRouteHandlerPost(c *gin.Context) {
	pd := DefaultPageData(c)
	pd.Title = pd.Trans("Configuration")

	prevCfg := config.LairInstance().GetConfig()

	newCookie := c.PostForm("cookie-secret")

	slog.Info("ConfigRouteHandlerPost", "prevCookie", prevCfg.CookieSecret)
	if newCookie != prevCfg.CookieSecret {
		slog.Info("ConfigRouteHandlerPost", "newCookie", newCookie)
		prevCfg.CookieSecret = newCookie
	}

	cd := ConfigPageData{
		PageData: pd,
		Config:   config.LairInstance().GetConfig(),
	}

	c.HTML(http.StatusOK, "config.html", cd)
}
