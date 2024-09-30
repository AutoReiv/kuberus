package server

import (
	"rbac/pkg/handlers/users"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

// registerAdminUserRoutes registers all the user-related routes under the /admin/users endpoint.
func registerAdminUserRoutes(e *echo.Echo, clientset *kubernetes.Clientset) {
	adminGroup := e.Group("/admin/users")
	adminGroup.POST("", func(c echo.Context) error { return users.HandleCreateUser(c, clientset) })
	adminGroup.PUT("", users.HandleUpdateUser)
	adminGroup.DELETE("", users.HandleDeleteUser)
	adminGroup.GET("", users.HandleListUsers)
	adminGroup.GET("/details", func(c echo.Context) error { return users.HandleGetUserDetails(c, clientset) })
}
