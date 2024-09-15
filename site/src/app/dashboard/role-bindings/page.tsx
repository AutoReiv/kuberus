"use client";
import { Skeleton } from "@/components/ui/skeleton";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import DataTable from "./_components/DataTable";
import { apiClient } from "@/lib/apiClient";

const RoleBindings = () => {
  // Get Roles
  const { data: roleBindings, isPending: isPendingRoles } = useQuery({
    queryKey: ["roleBindings"],
    queryFn: () => apiClient.getRoleBindings()
  });
  return (
    <div className="flex w-full flex-col">
      {isPendingRoles ? (
        <Skeleton className="h-full w-100 m-4"></Skeleton>
      ) : (
        <DataTable roleBindings={roleBindings}></DataTable>
      )}
    </div>
  );
};

export default RoleBindings;
