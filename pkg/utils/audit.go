package utils

import "log"

// LogAuditEvent logs an audit event for role and role binding changes.
func LogAuditEvent(action, resourceName, namespace string) {
	log.Printf("Audit event: action=%s, resource=%s, namespace=%s", action, resourceName, namespace)
}
