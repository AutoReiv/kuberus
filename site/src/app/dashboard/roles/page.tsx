"use client";

import React from "react";
import DataTable from "./_components/DataTable";
import { useQuery } from "@tanstack/react-query";
import { Skeleton } from "@/components/ui/skeleton";

const Roles = () => {
  // Get Roles
  const { data: roles , isPending: isPendingRoles } = useQuery({
    queryKey: ["roles"],
    queryFn: async () => {
      const URL = "http://localhost:8080/api/roles?namespace=all";
      const response = await fetch(URL, {
        method: "GET",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        credentials: "include",
      });
      const data = await response.json();
      return data;
    },
  }); 

  // Get namespaces
  const { data: namespace } = useQuery({
    queryKey: ["namespace"],
    queryFn: async () => {
      const URL = "http://localhost:8080/api/namespaces";
      const response = await fetch(URL, {
        method: "GET",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      });
      const data = await response.json();
      return data;
    },
  });
  
  return (
    <div className="flex w-full flex-col">
      {isPendingRoles ? <Skeleton className="h-full w-100 m-4"></Skeleton> : <DataTable roles={roles.items} namespace={namespace}></DataTable>}
    </div>
  );
};

export default Roles;
