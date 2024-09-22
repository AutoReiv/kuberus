package rbac

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
	"rbac/pkg/utils"
)

// APIResourcesHandler handles retrieving all Kubernetes API resources.
func APIResourcesHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Create a discovery client to list available API resources
		discoveryClient := clientset.Discovery()

		// Retrieve the list of preferred API resources
		apiResources, err := discoveryClient.ServerPreferredResources()
		if err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error retrieving API resources", err, "Failed to retrieve API resources")
		}

		// Collect the names of the API resources
		var resources []string
		for _, resourceList := range apiResources {
			for _, resource := range resourceList.APIResources {
				resources = append(resources, resource.Name+" ("+resourceList.GroupVersion+")")
			}
		}

		// Write the list of API resources as a JSON response
		return c.JSON(http.StatusOK, map[string][]string{"resources": resources})
	}
}
