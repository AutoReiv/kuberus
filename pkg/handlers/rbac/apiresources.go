package rbac

import (
	"net/http"
	"rbac/pkg/utils"

	"k8s.io/client-go/kubernetes"
)

// APIResourcesHandler handles retrieving all Kubernetes API resources.
// It uses the Kubernetes clientset to list the available API resources and returns them as a JSON response.
func APIResourcesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a discovery client to list available API resources
		discoveryClient := clientset.Discovery()

		// Retrieve the list of preferred API resources
		apiResources, err := discoveryClient.ServerPreferredResources()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Collect the names of the API resources
		var resources []string
		for _, resourceList := range apiResources {
			for _, resource := range resourceList.APIResources {
				resources = append(resources, resource.Name+" ("+resourceList.GroupVersion+")")
			}
		}

		// Write the list of API resources as a JSON response
		utils.WriteJSON(w, map[string][]string{"resources": resources})
	}
}
