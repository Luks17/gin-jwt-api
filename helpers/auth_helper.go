package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	user_type := c.GetString("user_type")

	err = nil

	if user_type != role {
		err = errors.New("Unauthorized acces to this resource")
		return err
	}

	return err
}

func MatchUserTypeToUUID(c *gin.Context, user_id string) (err error) {
	user_type := c.GetString("user_type")
	uid := c.GetString("user_id")

	err = nil

	if user_type == "USER" && uid != user_id {
		err = errors.New("Unauthorized access to this resource")
		return err
	}

	err = CheckUserType(c, user_type)
	return err
}
