package rbac

import (
	"context"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SimulateRequest represents the request payload for the simulation.
type SimulateRequest struct {
	Username  string   `json:"username" binding:"required"`
	RoleName  string   `json:"roleName" binding:"required"`
	Actions   []string `json:"actions" binding:"required"`
	Resources []string `json:"resources" binding:"required"`
	Namespace string   `json:"namespace" binding:"required"`
}

// SimulateResponse represents the response payload for the simulation.
type SimulateResponse struct {
	HasPermission bool              `json:"hasPermission"`
	Details       map[string]string `json:"details"`
}

// SimulateHandler handles the simulation of role assignments.
func SimulateHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SimulateRequest
		if err := c.Bind(&req); err != nil {
			utils.Logger.Error("Invalid request payload", zap.Error(err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload: " + err.Error()})
		}

		utils.Logger.Info("Simulation request received", zap.String("username", req.Username), zap.String("roleName", req.RoleName), zap.String("namespace", req.Namespace))

		// Validate inputs
		if !isValidUsername(req.Username) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid username"})
		}
		for _, action := range req.Actions {
			if !isValidAction(action) {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action: " + action})
			}
		}
		for _, resource := range req.Resources {
			if !isValidResource(resource, clientset) {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid resource: " + resource})
			}
		}

		// Fetch the user's current roles and role bindings
		roleBindings, err := clientset.RbacV1().RoleBindings(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			utils.Logger.Error("Error fetching role bindings", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching role bindings: " + err.Error()})
		}

		clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			utils.Logger.Error("Error fetching cluster role bindings", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching cluster role bindings: " + err.Error()})
		}

		// Fetch all roles and cluster roles
		roles, err := clientset.RbacV1().Roles(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			utils.Logger.Error("Error fetching roles", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching roles: " + err.Error()})
		}

		clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			utils.Logger.Error("Error fetching cluster roles", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching cluster roles: " + err.Error()})
		}

		// Validate role name
		if !isValidRoleName(req.RoleName, roles.Items, clusterRoles.Items) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role name"})
		}

		// Simulate the role assignment
		hasPermission, details := simulateRoleAssignment(req.Username, req.RoleName, req.Actions, req.Resources, roleBindings.Items, clusterRoleBindings.Items, roles.Items, clusterRoles.Items)

		utils.Logger.Info("Simulation result", zap.String("username", req.Username), zap.String("roleName", req.RoleName), zap.String("namespace", req.Namespace), zap.Bool("hasPermission", hasPermission))

		// Return the result
		return c.JSON(http.StatusOK, SimulateResponse{HasPermission: hasPermission, Details: details})
	}
}

// simulateRoleAssignment simulates the role assignment and checks if the user has the necessary permissions.
func simulateRoleAssignment(username, roleName string, actions, resources []string, roleBindings []rbacv1.RoleBinding, clusterRoleBindings []rbacv1.ClusterRoleBinding, roles []rbacv1.Role, clusterRoles []rbacv1.ClusterRole) (bool, map[string]string) {
	details := make(map[string]string)
	hasPermission := true

	for _, action := range actions {
		for _, resource := range resources {
			permission := false
			// Check current role bindings
			for _, rb := range roleBindings {
				if rb.RoleRef.Name == roleName {
					for _, subject := range rb.Subjects {
						if subject.Kind == rbacv1.UserKind && subject.Name == username {
							// Fetch the associated role
							for _, role := range roles {
								if role.Name == rb.RoleRef.Name && role.Namespace == rb.Namespace {
									// Check if the role has the necessary permissions
									for _, rule := range role.Rules {
										if contains(rule.Verbs, action) && contains(rule.Resources, resource) {
											permission = true
											break
										}
									}
								}
							}
						}
					}
				}
			}

			// Check current cluster role bindings
			for _, crb := range clusterRoleBindings {
				if crb.RoleRef.Name == roleName {
					for _, subject := range crb.Subjects {
						if subject.Kind == rbacv1.UserKind && subject.Name == username {
							// Fetch the associated cluster role
							for _, cr := range clusterRoles {
								if cr.Name == crb.RoleRef.Name {
									// Check if the cluster role has the necessary permissions
									for _, rule := range cr.Rules {
										if contains(rule.Verbs, action) && contains(rule.Resources, resource) {
											permission = true
											break
										}
									}
								}
							}
						}
					}
				}
			}

			if !permission {
				hasPermission = false
				details[action+"-"+resource] = "User does not have the necessary permissions"
			} else {
				details[action+"-"+resource] = "User has the necessary permissions"
			}
		}
	}

	return hasPermission, details
}

// contains checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Add validation functions
func isValidUsername(username string) bool {
	// Implement logic to check if the username exists
	// For example, you can fetch the user from a database or an external system
	// Here, we assume a function auth.GetAllUsers() that returns all users
	users, err := auth.GetAllUsers()
	if err != nil {
		return false
	}
	for _, user := range users {
		if user.Username == username {
			return true
		}
	}
	return false
}

func isValidRoleName(roleName string, roles []rbacv1.Role, clusterRoles []rbacv1.ClusterRole) bool {
	for _, role := range roles {
		if role.Name == roleName {
			return true
		}
	}
	for _, cr := range clusterRoles {
		if cr.Name == roleName {
			return true
		}
	}
	return false
}

func isValidAction(action string) bool {
	// List of valid Kubernetes verbs
	validActions := []string{"get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"}
	return contains(validActions, action)
}

func isValidResource(resource string, clientset *kubernetes.Clientset) bool {
	// Fetch the list of API resources
	apiResources, err := clientset.Discovery().ServerPreferredResources()
	if err != nil {
		return false
	}

	for _, resourceList := range apiResources {
		for _, apiResource := range resourceList.APIResources {
			if apiResource.Name == resource {
				return true
			}
		}
	}
	return false
}
