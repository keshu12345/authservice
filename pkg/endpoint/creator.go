package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keshucs12345/authservice/constants"
	"github.com/keshucs12345/authservice/utils"
)

func Creator(c *gin.Context) {
	user, exists := c.Get(constants.ApiKeyUserCtxKey.String())
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf("user not found"))
		return
	}
	utils.RespondSuccess(c, http.StatusOK, "User found", user)
}
