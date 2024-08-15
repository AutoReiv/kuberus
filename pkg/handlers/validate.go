package handlers

import (
	"bytes"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

// validateKubernetesYAML validates the YAML content for any Kubernetes resource
func validateKubernetesYAML(yamlContent []byte) error {
	// Create a YAML decoder
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlContent), 100)

	// Create a runtime.Scheme and a codec factory
	s := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(s)

	// Add the Kubernetes API types to the scheme
	if err := scheme.AddToScheme(s); err != nil {
		return fmt.Errorf("failed to add Kubernetes API types to scheme: %v", err)
	}

	// Decode the YAML content into a runtime.Object
	var obj runtime.Object
	if err := decoder.Decode(&obj); err != nil {
		return fmt.Errorf("failed to decode YAML: %v", err)
	}

	// Encode the object back to JSON to ensure it's valid
	encoder := codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...)
	if _, err := runtime.Encode(encoder, obj); err != nil {
		return fmt.Errorf("failed to encode object: %v", err)
	}

	return nil
}
