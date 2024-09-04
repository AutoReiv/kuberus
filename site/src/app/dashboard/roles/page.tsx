"use client";

import React from "react";
import DataTable from "./_components/DataTable";
import { useQuery } from "@tanstack/react-query";
import { Skeleton } from "@/components/ui/skeleton";

/**
 * Fetches a list of roles from the API.
 * @returns {Promise<any>} - A promise that resolves to the response data from the API.
 */
const getRoles = async () => {
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
};

/**
 * Fetches a list of namespaces from the API.
 * @returns {Promise<any>} - A promise that resolves to the response data from the API.
 */
const getNamespaces = async () => {
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
};

/**
 * Renders a component that displays a list of roles and namespaces.
 * 
 * The component uses the `useQuery` hook from `@tanstack/react-query` to fetch the list of roles and namespaces from the API.
 * If the data is still being fetched, a skeleton loader is displayed. Otherwise, a `DataTable` component is rendered with the fetched roles and namespaces.
 */
const Roles = () => {
  // Get Roles
  const { data: roles , isPending: isPendingRoles } = useQuery({
    queryKey: ["roles"],
    queryFn: getRoles
  }); 

  // Get Namespaces
  const { data: namespace } = useQuery({
    queryKey: ["namespace"],
    queryFn: getNamespaces
  });
  
  return (
    <div className="flex w-full flex-col">
      {isPendingRoles ? <Skeleton className="h-full w-100 m-4"></Skeleton> : <DataTable roles={roles.items} namespace={namespace}></DataTable>}
    </div>
  );
};

export default Roles;
