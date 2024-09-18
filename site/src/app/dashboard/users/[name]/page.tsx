"use client";

import React from "react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/apiClient";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Separator } from "@/components/ui/separator";
import { pageVariants } from "../../layout";
import GenericDataTable from "@/components/GenericDataTable";
import { SkeletonPage } from "@/components/SkeletonPage";
import { format } from "date-fns";

interface RoleBinding {
  metadata: {
    name: string;
    namespace: string;
    uid: string;
    resourceVersion: string;
    creationTimestamp: string;
  };
  subjects: {
    kind: string;
    apiGroup: string;
    name: string;
  }[];
  roleRef: {
    apiGroup: string;
    kind: string;
    name: string;
  };
}

interface UserDetails {
  userName: string;
  roleBindings: RoleBinding[] | null;
  clusterRoleBindings: any[] | null;
  clusterRoles: any[] | null;
}

const UserDetailsPage = ({ params }: { params: { name: string } }) => {
  const { name } = params;

  const {
    data: userDetails,
    isLoading,
    error,
  } = useQuery<UserDetails, Error>({
    queryKey: ["userDetails", name],
    queryFn: () => apiClient.getUserDetails(name),
  });

  if (isLoading) return <SkeletonPage />;
  if (error) return <div>Error loading user details: {error.message}</div>;
  if (!userDetails) return <div>No user details found</div>;

  const roleBindingsColumns = [
    { accessorKey: "metadata.name", header: "Role Binding Name" },
    { accessorKey: "metadata.namespace", header: "Namespace" },
    { accessorKey: "roleRef.name", header: "Role Name" },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Creation Time",
      cell: ({ getValue }: any) =>
        format(new Date(getValue()), "yyyy-MM-dd HH:mm:ss"),
    },
  ];
  const clusterRoleBindingsColumns = [
    { accessorKey: "name", header: "Cluster Role Binding Name" },
  ];

  const clusterRolesColumns = [
    { accessorKey: "name", header: "Cluster Role Name" },
  ];

  return (
    <motion.div
      initial="initial"
      animate="animate"
      exit="exit"
      variants={pageVariants}
      className="w-full h-full"
    >
      <div className="grid gap-6 md:grid-cols-1 lg:grid-cols-3">
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle>User Details</CardTitle>
            <CardDescription>
              Information about {userDetails.userName}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ul className="grid gap-3">
              <li className="flex items-center justify-between">
                <span className="text-muted-foreground">Username:</span>
                <span>{userDetails.userName}</span>
              </li>
            </ul>
          </CardContent>
        </Card>

        <div className="lg:col-span-2">
          <Tabs defaultValue="roleBindings">
            <TabsList>
              <TabsTrigger value="roleBindings">Role Bindings</TabsTrigger>
              <TabsTrigger value="clusterRoleBindings">
                Cluster Role Bindings
              </TabsTrigger>
              <TabsTrigger value="clusterRoles">Cluster Roles</TabsTrigger>
            </TabsList>
            <TabsContent value="roleBindings">
              <Card>
                <CardHeader>
                  <CardTitle>Role Bindings</CardTitle>
                  <CardDescription>
                    Role bindings associated with this user
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <GenericDataTable
                    data={userDetails.roleBindings || []}
                    columns={roleBindingsColumns}
                    enableSorting
                    enableFiltering
                    enablePagination
                  />
                </CardContent>
              </Card>
            </TabsContent>
            <TabsContent value="clusterRoleBindings">
              <Card>
                <CardHeader>
                  <CardTitle>Cluster Role Bindings</CardTitle>
                  <CardDescription>
                    Cluster role bindings associated with this user
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <GenericDataTable
                    data={userDetails.clusterRoleBindings || []}
                    columns={clusterRoleBindingsColumns}
                    enableSorting
                    enableFiltering
                    enablePagination
                  />
                </CardContent>
              </Card>
            </TabsContent>
            <TabsContent value="clusterRoles">
              <Card>
                <CardHeader>
                  <CardTitle>Cluster Roles</CardTitle>
                  <CardDescription>
                    Cluster roles associated with this user
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <GenericDataTable
                    data={userDetails.clusterRoles || []}
                    columns={clusterRolesColumns}
                    enableSorting
                    enableFiltering
                    enablePagination
                  />
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </motion.div>
  );
};

export default UserDetailsPage;
