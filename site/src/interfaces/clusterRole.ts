export interface ClusterRole {
  apiVersion: string;
  kind: string;
  metadata: ClusterRoleMetadata;
  rules: ClusterRoleRules[];
}

export interface ClusterRoleMetadata {
  name: string;
  uid: string;
  resourceVersion: string;
  creationTimestamp: string;
}

export interface ClusterRoleRules {
  apiGroups: string[];
  resources: string[];
  verbs: string[];
}
