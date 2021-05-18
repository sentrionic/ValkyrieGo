package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// used to help extract validation errors
type invalidArgument struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// bindData is helper function, returns false if data is not bound
func bindData(c *gin.Context, req interface{}) bool {
	// Bind incoming json to struct and check for validation errors
	if err := c.ShouldBind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)

		if errs, ok := err.(validator.ValidationErrors); ok {
			// could probably extract this, it is also in middleware_auth_user
			var invalidArgs []invalidArgument

			for _, err := range errs {
				fmt.Println(err)
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					fmt.Sprintf("Error: Value must be %s %s, got: %s", err.Tag(), err.Param(), err.Value().(string)),
				})
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request parameters. See errors",
				"errors":  invalidArgs,
			})
			return false
		}
	}
	return true
}
