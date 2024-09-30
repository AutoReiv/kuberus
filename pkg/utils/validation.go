package utils

import (
	"errors"

	rbacv1 "k8s.io/api/rbac/v1"
)

// ValidateRole ensures that the role is valid.
func ValidateRole(role *rbacv1.Role) error {
	if role.Name == "" {
		return errors.New("role name is required")
	}
	if len(role.Rules) == 0 {
		return errors.New("at least one rule is required")
	}
	return nil
}

// ValidateRoleBinding ensures that the role binding is valid.
func ValidateRoleBinding(roleBinding *rbacv1.RoleBinding) error {
	if roleBinding.Name == "" {
		return errors.New("role binding name is required")
	}
	if roleBinding.RoleRef.Name == "" {
		return errors.New("role reference name is required")
	}
	if len(roleBinding.Subjects) == 0 {
		return errors.New("at least one subject is required")
	}
	return nil
}

// ValidateClusterRoleBinding ensures that the cluster role binding is valid.
func ValidateClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	if clusterRoleBinding.Name == "" {
		return errors.New("cluster role binding name is required")
	}
	if clusterRoleBinding.RoleRef.Name == "" {
		return errors.New("role reference name is required")
	}
	if len(clusterRoleBinding.Subjects) == 0 {
		return errors.New("at least one subject is required")
	}
	return nil
}
