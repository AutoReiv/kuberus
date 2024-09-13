'use client';
import { Skeleton } from "@/components/ui/skeleton";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import DataTable from "./_components/DataTable";

/**
 * Fetches a list of namespaces from the API.
 * @returns {Promise<any>} - A promise that resolves to the response data from the API.
 */
const getClusterRoleBindings = async () => {
  const URL = "http://localhost:8080/api/clusterrolebindings?namespaces=all";
  const response = await fetch(URL, {
    method: "GET",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  });
  const data = await response.json();
  return data;
};

const ClusterRoleBindings = () => {
  // Get Cluster Role Bindings
  const { data: clusterRoleBindings, isPending: isPendingRoles } = useQuery({
    queryKey: ["clusterRoleBindings"],
    queryFn: getClusterRoleBindings,
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
