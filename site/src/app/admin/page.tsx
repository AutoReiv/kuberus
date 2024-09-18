"use client";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { ArrowUpDown, Package2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";

interface AuditLog {
  id: number;
  action: string;
  resource_name: string;
  namespace: string;
  timestamp: string;
  hash: string;
}

export interface User {
  username: string;
}

const getActionVariant = (
  action: string
): "default" | "secondary" | "destructive" | "outline" => {
  switch (action.toLowerCase()) {
    case "create":
      return "default";
    case "update":
      return "secondary";
    case "delete":
      return "destructive";
    default:
      return "outline";
  }
};

const formatRelativeTime = (timestamp: string): string => {
  const now = new Date();
  const logTime = new Date(timestamp);
  const diffInSeconds = Math.floor((now.getTime() - logTime.getTime()) / 1000);

  if (diffInSeconds < 60) return `${diffInSeconds}s ago`;
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`;
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`;
  return `${Math.floor(diffInSeconds / 86400)}d ago`;
};

const Admin = () => {
  const {
    data: auditLogs,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["auditLog"],
    queryFn: () => apiClient.getAuditLogs(),
  });

  const { data: users, isFetched: usersFetched } = useQuery({
    queryKey: ["users"],
    queryFn: () => apiClient.getAdminUsers(),
  });

  const columns: ColumnDef<User>[] = [
    {
      accessorKey: "username",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Name
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
      cell: ({ row }) => <div>{row.original.username}</div>,
    },
  ];

  return (
    <div className="flex min-h-screen w-full flex-col">
      <header className="sticky top-0 flex h-16 items-center gap-4 border-b bg-background px-4 md:px-6">
        <nav className="hidden flex-col gap-6 text-lg font-medium md:flex md:flex-row md:items-center md:gap-5 md:text-sm lg:gap-6">
          <Link
            href="#"
            className="flex items-center gap-2 text-lg font-semibold md:text-base"
          >
            <Package2 className="h-6 w-6" />
            <span className="sr-only">Acme Inc</span>
          </Link>
          <Link
            href="#"
            className="text-foreground transition-colors hover:text-foreground"
          >
            Dashboard
          </Link>
          <Link
            href="#"
            className="text-muted-foreground transition-colors hover:text-foreground"
          >
            Orders
          </Link>
          <Link
            href="#"
            className="text-muted-foreground transition-colors hover:text-foreground"
          >
            Products
          </Link>
          <Link
            href="#"
            className="text-muted-foreground transition-colors hover:text-foreground"
          >
            Customers
          </Link>
          <Link
            href="#"
            className="text-muted-foreground transition-colors hover:text-foreground"
          >
            Analytics
          </Link>
        </nav>
      </header>
      <main className="flex flex-1 flex-col gap-4 p-4 md:gap-8 md:p-8">
        <div className="grid gap-4 md:gap-8 lg:grid-cols-2 xl:grid-cols-3">
          {isLoading ? (
            <p>Loading Table...</p>
          ) : error ? (
            <p>Error loading users</p>
          ) : (
            <GenericDataTable
              className="col-span-2"
              data={users}
              columns={columns}
              title="Users"
              description="Manage and view all registered users in the system"
            ></GenericDataTable>
          )}
          <Card>
            <CardHeader>
              <CardTitle>Audit Log</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-1">
              {isLoading ? (
                <p>Loading audit logs...</p>
              ) : error ? (
                <p>Error loading audit logs</p>
              ) : (
                auditLogs?.slice(0, 10).map((log: AuditLog) => (
                  <div
                    key={log.id}
                    className="flex justify-between items-center text-sm hover:bg-muted rounded-md p-1"
                  >
                    <span>{log.action}</span>
                    <span>{log.resource_name}</span>
                    <Badge variant={getActionVariant(log.action)}>
                      {formatRelativeTime(log.timestamp)}
                    </Badge>
                  </div>
                ))
              )}
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
};

export default Admin;
