"use client";

import React from "react";
import GenericDataTable from "@/components/GenericDataTable";
import { ArrowUpDown, Eye, Plus, Trash } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { ColumnDef } from "@tanstack/react-table";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useDeleteRoleBinding, useRoleBindings } from "@/hooks/useRoleBindings";
import { RoleBinding } from "@/interfaces/roleBinding";
import CreateRoleBindingsDialog from "./_components/CreateRoleBindingsDialog";

const RoleBindings = () => {
  const router = useRouter();
  const [isCreateRoleBindingDialogOpen, setIsCreateRoleBindingDialogOpen] =
    React.useState(false);
  const { data: roleBindings, isLoading } = useRoleBindings();
  const deleteRoleBindingMutation = useDeleteRoleBinding();

  const columns: ColumnDef<RoleBinding>[] = [
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
      accessorKey: "metadata.name",
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
    },
    {
      accessorKey: "metadata.namespace",
      header: "Namespace",
    },
    {
      accessorKey: "roleRef.name",
      header: "Role",
    },
    {
      accessorKey: "subjects",
      header: "Subjects",
      cell: ({ row }) => {
        const subjects = row.original.subjects;
        return (
          <div>
            {subjects.map((subject, index) => (
              <div key={index}>{`${subject.kind}: ${subject.name}`}</div>
            ))}
          </div>
        );
      },
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
      cell: ({ row }) => {
        return new Date(
          row.original.metadata.creationTimestamp
        ).toLocaleString();
      },
    },
  ];

  const routeToDetails = (roleBinding: RoleBinding) => {
    router.push(
      `/dashboard/role-bindings/${roleBinding.metadata.namespace}/${roleBinding.metadata.name}`
    );
    toast(
      `Routing to details page for ${roleBinding.metadata.name} in ${roleBinding.metadata.namespace}`
    );
  };

  const onDelete = async (roleBinding: RoleBinding) => {
    await deleteRoleBindingMutation.mutateAsync({
      namespace: roleBinding.metadata.namespace,
      name: roleBinding.metadata.name,
    });
  };

  return (
    <motion.div
      className="flex w-full flex-col"
      initial="initial"
      animate="animate"
      exit="exit"
      variants={pageVariants}
    >
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-3xl font-bold tracking-tight">Role Bindings</h2>
        <Button onClick={() => setIsCreateRoleBindingDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Create Role Binding
        </Button>
        <CreateRoleBindingsDialog
          isOpen={isCreateRoleBindingDialogOpen}
          setIsOpen={setIsCreateRoleBindingDialogOpen}
        />
      </div>
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={roleBindings}
          columns={columns}
          enableGridView={false}
          rowActions={(row) => [
            <Trash key="delete" size={20} onClick={() => onDelete(row)}>
              Delete
            </Trash>,
            <Eye key="view" size={20} onClick={() => routeToDetails(row)}>
              View Details
            </Eye>,
          ]}
        ></GenericDataTable>
      )}
    </motion.div>
  );
};
export default RoleBindings;
