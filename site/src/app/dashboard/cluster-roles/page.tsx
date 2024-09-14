"use client";

import { Skeleton } from "@/components/ui/skeleton";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import DataTable from "./_components/DataTable";
import { apiClient } from "@/lib/apiClient";

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

  return (
    <div className="flex w-full flex-col">
      {isLoading ? (
        <Skeleton className="h-full w-100 m-4"></Skeleton>
      ) : (
        <DataTable clusterRoles={clusterRoles}></DataTable>
      )}
    </div>
  );
};

export default ClusterRoles;
