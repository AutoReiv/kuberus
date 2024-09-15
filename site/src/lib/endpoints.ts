// Constants for common URL segments
const ALL_NAMESPACES = 'all';

/**
 * API Endpoints for the application
 */
export const ENDPOINTS = {
  // Kubernetes resource endpoints
  K8S_RESOURCES: {
    RESOURCES: '/api/resources',
    NAMESPACES: '/api/namespaces',
  },

  // RBAC-related endpoints
  RBAC: {
    ROLES: {
      BASE: `/api/roles?namespace=${ALL_NAMESPACES}`,
      DETAILS: (namespace: string, name: string) =>
        `/api/roles/details?roleName=${name}&namespace=${namespace}`,
      UPDATE: (namespace: string, name: string) =>
        `/api/roles?namespace=${namespace}&name=${name}`,
    },
    ROLEBINDINGS: {
      BASE: '/api/rolebindings',
      DETAILS: (namespace: string, name: string) =>
        `/api/rolebinding/details?name=${name}&namespace=${namespace}`,
    },
    CLUSTERROLES: {
      BASE: '/api/clusterroles',
      DETAILS: (name: string) => `/api/clusterroles/details?clusterRoleName=${name}`,
    },
    CLUSTERROLEBINDINGS: {
      BASE: '/api/clusterrolebindings',
      DETAILS: (name: string) => `/api/clusterrolebindings/details?name=${name}`,
    },
  },

  // User and group management endpoints
  USER_MANAGEMENT: {
    USERS: {
      BASE: '/api/users',
      DETAILS: '/api/user-details',
    },
    GROUPS: {
      BASE: '/api/groups',
      DETAILS: '/api/group-details',
    },
  },

  ADMIN: {
    CREATION: () => '/admin/create',
    USERS: (username?: string) => '/admin/users',
  },

  // Audit and logging endpoints
  AUDIT: {
    LOGS: '/api/audit-logs',
  },
};