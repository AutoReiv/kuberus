"use client";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useMutation, useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { useState } from "react";
import {
  Activity,
  ArrowUpRight,
  CreditCard,
  DollarSign,
  Package2,
  Users,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import UserManagement from "./_components/UserManagement";
import { apiClient } from "@/lib/apiClient";
import { DataTable } from "./_components/DataTable";
import { CreateUserDialog } from "./_components/CreateUserDialog";

interface AuditLog {
  id: number;
  action: string;
  resource_name: string;
  namespace: string;
  timestamp: string;
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

  console.log("auditLogs", auditLogs); 
  const { data: users, isFetched: usersFetched } = useQuery({
    queryKey: ["users"],
    queryFn: () => apiClient.getAdminUsers(),
  });

  console.log("users", users);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

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
          {usersFetched ? <DataTable users={users} /> : <div>Loading...</div>}
          <Card>
            <CardHeader>
              <CardTitle>Audit Log</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-8">
              <ul className="audit-logs-list space-y-4">
                {auditLogs
                  ? auditLogs.map((log) => (
                      <li
                        key={log.id}
                        className="audit-log-item p-4 border rounded-lg shadow-sm"
                      >
                        <div className="flex items-center justify-between text-sm py-2 border-b last:border-b-0">
                          <div className="flex items-center space-x-3">
                            <Badge
                              variant={getActionVariant(log.action)}
                              className="w-16 justify-center"
                            >
                              {log.action}
                            </Badge>
                            <span
                              className="font-medium truncate max-w-[150px]"
                              title={log.resource_name}
                            >
                              {log.resource_name}
                            </span>
                          </div>
                          <div className="flex items-center space-x-2 text-muted-foreground">
                            <span className="text-xs bg-secondary px-2 py-1 rounded">
                              {log.namespace}
                            </span>
                            <span
                              title={new Date(log.timestamp).toLocaleString()}
                            >
                              {formatRelativeTime(log.timestamp)}
                            </span>
                          </div>
                        </div>
                      </li>
                    ))
                  : "No Logs"}
              </ul>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
};

export default Admin;
