export interface ClusterRoleBinding {
  metadata: {
    name: string;
    uid: string;
    resourceVersion: string;
    creationTimestamp: string;
  };
  subjects: {
    kind: string;
    apiGroup: string;
    name: string;
    namespace?: string;
  }[];
  roleRef: {
    apiGroup: string;
    kind: string;
    name: string;
  };
}
