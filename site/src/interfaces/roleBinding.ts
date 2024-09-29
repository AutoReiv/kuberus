export interface RoleBinding {
  metadata: {
    name: string;
    namespace: string;
    uid: string;
    resourceVersion: string;
    creationTimestamp: string;
  };
  subjects: {
    kind: string;
    apiGroup: string;
    name: string;
  }[];
  roleRef: {
    apiGroup: string;
    kind: string;
    name: string;
  };
}
