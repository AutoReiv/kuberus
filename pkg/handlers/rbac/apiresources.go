package rbac

import (
	"net/http"
	"rbac/pkg/utils"

	"k8s.io/client-go/kubernetes"
)

// APIResourcesHandler handles retrieving all Kubernetes API resources.
func APIResourcesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		discoveryClient := clientset.Discovery()
		apiResources, err := discoveryClient.ServerPreferredResources()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var resources []string
		for _, resourceList := range apiResources {
			for _, resource := range resourceList.APIResources {
				resources = append(resources, resource.Name+" ("+resourceList.GroupVersion+")")
			}
		}

		utils.WriteJSON(w, map[string][]string{"resources": resources})
	}
}
