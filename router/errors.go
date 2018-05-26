package router

import (
	"github.com/gin-gonic/gin"
)

var errors = struct {
	BadRequest   gin.H
	Unauthorized gin.H
	Internal     gin.H
}{
	BadRequest:   gin.H{"error": "StatusBadRequest"},
	Unauthorized: gin.H{"error": "StatusUnauthorized"},
	Internal:     gin.H{"error": "StatusInternalServerError"},
}
