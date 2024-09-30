export interface User {
  metadata: {
    name: string;
    uid: string;
    resourceVersion: string;
    creationTimestamp: string;
  };
  groups: string[];
}
