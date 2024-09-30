"use client";

import React from "react";
import { Button } from "@/components/ui/button";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
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
import { useCreateRoleBinding } from "@/hooks/useRoleBindings";

const formSchema = z.object({
  roleBindingName: z.string().min(1, "Name is required"),
  namespace: z.string().min(1, "Namespace is required"),
  roleName: z.string().min(1, "Role name is required"),
  subjectKind: z.enum(["User", "Group", "ServiceAccount"]),
  subjectName: z.string().min(1, "Subject name is required"),
});

const CreateRoleBindingsDialog = ({
  isOpen,
  setIsOpen,
}: {
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
}) => {
  const createRoleBindingMutation = useCreateRoleBinding();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      roleBindingName: "",
      namespace: "",
      roleName: "",
      subjectKind: "User",
      subjectName: "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    const roleBindingData = {
      metadata: {
        name: values.roleBindingName,
        namespace: values.namespace,
      },
      roleRef: {
        kind: "Role",
        name: values.roleName,
        apiGroup: "rbac.authorization.k8s.io",
      },
      subjects: [
        {
          kind: values.subjectKind,
          name: values.subjectName,
          apiGroup: "rbac.authorization.k8s.io",
        },
      ],
    };
    await createRoleBindingMutation.mutateAsync(roleBindingData);
    form.reset();
  };

  return (
    <ResponsiveDialog
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      title="Create New Role Binding"
    >
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
            control={form.control}
            name="roleBindingName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input placeholder="Enter role binding name" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="namespace"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Namespace</FormLabel>
                <FormControl>
                  <Input placeholder="Enter namespace" {...field} />
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
          <Button type="submit">Create Role Binding</Button>
        </form>
      </Form>
    </ResponsiveDialog>
  );
};

export default CreateRoleBindingsDialog;
