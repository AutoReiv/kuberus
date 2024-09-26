"use client";

import React, { useState } from "react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import GenericDataTable from "@/components/GenericDataTable";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/apiClient";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, Eye, Plus, Trash } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
import { useRouter } from "next/navigation";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface ClusterRoleBinding {
  metadata: {
    name: string;
    uid: string;
    creationTimestamp: string;
  };
  subjects: {
    kind: string;
    name: string;
  }[];
  roleRef: {
    kind: string;
    name: string;
  };
}

const formSchema = z.object({
  name: z.string().min(1, "Name is required"),
  roleName: z.string().min(1, "Role name is required"),
  subjectKind: z.enum(["User", "Group", "ServiceAccount"]),
  subjectName: z.string().min(1, "Subject name is required"),
});

const ClusterRoleBindings = () => {
  const router = useRouter();
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const queryClient = useQueryClient();

  // Get Cluster Role Bindings
  const { data: clusterRoleBindings, isPending: isPendingRoles } = useQuery({
    queryKey: ["clusterRoleBindings"],
    queryFn: () => apiClient.getClusterRoleBindings(),
  });

  const columns: ColumnDef<ClusterRoleBinding>[] = [
    {
      accessorKey: "name",
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
      cell: ({ row }) => <div>{row.original.metadata.name}</div>,
    },
    {
      accessorKey: "roleRef.name",
      header: "Role",
      cell: ({ row }) => <div>{row.original.roleRef.name}</div>,
    },
    {
      accessorKey: "subjects",
      header: "Subjects",
      cell: ({ row }) => (
        <div className="flex flex-wrap gap-1">
          {row.original.subjects.map((subject, index) => (
            <Badge key={index} variant="secondary">
              {subject.kind}: {subject.name}
            </Badge>
          ))}
        </div>
      ),
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
      cell: ({ row }) => (
        <div>
          {new Date(row.original.metadata.creationTimestamp).toLocaleString()}
        </div>
      ),
    },
  ];

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      roleName: "",
      subjectKind: "User",
      subjectName: "",
    },
  });

  const createClusterRoleBindingMutation = useMutation({
    mutationFn: (data: z.infer<typeof formSchema>) =>
      apiClient.createClusterRoleBinding(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["clusterRoleBindings"] });
      setIsCreateDialogOpen(false);
      form.reset();
    },
  });

  const onSubmit = (values: z.infer<typeof formSchema>) => {
    const payload = {
      metadata: {
        name: values.name,
      },
      subjects: [
        {
          kind: values.subjectKind,
          name: values.subjectName,
          apiGroup: "rbac.authorization.k8s.io",
        },
      ],
      roleRef: {
        kind: "ClusterRole",
        name: values.roleName,
        apiGroup: "rbac.authorization.k8s.io",
      },
    };
    createClusterRoleBindingMutation.mutate(payload);
  };

  const deleteClusterRoleBindingMutation = useMutation({
    mutationFn: (name: string) => apiClient.deleteClusterRoleBinding(name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["clusterRoleBindings"] });
    },
  });

  const handleDeleteClusterRoleBinding = (row: ClusterRoleBinding) => {
    if (
      confirm(
        `Are you sure you want to delete the cluster role binding "${row.metadata.name}"?`
      )
    ) {
      deleteClusterRoleBindingMutation.mutate(row.metadata.name);
    }
  };

  const routeToDetails = (name: string) => {
    router.push(`/dashboard/cluster-role-bindings/${name}`);
  };

  return (
    <motion.div
      className="flex w-full flex-col"
      initial="initial"
      animate="animate"
      exit="exit"
      variants={pageVariants}
    >
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight">
          Cluster Role Bindings
        </h2>
        <Button onClick={() => setIsCreateDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Create Cluster Role Binding
        </Button>
        <ResponsiveDialog
          isOpen={isCreateDialogOpen}
          setIsOpen={setIsCreateDialogOpen}
          title="Create New Cluster Role Binding"
        >
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Enter cluster role binding name"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="roleName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Role Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter role name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="subjectKind"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Subject Kind</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select subject kind" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="User">User</SelectItem>
                        <SelectItem value="Group">Group</SelectItem>
                        <SelectItem value="ServiceAccount">
                          Service Account
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="subjectName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Subject Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter subject name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit">Create Cluster Role Binding</Button>
            </form>
          </Form>
        </ResponsiveDialog>
      </div>
      {isPendingRoles ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={clusterRoleBindings}
          columns={columns}
          rowActions={(row) => [
            <Trash
              key="delete"
              size={20}
              onClick={() => handleDeleteClusterRoleBinding(row)}
            >
              Delete
            </Trash>,
            <Eye
              key="view"
              onClick={() => routeToDetails(row.metadata.name)}
              size={20}
            >
              View Details
            </Eye>,
          ]}
        ></GenericDataTable>
      )}
    </motion.div>
  );
};

export default ClusterRoleBindings;
