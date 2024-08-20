package handlers

import (
	"encoding/json"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

// ValidateRoleJSON validates the JSON content and converts it to a Kubernetes Role object
func ValidateRoleJSON(body []byte) (*rbacv1.Role, error) {
	// Validate JSON content
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return nil, fmt.Errorf("invalid JSON content: %v", err)
	}

	// Convert JSON to Kubernetes object
	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer()
	obj, _, err := decoder.Decode(body, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid Kubernetes object: %v", err)
	}

	// Assuming the object is a Role, validate it
	role, ok := obj.(*rbacv1.Role)
	if !ok {
		return nil, fmt.Errorf("invalid Role object")
	}

	// Additional field validation
	if role.Name == "" || role.Namespace == "" {
		return nil, fmt.Errorf("role must have a name and namespace")
	}

	return role, nil
}
