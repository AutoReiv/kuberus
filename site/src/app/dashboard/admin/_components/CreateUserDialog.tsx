"use client";

import { useQuery, useQueryClient, useMutation } from "@tanstack/react-query";
import { ArrowUpDown, Package2, Plus, Trash } from "lucide-react";
import { Button } from "@/components/ui/button";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
import { useState } from "react";
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
import { toast } from "sonner";
import { motion } from "framer-motion";


const createUserSchema = z
  .object({
    username: z.string().min(3, "Username must be at least 3 characters"),
    password: z.string().min(8, "Password must be at least 8 characters"),
    passwordConfirm: z.string(),
  })
  .refine((data) => data.password === data.passwordConfirm, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

const CreateUserDialog = ({ isOpen, setIsOpen, onCreateUser }) => {
  const queryClient = useQueryClient();

  const form = useForm({
    resolver: zodResolver(createUserSchema),
    defaultValues: {
      username: "",
      password: "",
      passwordConfirm: "",
    },
  });

  const createUserMutation = useMutation({
    mutationFn: (userData: {
      username: string;
      password: string;
      passwordConfirm: string;
    }) => apiClient.getAdminCreateUser(userData),
    onMutate: () => {
      toast.loading("Creating user...");
    },
    onSuccess: (_, variables) => {
      toast.dismiss();
      toast.success(
        `User ${variables.username} has been created successfully.`
      );
      queryClient.invalidateQueries({ queryKey: ["users"] });
      setIsOpen(false);
      form.reset();
    },
    onError: (error) => {
      toast.dismiss();
      toast.error(
        error.message || "An unexpected error occurred while creating the user."
      );
    },
  });

  const onSubmit = (data) => {
    createUserMutation.mutate({
      username: data.username,
      password: data.password,
      passwordConfirm: data.passwordConfirm,
    });
  };

  return (
    <ResponsiveDialog
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      title="Create New User"
    >
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter username" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input
                      type="password"
                      placeholder="Enter password"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="passwordConfirm"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Confirm Password</FormLabel>
                  <FormControl>
                    <Input
                      type="password"
                      placeholder="Confirm password"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit">Create User</Button>
          </form>
        </Form>
      </motion.div>
    </ResponsiveDialog>
  );
};


export default CreateUserDialog;