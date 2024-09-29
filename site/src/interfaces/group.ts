export interface Group {
    metadata: {
      name: string;
      creationTimestamp: string;
      resourceVersion: string;
    };
    users: string[];
  }
  