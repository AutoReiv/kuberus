export interface Role {
  metadata: RoleMetaData;
  rules: RoleRules[];
  active: boolean;
}

export interface RoleMetaData {
  name: string;
  namespace: string;
  uid?: string;
  resourceVersion?: string;
  creationTimestamp?: string;
  managedFields?: [];
}

export interface RoleRules {
  verbs?: string[];
  apiGroups?: string[];
  resources?: string[];
}
