"use client";

import React from "react";
import { apiClient } from "@/lib/apiClient";
import { useQuery } from "@tanstack/react-query";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Copy,
  Edit,
  FileText,
  MoreHorizontal,
  Trash
} from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { format } from "date-fns";
import { SkeletonPage } from "@/components/SkeletonPage";

interface ServiceAccount {
  metadata: {
    name: string;
    namespace: string;
    creationTimestamp: string;
  };
}

const ServiceAccounts = () => {
  // Get Service Accounts
  const {
    data: serviceAccounts,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["serviceAccounts"],
    queryFn: () => apiClient.getServiceAccounts(),
  });

  const columns: ColumnDef<ServiceAccount>[] = [
    {
      accessorKey: "metadata.name",
      header: "Name",
    },
    {
      accessorKey: "metadata.namespace",
      header: "Namespace",
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
      cell: ({ row }) => {
        return format(
          new Date(row.original.metadata.creationTimestamp),
          "MM/dd - hh:mm:ss a"
        );
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const serviceAccount = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuItem>
                <FileText className="mr-2 h-4 w-4" />
                <Link
                  href={`/dashboard/service-accounts/${serviceAccount.metadata.namespace}/${serviceAccount.metadata.name}`}
                >
                  View
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Clone", serviceAccount.metadata.name)
                }
              >
                <Copy className="mr-2 h-4 w-4" />
                Clone
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Edit", serviceAccount.metadata.name)
                }
              >
                <Edit className="mr-2 h-4 w-4" />
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Delete", serviceAccount.metadata.name)
                }
              >
                <Trash className="mr-2 h-4 w-4" />
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  if (isError) {
    return <div>Error</div>;
  }

  return (
    <div className="flex w-full flex-col">
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={serviceAccounts}
          columns={columns}
          title="Service Accounts"
          description="Manage and configure service accounts in your Kubernetes cluster"
        />
      )}
    </div>
  );
};

export default ServiceAccounts;
