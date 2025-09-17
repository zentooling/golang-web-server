package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/infra"
	"github.com/uberswe/golang-base-project/models"
	"golang.org/x/crypto/bcrypt"
)

// ResetPasswordPageData defines additional data needed to render the reset password page
type ResetPasswordPageData struct {
	PageData
	Token string
}

// ResetPassword renders the HTML page for resetting the users password
func ResetPassword(c *gin.Context) {
	token := c.Param("token")
	pdPre := DefaultPageData(c)
	pdPre.Title = pdPre.Trans("Reset Password")
	pd := ResetPasswordPageData{
		PageData: pdPre,
		Token:    token,
	}
	c.HTML(http.StatusOK, "resetpassword.gohtml", pd)
}

// ResetPasswordPost handles post request used to reset users passwords
func ResetPasswordPost(c *gin.Context) {
	pdPre := DefaultPageData(c)
	passwordError := pdPre.Trans("Your password must be 8 characters in length or longer")
	resetError := pdPre.Trans("Could not reset password, please try again")

	token := c.Param("token")
	pdPre.Title = pdPre.Trans("Reset Password")
	pd := ResetPasswordPageData{
		PageData: pdPre,
		Token:    token,
	}
	password := c.PostForm("password")

	if len(password) < 8 {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: passwordError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.gohtml", pd)
		return
	}

	forgotPasswordToken := models.Token{
		Value: token,
		Type:  models.TokenPasswordReset,
	}

	db := infra.LairInstance().GetDb()

	res := db.Where(&forgotPasswordToken).First(&forgotPasswordToken)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.gohtml", pd)
		return
	}

	if forgotPasswordToken.HasExpired() {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.gohtml", pd)
		return
	}

	user := models.User{}
	user.ID = uint(forgotPasswordToken.ModelID)
	res = db.Where(&user).First(&user)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.gohtml", pd)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("ResetPasswordPost", "error", err)
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.gohtml", pd)
		return
	}

	user.Password = string(hashedPassword)

	res = db.Save(&user)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.gohtml", pd)
		return
	}

	res = db.Delete(&forgotPasswordToken)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.gohtml", pd)
		return
	}

	pd.Messages = append(pd.Messages, Message{
		Type:    "success",
		Content: pdPre.Trans("Your password has successfully been reset."),
	})

	c.HTML(http.StatusOK, "resetpassword.gohtml", pd)
}
