export interface Role {
  metadata: RoleMetaData,
  rules: RoleRules[],
  active: boolean
}

interface RoleMetaData {
    name: string,
    namespace: string,
    uid: string,
    resourceVersion: string,
    creationTimestamp: string,
    managedFields: []
}

export interface RoleRules {
    verbs: [],
    apiGroups: [],
    resources: []
}

export interface Namespace {
    metadata: NamespaceMetadata
}

interface NamespaceMetadata {
    name: string,
    uid: string,
    resourceVersion: string,
    labels: {},
    managedFields :[]
}