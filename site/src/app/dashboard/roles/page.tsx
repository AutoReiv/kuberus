"use client";

import React from "react";
import DataTable from "./_components/DataTable";
import { useQuery } from "@tanstack/react-query";

const Roles = () => {
  const { data, isPending, isError } = useQuery({
    queryKey: ["roles"],
    queryFn: async () => {
      const URL = "http://localhost:8080/api/roles";
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
      {isPending ? <div>Loading...</div> : <DataTable roles={data.items}></DataTable>}
    </div>
  );
};

export default Roles;
