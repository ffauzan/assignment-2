package rest

import (
	"fga-asg-2/pkg/order"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewRouter(orderService order.Service) *gin.Engine {
	r := gin.Default()
	RegisterOrderHandler(r, orderService)
	return r
}

// Function to send error response
func SendErrorResponse(c *gin.Context, err error, code int) {
	c.JSON(code, BaseResponse{
		Status:  "error",
		Message: err.Error(),
		Data:    nil,
	})
}
