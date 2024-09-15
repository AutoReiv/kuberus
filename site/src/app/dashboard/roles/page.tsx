"use client";

import React from "react";
import DataTable from "./_components/DataTable";
import { useQuery } from "@tanstack/react-query";
import { Skeleton } from "@/components/ui/skeleton";
import { apiClient } from "@/lib/apiClient";

/**
 * Renders a component that displays a list of roles and namespaces.
 * 
 * The component uses the `useQuery` hook from `@tanstack/react-query` to fetch the list of roles and namespaces from the API.
 * If the data is still being fetched, a skeleton loader is displayed. Otherwise, a `DataTable` component is rendered with the fetched roles and namespaces.
 */
const Roles = () => {
  // Get Roles
  const { data: roles, isLoading, isError,  } = useQuery({
    queryKey: ["roles"],
    queryFn: () => apiClient.getRoles()
  }); 

  // Get Namespaces
  const { data: namespace } = useQuery({
    queryKey: ["namespace"],
    queryFn: ()=> apiClient.getNamespaces()
  });
  
  if(isError){
    return <div>Error</div>
  }

  return (
    <div className="flex w-full flex-col">
      {isLoading ? <Skeleton className="h-full w-100 m-4"></Skeleton> : <DataTable roles={roles} namespace={namespace}></DataTable>}
    </div>
  );
};

export default Roles;
