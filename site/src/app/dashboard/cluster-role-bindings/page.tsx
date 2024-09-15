"use client";

import { useQuery } from "@tanstack/react-query";
import React from "react";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { ArrowUpDown } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { SkeletonPage } from "@/components/SkeletonPage";

interface ClusterRoleBinding {
  metadata: {
    name: string;
    uid: string;
    creationTimestamp: string;
  };
  subjects: {
    kind: string;
    name: string;
  }[];
  roleRef: {
    kind: string;
    name: string;
  };
}

const ClusterRoleBindings = () => {
  // Get Cluster Role Bindings
  const { data: clusterRoleBindings, isPending: isPendingRoles } = useQuery({
    queryKey: ["clusterRoleBindings"],
    queryFn: () => apiClient.getClusterRoleBindings(),
  });

  const columns: ColumnDef<ClusterRoleBinding>[] = [
    {
      accessorKey: "name",
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
      cell: ({ row }) => <div>{row.original.metadata.name}</div>,
    },
    {
      accessorKey: "roleRef.name",
      header: "Role",
      cell: ({ row }) => <div>{row.original.roleRef.name}</div>,
    },
    {
      accessorKey: "subjects",
      header: "Subjects",
      cell: ({ row }) => (
        <div className="flex flex-wrap gap-1">
          {row.original.subjects.map((subject, index) => (
            <Badge key={index} variant="secondary">
              {subject.kind}: {subject.name}
            </Badge>
          ))}
        </div>
      ),
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
      cell: ({ row }) => (
        <div>
          {new Date(row.original.metadata.creationTimestamp).toLocaleString()}
        </div>
      ),
    },
  ];

  return (
    <div className="flex w-full flex-col">
      {isPendingRoles ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={clusterRoleBindings}
          columns={columns}
          title="Cluster Role Bindings"
          description="Manage and view Cluster Role Bindings, which define the association between a set of permissions and a user or set of users across the entire cluster. These bindings are crucial for controlling access and maintaining security in your Kubernetes environment."
        ></GenericDataTable>
      )}
    </div>
  );
};

export default ClusterRoleBindings;
