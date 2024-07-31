package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindJSON() is a wrapper for gin.Context.ShouldBindJSON that has better error handling capabilities (imo)
func BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		if err.Error() == "EOF" {
			return errors.New("request body cannot be empty")
		} else if validator_errors, ok := err.(validator.ValidationErrors); ok {
			out := make([]string, len(validator_errors))
			for i, fe := range validator_errors {
				out[i] = fe.Field() + ": " + fe.Tag()
			}
			return errors.New("invalid request: " + strings.Join(out, ", "))
		} else {
			return errors.New("invalid request: " + err.Error())
		}
	}
	return nil
}
