'use client';

import { Skeleton } from '@/components/ui/skeleton';
import { useQuery } from '@tanstack/react-query';
import React from 'react'
import DataTable from './_components/DataTable';

/**
 * Fetches a list of namespaces from the API.
 * @returns {Promise<any>} - A promise that resolves to the response data from the API.
 */
const getClusterRoles = async () => {
  const URL = "http://localhost:8080/api/clusterroles?namespaces=all";
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

const ClusterRoles = () => {
 // Get ClusterRoles
 const { data: clusterRoles, isLoading, isError,  } = useQuery({
  queryKey: ["roles"],
  queryFn: getClusterRoles
}); 

if(isError){
  return <div>Error</div>
}

return (
  <div className="flex w-full flex-col">
    {isLoading ? <Skeleton className="h-full w-100 m-4"></Skeleton> : <DataTable clusterRoles={clusterRoles}></DataTable>}
  </div>
);
}

export default ClusterRoles