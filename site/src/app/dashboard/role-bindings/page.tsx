"use client";
import { Skeleton } from "@/components/ui/skeleton";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import DataTable from "./_components/DataTable";

/**
 * Fetches a list of namespaces from the API.
 * @returns {Promise<any>} - A promise that resolves to the response data from the API.
 */
const getRoleBindings = async () => {
  const URL = "http://localhost:8080/api/rolebindings?namespaces=all";
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

const RoleBindings = () => {
  // Get Roles
  const { data: roleBindings, isPending: isPendingRoles } = useQuery({
    queryKey: ["roleBindings"],
    queryFn: getRoleBindings,
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
