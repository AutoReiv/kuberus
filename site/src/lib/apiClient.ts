import { ENDPOINTS } from "./endpoints";

const API_BASE_URL = "http://localhost:8080";

/**
 * ApiClient class for handling API requests
 */
class ApiClient {
  /**
   * Generic fetch method for making API requests
   * @param endpoint - API endpoint
   * @param options - Request options
   * @returns Promise with the JSON response
   */
  private async fetch(endpoint: string, options: RequestInit = {}): Promise<any> {
    const token = localStorage.getItem("authToken");
    const headers = {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...options.headers,
    };

    try {
      const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        ...options,
        headers,
      });

      if (!response.ok) {
        throw new Error(`API request failed: ${response.statusText}`);
      }

      return response.json();
    } catch (error) {
      console.error("API request error:", error);
      throw error;
    }
  }

  // Role-related methods
  async getRoles() {
    return this.fetch(ENDPOINTS.RBAC.ROLES.BASE);
  }

  async getRoleDetails(namespace: string, name: string) {
    return this.fetch(ENDPOINTS.RBAC.ROLES.DETAILS(namespace, name));
  }

  async updateRole(namespace: string, name: string, roleData: any) {
    return this.fetch(ENDPOINTS.RBAC.ROLES.UPDATE(namespace, name), {
      method: "PUT",
      body: JSON.stringify(roleData),
    });
  }

  async duplicateRole(newNamespace: string, newName: string, roleData: any) {
    return this.fetch(ENDPOINTS.RBAC.ROLES.UPDATE(newNamespace, newName), {
      method: "POST",
      body: JSON.stringify(roleData),
    });
  }

  // RoleBinding-related methods
  async getRoleBindings() {
    return this.fetch(ENDPOINTS.RBAC.ROLEBINDINGS.BASE);
  }

  async getRoleBindingDetails(namespace: string, name: string) {
    return this.fetch(ENDPOINTS.RBAC.ROLEBINDINGS.DETAILS(namespace, name));
  }

  // ClusterRole-related methods
  async getClusterRoles() {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLES.BASE);
  }

  async getClusterRoleDetails(name: string) {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLES.DETAILS(name));
  }

  // ClusterRoleBinding-related methods
  async getClusterRoleBindings() {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLEBINDINGS.BASE);
  }

  async getClusterRoleBindingDetails(name: string) {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLEBINDINGS.DETAILS(name));
  }

  // User-related methods
  async getUsers() {
    return this.fetch(ENDPOINTS.USER_MANAGEMENT.USERS.BASE);
  }

  async createUser(userData: any) {
    return this.fetch(ENDPOINTS.ADMIN.USERS(), {
      method: "POST",
      body: JSON.stringify(userData),
    });
  }

  async getAdminUsers() {
    return this.fetch(ENDPOINTS.ADMIN.USERS(), {
      method: "GET",
    });
  }

  async getAdminUserDelete(username: string) {
    return this.fetch(ENDPOINTS.ADMIN.USERS(username), {
      method: "DELETE",
    });
  }

  async getAdminUserUpdate(username: string, userData: any) {
    return this.fetch(ENDPOINTS.ADMIN.USERS(username), {
      method: "PUT",
      body: JSON.stringify(userData),
    });
  }

  // Resource-related methods
  async getResources() {
    return this.fetch(ENDPOINTS.K8S_RESOURCES.RESOURCES);
  }

  async getNamespaces() {
    return this.fetch(ENDPOINTS.K8S_RESOURCES.NAMESPACES);
  }

  // Audit-related methods
  async getAuditLogs() {
    return this.fetch(ENDPOINTS.AUDIT.LOGS);
  }
}

export const apiClient = new ApiClient();
