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
  private async fetch(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<any> {
    const headers = {
      "Content-Type": "application/json",
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

  async deleteRoles(namespace: string, name: string) {
    const response = await this.fetch(
      ENDPOINTS.RBAC.ROLES.DELETE(namespace, name),
      {
        method: "DELETE",
      }
    );

    return response;
  }

  async getRoleDetails(namespace: string, name: string) {
    return this.fetch(ENDPOINTS.RBAC.ROLES.DETAILS(namespace, name));
  }

  async createRoles(roleData: any) {
    return this.fetch(ENDPOINTS.RBAC.ROLES.CREATE, {
      method: "POST",
      body: JSON.stringify(roleData),
    });
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

  async getRoleBindingsDetails(namespace: string, name: string) {
    return this.fetch(ENDPOINTS.RBAC.ROLEBINDINGS.DETAILS(namespace, name));
  }

  async deleteRoleBindings(namespace: string, name: string) {
    const response = await this.fetch(
      ENDPOINTS.RBAC.ROLEBINDINGS.DELETE(namespace, name),
      {
        method: "DELETE",
      }
    );

    return response;
  }

  async createRoleBindings(roleBindingData: any) {
    return this.fetch(ENDPOINTS.RBAC.ROLEBINDINGS.BASE, {
      method: "POST",
      body: JSON.stringify(roleBindingData),
    });
  }

  // ClusterRole-related methods
  async getClusterRoles() {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLES.BASE);
  }

  async getClusterRolesDetails(name: string) {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLES.DETAILS(name));
  }

  async createClusterRoles(clusterRoleData: any) {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLES.CREATE, {
      method: "POST",
      body: JSON.stringify(clusterRoleData),
    });
  }

  async deleteClusterRoles(name: string) {
    const response = await this.fetch(
      ENDPOINTS.RBAC.CLUSTERROLES.DELETE(name),
      {
        method: "DELETE",
      }
    );

    return response;
  }

  // ClusterRoleBinding-related methods
  async getClusterRoleBindings() {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLEBINDINGS.BASE);
  }

  async getClusterRoleBindingsDetails(name: string) {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLEBINDINGS.DETAILS(name));
  }

  async createClusterRoleBindings(clusterRoleBindingData: any) {
    return this.fetch(ENDPOINTS.RBAC.CLUSTERROLEBINDINGS.CREATE, {
      method: "POST",
      body: JSON.stringify(clusterRoleBindingData),
    });
  }

  async deleteClusterRoleBindings(name: string) {
    const response = await this.fetch(
      ENDPOINTS.RBAC.CLUSTERROLEBINDINGS.DELETE(name),
      {
        method: "DELETE",
      }
    );

    return response;
  }

  // User-related methods
  async getUsers() {
    return this.fetch(ENDPOINTS.USER_MANAGEMENT.USERS.BASE);
  }

  async getAdminCreateUser(userData: any) {
    return this.fetch(ENDPOINTS.ADMIN.CREATE, {
      method: "POST",
      body: JSON.stringify(userData),
    });
  }

  async getAdminUsers() {
    return this.fetch(ENDPOINTS.ADMIN.LIST, {
      method: "GET",
    });
  }

  async getAdminUserDelete(username: string) {
    return this.fetch(ENDPOINTS.ADMIN.DELETE, {
      method: "DELETE",
    });
  }

  async getAdminUserUpdate(username: string, userData: any) {
    return this.fetch(ENDPOINTS.ADMIN.UPDATE, {
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

  async createNamespaces(namespaceData: any) {
    return this.fetch(ENDPOINTS.K8S_RESOURCES.NAMESPACES, {
      method: "POST",
      body: JSON.stringify({ metadata: { name: namespaceData } }),
    });
  }

  async deleteNamespaces(namespace) {
    const response = await this.fetch(
      ENDPOINTS.K8S_RESOURCES.DELETENAMESPACE(namespace),
      {
        method: "DELETE",
      }
    );

    return response;
  }

  // Audit-related methods
  async getAuditLogs() {
    return this.fetch(ENDPOINTS.AUDIT.LOGS, {
      method: "GET",
    });
  }

  async getServiceAccounts() {
    return this.fetch(ENDPOINTS.USER_MANAGEMENT.SERVICEACCOUNTS.BASE);
  }

  async getServiceAccountDetails(namespace: string, name: string) {
    return this.fetch(
      ENDPOINTS.USER_MANAGEMENT.SERVICEACCOUNTS.DETAILS(namespace, name)
    );
  }

  async getGroups() {
    return this.fetch(ENDPOINTS.USER_MANAGEMENT.GROUPS.BASE);
  }

  async getGroupDetails(name: string) {
    return this.fetch(ENDPOINTS.USER_MANAGEMENT.GROUPS.DETAILS(name));
  }

  async getUserDetails(username: string) {
    return this.fetch(ENDPOINTS.USER_MANAGEMENT.USERS.DETAILS(username));
  }
}

export const apiClient = new ApiClient();
