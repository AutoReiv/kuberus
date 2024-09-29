export interface Metadata {
  name: string;
  namespace?: string;
  uid: string;
  resourceVersion: string;
  creationTimestamp: string;
}

export interface KubernetesResource {
  apiVersion: string;
  kind: string;
  metadata: Metadata;
}
