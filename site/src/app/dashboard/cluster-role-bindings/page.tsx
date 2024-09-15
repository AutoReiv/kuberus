'use client';
import { Skeleton } from "@/components/ui/skeleton";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import DataTable from "./_components/DataTable";
import { apiClient } from "@/lib/apiClient";

const ClusterRoleBindings = () => {
  // Get Cluster Role Bindings
  const { data: clusterRoleBindings, isPending: isPendingRoles } = useQuery({
    queryKey: ["clusterRoleBindings"],
    queryFn: () => apiClient.getClusterRoleBindings(),
  });
  return (
    <div className="flex w-full flex-col">
      {isPendingRoles ? (
        <Skeleton className="h-full w-100 m-4"></Skeleton>
      ) : (
        <DataTable clusterRoleBindings={clusterRoleBindings}></DataTable>
      )}
    </div>
  );
};

export default ClusterRoleBindings;
