"use client";

import { useQuery } from "@tanstack/react-query";
import React from "react";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  ArrowUpDown,
  Clock,
  Copy,
  FileCode,
  FileText,
  MoreHorizontal,
  RefreshCw,
  UserMinus,
  UserPlus,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { Checkbox } from "@/components/ui/checkbox";
import { ColumnDef } from "@tanstack/react-table";
import { SkeletonPage } from "@/components/SkeletonPage";

interface RoleBinding {
  metadata: {
    name: string;
    namespace: string;
    uid: string;
    resourceVersion: string;
    creationTimestamp: string;
    managedFields: {
      manager: string;
      operation: string;
      apiVersion: string;
      time: string;
      fieldsType: string;
      fieldsV1: {
        [key: string]: any;
      };
    }[];
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

const RoleBindings = () => {
  // Get Roles
  const { data: roleBindings, isLoading } = useQuery({
    queryKey: ["roleBindings"],
    queryFn: () => apiClient.getRoleBindings(),
  });

  const columns: ColumnDef<RoleBinding>[] = [
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={table.getIsAllPageRowsSelected()}
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "metadata.name",
      id: "name", // Add this line
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
    },
    {
      accessorKey: "metadata.namespace",
      header: "Namespace",
    },
    {
      accessorKey: "roleRef.name",
      header: "Role",
    },
    {
      accessorKey: "subjects",
      header: "Subjects",
      cell: ({ row }) => {
        const subjects = row.original.subjects;
        return (
          <div>
            {subjects.map((subject, index) => (
              <div key={index}>{`${subject.kind}: ${subject.name}`}</div>
            ))}
          </div>
        );
      },
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
      cell: ({ row }) => {
        return new Date(
          row.original.metadata.creationTimestamp
        ).toLocaleString();
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const roleBinding = row.original;
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
                  href={`/dashboard/role-bindings/${roleBinding.metadata.namespace}/${roleBinding.metadata.name}`}
                >
                  View
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => console.log("Clone", roleBinding.metadata.name)}
              >
                <Copy className="mr-2 h-4 w-4" />
                Clone
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Add Subject", roleBinding.metadata.name)
                }
              >
                <UserPlus className="mr-2 h-4 w-4" />
                Add Subject
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Remove Subject", roleBinding.metadata.name)
                }
              >
                <UserMinus className="mr-2 h-4 w-4" />
                Remove Subject
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Change Role", roleBinding.metadata.name)
                }
              >
                <RefreshCw className="mr-2 h-4 w-4" />
                Change Role
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Export YAML", roleBinding.metadata.name)
                }
              >
                <FileCode className="mr-2 h-4 w-4" />
                Export YAML
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Temporary Access", roleBinding.metadata.name)
                }
              >
                <Clock className="mr-2 h-4 w-4" />
                Temporary Access
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => console.log("Edit", roleBinding.metadata.name)}
              >
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => console.log("Delete", roleBinding.metadata.name)}
              >
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  return (
    <div className="flex w-full flex-col">
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={roleBindings}
          columns={columns}
          title="Role Bindings"
          description="Manage and configure service accounts to control access and authentication for applications and processes within your Kubernetes cluster"
          // Add in row action to route to details
        ></GenericDataTable>
      )}
    </div>
  );
};

export default RoleBindings;
