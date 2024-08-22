"use client";

import React from "react";
import DataTable from "./_components/DataTable";
import { useQuery } from "@tanstack/react-query";

const Roles = () => {
  // Get Roles
  const { data: roles , isPending: isPendingRoles, isError: isErrorRoles } = useQuery({
    queryKey: ["roles"],
    queryFn: async () => {
      const URL = "http://localhost:8080/api/roles?namespace=all";
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

  // Get namespaces
  const { data: namespace, isError: isErrorNamespace } = useQuery({
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
      {isPendingRoles ? <div>Loading...</div> : <DataTable roles={roles.items} namespace={namespace}></DataTable>}
    </div>
  );
};

export default Roles;
