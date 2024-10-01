"use client";

import React, { useState } from "react";
import GenericDataTable from "@/components/GenericDataTable";
import { Button } from "@/components/ui/button";
import { Eye, Plus, Trash } from "lucide-react";
import { ColumnDef } from "@tanstack/react-table";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useRouter } from "next/navigation";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  useCreateNamespace,
  useDeleteNamespaces,
  useNamespaces,
} from "@/hooks/useNamespace";
import { Namespace } from "@/interfaces/namespace";

const namespaceFormSchema = z.object({
  namespace: z.string().min(1, "Namespace is required"),
});
type NamespaceFormData = z.infer<typeof namespaceFormSchema>;

const Namespaces = () => {
  const router = useRouter();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const { data: namespaces, isLoading, error } = useNamespaces();
  const createNamespace = useCreateNamespace({
    onSuccess: () => {
      setIsDialogOpen(false);
    },
  });
  const deleteNamespace = useDeleteNamespaces();
  const columns: ColumnDef<Namespace>[] = [
    {
      accessorKey: "metadata.name",
      header: "Name",
    },
    {
      accessorKey: "metadata.uid",
      header: "UID",
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
    },
  ];

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<NamespaceFormData>({
    resolver: zodResolver(namespaceFormSchema),
  });

  const onSubmit = (data: NamespaceFormData) => {
    createNamespace.mutate(data.namespace);
  };

  const onDelete = (row: Namespace) => {
    deleteNamespace.mutate(row);
  };

  const routeToDetails = (row: Namespace) => {
    router.push(`/dashboard/namespaces/${row.metadata.name}`);
  };

  if (isLoading) return <SkeletonPage />;
  if (error) return <div>Error loading namespaces</div>;

  return (
    <motion.div
      variants={pageVariants}
      initial="hidden"
      animate="visible"
      exit="exit"
      className="w-full"
    >
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Namespaces</h1>
        <Button onClick={() => setIsDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" /> Create Namespace
        </Button>
        <ResponsiveDialog
          isOpen={isDialogOpen}
          setIsOpen={setIsDialogOpen}
          title="Create RBAC Rule"
        >
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <Label htmlFor="namespace">Namespace</Label>
              <Input
                id="namespace"
                {...register("namespace")}
                placeholder="Enter namespace"
              />
              {errors.namespace && (
                <p className="text-red-500">{errors.namespace.message}</p>
              )}
            </div>
            <Button type="submit">Create</Button>
          </form>
        </ResponsiveDialog>
      </div>
      <GenericDataTable
        columns={columns}
        data={namespaces.items}
        enableGridView={false}
        rowActions={(row) => [
          <Trash key="delete" size={20} onClick={() => onDelete(row)}>
            Delete
          </Trash>,
          <Eye key="view" size={20} onClick={() => routeToDetails(row)}>
            View Details
          </Eye>,
        ]}
      />
    </motion.div>
  );
};
export default Namespaces;
