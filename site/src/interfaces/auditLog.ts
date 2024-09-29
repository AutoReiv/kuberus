export interface AuditLog {
  action: string;
  resource: string;
  namespace: string;
  timestamp: string;
  hash: string;
}
