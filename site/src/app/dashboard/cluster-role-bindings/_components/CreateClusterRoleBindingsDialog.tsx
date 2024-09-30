"use client";

import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
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
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateClusterRoleBinding } from "@/hooks/useClusterRoleBinding";

const formSchema = z.object({
  name: z.string().min(1, "Name is required"),
  roleName: z.string().min(1, "Role name is required"),
  subjectKind: z.enum(["User", "Group", "ServiceAccount"]),
  subjectName: z.string().min(1, "Subject name is required"),
});

const CreateClusterRoleBindingsDialog = ({
  isOpen,
  setIsOpen,
}: {
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
}) => {
  const createClusterRoleBindingMutation = useCreateClusterRoleBinding({
    onSuccess: () => {
      setIsOpen(false);
    },
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      roleName: "",
      subjectKind: "User",
      subjectName: "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
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
    await createClusterRoleBindingMutation.mutate(payload);
  };
  return (
    <ResponsiveDialog
      isOpen={isOpen}
      setIsOpen={setIsOpen}
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
  );
};

export default CreateClusterRoleBindingsDialog;
