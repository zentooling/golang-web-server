package routes

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/infra"
	"github.com/uberswe/golang-base-project/models"
)

// Activate handles requests used to activate a users account
func Activate(c *gin.Context) {

	pd := DefaultPageData(c)
	activationError := pd.Trans("Please provide a valid activation token")
	activationSuccess := pd.Trans("Account activated. You may now proceed to login to your account.")
	pd.Title = pd.Trans("Activate")
	token := c.Param("token")
	activationToken := models.Token{
		Value: token,
		Type:  models.TokenUserActivation,
	}

	db := infra.LairInstance().GetDb()

	res := db.Where(&activationToken).First(&activationToken)
	if res.Error != nil {
		pd.AddMessage(Error, activationError)
		slog.Error("Activate:TokenNotFound", "error", res.Error)
		c.HTML(http.StatusBadRequest, "activate.gohtml", pd)
		return
	}

	if activationToken.HasExpired() {
		pd.AddMessage(Error, activationError)
		slog.Info("Activate:TokenHasExpired", "error", res.Error)
		c.HTML(http.StatusBadRequest, "activate.gohtml", pd)
		return
	}

	user := models.User{}
	user.ID = uint(activationToken.ModelID)

	res = db.Where(&user).First(&user)
	if res.Error != nil {
		pd.AddMessage(Error, activationError)
		slog.Error("Activate:UserNotFound", "error", res.Error)
		c.HTML(http.StatusBadRequest, "activate.gohtml", pd)
		return
	}

	now := time.Now()
	user.ActivatedAt = &now

	res = db.Save(&user)
	if res.Error != nil {
		pd.AddMessage(Error, activationError)
		slog.Error("Activate:SaveUser", "error", res.Error)
		c.HTML(http.StatusBadRequest, "activate.gohtml", pd)
		return
	}

	// We don't need to check for an error here, even if it's not deleted it will not really affect application logic
	db.Delete(&activationToken)

	pd.AddMessage(Success, activationSuccess)
	slog.Info("Activate:Success", "token", activationToken.Value)
	c.HTML(http.StatusOK, "activate.gohtml", pd)
}
