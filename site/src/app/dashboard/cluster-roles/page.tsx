"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import React from "react";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, Eye, Plus, Trash } from "lucide-react";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useFieldArray, useForm } from "react-hook-form";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";

interface ClusterRole {
  metadata: {
    name: string;
    creationTimestamp: string;
  };
  rules: {
    apiGroups: string[];
    resources: string[];
    verbs: string[];
  }[];
}

const clusterRoleSchema = z.object({
  name: z.string().min(1, "Name is required"),
  rules: z
    .array(
      z.object({
        apiGroups: z.array(z.string()),
        resources: z.array(z.string()),
        verbs: z.array(z.string()),
      })
    )
    .min(1, "At least one rule is required"),
});

type ClusterRoleFormValues = z.infer<typeof clusterRoleSchema>;

const ClusterRoles = () => {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [isCreateDialogOpen, setIsCreateDialogOpen] = React.useState(false);
  // Get ClusterRoles
  const {
    data: clusterRoles,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["roles"],
    queryFn: () => apiClient.getClusterRoles(),
  });

  const columns: ColumnDef<ClusterRole>[] = [
    {
      id: "metadata.name",
      accessorKey: "metadata.name",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Name
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
    },
    {
      accessorKey: "rules",
      header: "Rules",
      cell: ({ row }) => {
        const rules = row.original.rules;
        return <span>{rules.length} rule(s)</span>;
      },
    },
  ];

  const deleteClusterRoleMutation = useMutation({
    mutationFn: (selectedRows: ClusterRole[]) =>
      Promise.all(
        selectedRows.map((row) =>
          apiClient.deleteClusterRole(row.metadata.name)
        )
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
  });

  const handleDeleteRole = (selectedRows: ClusterRole[]) => {
    if (
      confirm("Are you sure you want to delete the selected cluster role(s)?")
    ) {
      deleteClusterRoleMutation.mutate(selectedRows);
    }
  };

  const routeToDetails = (name: string) => {
    router.push(`/dashboard/cluster-roles/${name}`);
  };

  const CreateClusterRoleDialog = ({
    isOpen,
    setIsOpen,
  }: {
    isOpen: boolean;
    setIsOpen: (open: boolean) => void;
  }) => {
    const form = useForm<ClusterRoleFormValues>({
      resolver: zodResolver(clusterRoleSchema),
      defaultValues: {
        name: "",
        rules: [{ apiGroups: [], resources: [], verbs: [] }],
      },
    });

    const { fields, append, remove } = useFieldArray({
      control: form.control,
      name: "rules",
    });

    // @TODO: Look into why this is working but not showing in the table.
    const onSubmit = async (data: ClusterRoleFormValues) => {
      const payload = {
        apiVersion: "rbac.authorization.k8s.io/v1",
        kind: "ClusterRole",
        metadata: {
          name: data.name,
        },
        rules: data.rules.map((rule) => ({
          apiGroups: rule.apiGroups,
          resources: rule.resources,
          verbs: rule.verbs,
        })),
      };

      try {
        await apiClient.createClusterRole(payload);
        toast.success(`Cluster role ${data.name} created successfully`);
        queryClient.invalidateQueries({ queryKey: ["clusterRoles"] });
        setIsOpen(false);
      } catch (error) {
        toast.error(`Failed to create cluster role: ${error.message}`);
      }
    };
    return (
      <ResponsiveDialog
        isOpen={isOpen}
        setIsOpen={setIsOpen}
        title="Create New Cluster Role"
        className="!max-w-none w-[60%] h-[60]"
      >
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <div className="flex items-center gap-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem className="flex-1">
                    <FormLabel>Cluster Role Name *</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter cluster role name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {fields.map((field, index) => (
              <div key={field.id} className="space-y-4 p-4 border rounded-md">
                <div className="flex justify-between items-center">
                  <h4 className="text-lg font-semibold">Rule {index + 1}</h4>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => remove(index)}
                  >
                    <Trash className="h-4 w-4" />
                  </Button>
                </div>

                <FormField
                  control={form.control}
                  name={`rules.${index}.apiGroups`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>API Groups *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Enter API groups (comma-separated)"
                          {...field}
                          onChange={(e) =>
                            field.onChange(e.target.value.split(","))
                          }
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name={`rules.${index}.resources`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Resources *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Enter resources (comma-separated)"
                          {...field}
                          onChange={(e) =>
                            field.onChange(e.target.value.split(","))
                          }
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name={`rules.${index}.verbs`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Verbs *</FormLabel>
                      <FormControl>
                        <Select
                          onValueChange={(value) =>
                            field.onChange([...field.value, value])
                          }
                          value=""
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Select verbs" />
                          </SelectTrigger>
                          <SelectContent>
                            {[
                              "get",
                              "list",
                              "watch",
                              "create",
                              "update",
                              "patch",
                              "delete",
                            ].map((verb) => (
                              <SelectItem key={verb} value={verb}>
                                {verb}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </FormControl>
                      <div className="flex flex-wrap gap-2 mt-2">
                        {field.value.map((verb, verbIndex) => (
                          <div
                            key={verbIndex}
                            className="bg-secondary text-secondary-foreground px-2 py-1 rounded-md flex items-center gap-2"
                          >
                            {verb}
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              onClick={() =>
                                field.onChange(
                                  field.value.filter((_, i) => i !== verbIndex)
                                )
                              }
                            >
                              <Trash className="h-3 w-3" />
                            </Button>
                          </div>
                        ))}
                      </div>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            ))}

            <Button
              type="button"
              variant="outline"
              onClick={() =>
                append({ apiGroups: [], resources: [], verbs: [] })
              }
            >
              Add Rule
            </Button>

            <Button type="submit">Create Cluster Role</Button>
          </form>
        </Form>
      </ResponsiveDialog>
    );
  };

  if (isError) {
    return <div>Error</div>;
  }

  return (
    <motion.div
      className="flex w-full flex-col"
      initial="initial"
      animate="animate"
      exit="exit"
      variants={pageVariants}
    >
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Cluster Roles</h1>
        <Button onClick={() => setIsCreateDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" /> Create Cluster Role
        </Button>
        <CreateClusterRoleDialog
          isOpen={isCreateDialogOpen}
          setIsOpen={setIsCreateDialogOpen}
        />
      </div>
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={clusterRoles}
          columns={columns}
          enableGridView={false}
          rowActions={(row) => [
            <Trash
              key="delete"
              size={20}
              onClick={() => handleDeleteRole([row])}
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
          enableRowSelection={true}
          bulkActions={(selectedRows) => [
            <Button
              key="delete"
              onClick={() => handleDeleteRole(selectedRows)}
              variant="destructive"
            >
              Delete Selected Cluster Roles
            </Button>,
          ]}
        ></GenericDataTable>
      )}
    </motion.div>
  );
};
export default ClusterRoles;
