"use client";

import { useQuery } from "@tanstack/react-query";
import React from "react";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { ArrowUpDown } from "lucide-react";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";

interface ClusterRole {
  metadata: {
    name: string;
    creationTimestamp: string;
  };
  rules: {
    apiGroups: string[];
    resources: string[];
    verbs: string[];
  }[];
}

const ClusterRoles = () => {
  // Get ClusterRoles
  const {
    data: clusterRoles,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["roles"],
    queryFn: () => apiClient.getClusterRoles(),
  });

  if (isError) {
    return <div>Error</div>;
  }

  const columns: ColumnDef<ClusterRole>[] = [
    {
      id: "metadata.name",
      accessorKey: "metadata.name",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Name
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
    },
    {
      accessorKey: "rules",
      header: "Rules",
      cell: ({ row }) => {
        const rules = row.original.rules;
        return <span>{rules.length} rule(s)</span>;
      },
    },
  ];

  return (
    <motion.div
      className="flex w-full flex-col"
      initial="initial"
      animate="animate"
      exit="exit"
      variants={pageVariants}
    >
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={clusterRoles}
          columns={columns}
          title="Cluster Roles"
          description="Cluster roles are used to grant permissions to users, groups, or service accounts across the entire cluster. They define the actions that can be performed on resources within the cluster. Cluster roles are typically used for global permissions that apply to all namespaces."
        ></GenericDataTable>
      )}
    </motion.div>
  );
};

export default ClusterRoles;
