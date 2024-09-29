export interface Namespace {
  metadata: NamespaceMetadata;
}

export interface NamespaceMetadata {
  name: string;
  uid: string;
  resourceVersion: string;
  labels: Record<string, string>;
  managedFields: any[];
}
