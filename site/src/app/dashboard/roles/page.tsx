"use client";

import React from "react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import { Checkbox } from "@/components/ui/checkbox";
import { ColumnDef } from "@tanstack/react-table";
import { Role } from "../_interfaces/role";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, MoreHorizontal } from "lucide-react";
import Link from "next/link";
import { ResponsiveDialog } from "./_components/ResponsiveDialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { format } from "date-fns";
import { SkeletonPage } from "@/components/SkeletonPage";

/**
 * Renders a component that displays a list of roles and namespaces.
 *
 * The component uses the `useQuery` hook from `@tanstack/react-query` to fetch the list of roles and namespaces from the API.
 * If the data is still being fetched, a skeleton loader is displayed. Otherwise, a `DataTable` component is rendered with the fetched roles and namespaces.
 */
const Roles = () => {
  // Get Roles
  const {
    data: roles,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["roles"],
    queryFn: () => apiClient.getRoles(),
  });

  // Get Namespaces
  const { data: namespace } = useQuery({
    queryKey: ["namespace"],
    queryFn: () => apiClient.getNamespaces(),
  });

  if (isError) {
    return <div>Error</div>;
  }

  const columns: ColumnDef<Role>[] = [
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={table.getIsAllPageRowsSelected()}
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      id: "name",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Name
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
      accessorKey: "metadata.name",
      cell: ({ row }) => {
        const name = row.getValue("name") as string;
        const namespace = row.original.metadata.namespace;
        return (
          <Link
            href={`/dashboard/roles/${namespace}/${name}`}
            className="hover:underline"
          >
            {name}
          </Link>
        );
      },
    },
    {
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Namespace
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
      accessorKey: "metadata.namespace",
    },
    {
      id: "createdAt",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Created At
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
      accessorKey: "metadata.creationTimestamp",
      cell: ({ getValue }) => {
        const timestamp: any = getValue();
        return format(new Date(timestamp), "MM/dd - hh:mm:ss a");
      },
    },
    {
      id: "actions",
      header: "",
      accessorKey: "metadata",
      cell: ({ row }) => {
        const { namespace, name } = row.original.metadata;
        const dialogKey = `${namespace}-${name}`;
        // return (
        //   <>
        //     <ResponsiveDialog
        //       isOpen={deleteConfirmationDialogs[dialogKey] || false}
        //       setIsOpen={(isOpen) =>
        //         setDeleteConfirmationDialog(dialogKey, isOpen === true)
        //       }
        //       title="Delete Role"
        //       description="Are you sure you want to delete this role?"
        //     >
        //       <Button
        //         variant="destructive"
        //         onClick={() => {
        //           deleteRole(namespace, name);
        //           setDeleteConfirmationDialog(dialogKey, false);
        //         }}
        //       >
        //         Confirm
        //       </Button>
        //       <Button
        //         variant="ghost"
        //         onClick={() => setDeleteConfirmationDialog(dialogKey, false)}
        //       >
        //         Cancel
        //       </Button>
        //     </ResponsiveDialog>
        //     <DropdownMenu>
        //       <DropdownMenuTrigger asChild>
        //         <Button aria-haspopup="true" size="icon" variant="ghost">
        //           <MoreHorizontal className="h-4 w-4" />
        //           <span className="sr-only">Toggle menu</span>
        //         </Button>
        //       </DropdownMenuTrigger>
        //       <DropdownMenuContent align="end">
        //         <DropdownMenuLabel>Actions</DropdownMenuLabel>
        //         <DropdownMenuItem>Edit</DropdownMenuItem>
        //         <DropdownMenuItem>
        //           <button
        //             onClick={() => setDeleteConfirmationDialog(dialogKey, true)}
        //           >
        //             Delete
        //           </button>
        //         </DropdownMenuItem>
        //       </DropdownMenuContent>
        //     </DropdownMenu>
        //   </>
        // );
      },
    },
  ];

  return (
    <div className="flex w-full flex-col">
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={roles}
          columns={columns}
          title="Roles"
          description="Manage and configure user roles to control access and permissions across your Kubernetes cluster"
          // Add in row action to route to details
        ></GenericDataTable>
      )}
    </div>
  );
};

export default Roles;
