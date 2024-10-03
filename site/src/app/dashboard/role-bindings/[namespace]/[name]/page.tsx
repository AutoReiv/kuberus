"use client";

import { format } from "date-fns";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/apiClient";
import { RoleBinding } from "@/interfaces/roleBinding";

const RoleBindingDetailsPage = ({
  params,
}: {
  params: { namespace: string; name: string };
}) => {
  const { namespace, name } = params;

  const {
    data: roleBindingDetails,
    isLoading,
    error
  } = useQuery<RoleBinding, Error>({
    queryKey: ["roleDetails", namespace, name],
    queryFn: () => apiClient.getRoleBindingsDetails(namespace, name),
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <div className="flex min-h-screen w-full flex-col bg-muted/40">
      <div className="flex flex-col sm:gap-4">
        <main className="grid flex-1 items-start gap-4 p-4 sm:px-6 md:gap-8 lg:grid-cols-3 xl:grid-cols-3">
          <div>
            <Card className="overflow-hidden">
              <CardContent className="p-6 text-sm">
                <div className="grid gap-3">
                  <div className="font-semibold flex items-center justify-between">
                    Role Binding Details
                  </div>
                  {roleBindingDetails && (
                    <ul className="grid gap-3">
                      <li className="flex items-center justify-between">
                        <span className="text-muted-foreground">Name:</span>
                        <span>{roleBindingDetails.metadata.name}</span>
                      </li>
                      <li className="flex items-center justify-between">
                        <span className="text-muted-foreground">
                          Namespace:
                        </span>
                        <span>{roleBindingDetails.metadata.namespace}</span>
                      </li>
                      <li className="flex items-center justify-between">
                        <span className="text-muted-foreground">
                          Creation Date:
                        </span>
                        <span>
                          {format(
                            new Date(
                              roleBindingDetails.metadata.creationTimestamp
                            ),
                            "MM/dd - hh:mm:ss a"
                          )}
                        </span>
                      </li>
                      <li className="flex items-center justify-between">
                        <span className="text-muted-foreground">
                          Resource Version:
                        </span>
                        <span>
                          {roleBindingDetails.metadata.resourceVersion}
                        </span>
                      </li>
                    </ul>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>
          <div className="grid auto-rows-max items-start gap-4 md:gap-8 lg:col-span-2">
            <Tabs defaultValue="subjects">
              <div className="flex items-center">
                <TabsList>
                  <TabsTrigger value="subjects">Subjects</TabsTrigger>
                  <TabsTrigger value="roleRef">Role Reference</TabsTrigger>
                  <TabsTrigger value="graph">Graph</TabsTrigger>
                </TabsList>
              </div>
              <TabsContent value="subjects">
                <Card>
                  <CardHeader className="px-7 flex-row items-center justify-between">
                    <div>
                      <CardTitle>Subjects</CardTitle>
                      <CardDescription>
                        Subjects bound to this role binding.
                      </CardDescription>
                    </div>
                  </CardHeader>
                  <CardContent>
                    {roleBindingDetails && (
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableCell>Kind</TableCell>
                            <TableCell>Name</TableCell>
                            <TableCell>Namespace</TableCell>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {roleBindingDetails.subjects.map((subject, index) => (
                            <TableRow key={index}>
                              <TableCell>{subject.kind}</TableCell>
                              <TableCell>{subject.name}</TableCell>
                              <TableCell>
                                {subject.namespace || "N/A"}
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    )}
                  </CardContent>
                </Card>
              </TabsContent>
              <TabsContent value="roleRef">
                <Card>
                  <CardHeader className="px-7 flex-row items-center justify-between">
                    <div>
                      <CardTitle>Role Reference</CardTitle>
                      <CardDescription>
                        Reference to the role being bound.
                      </CardDescription>
                    </div>
                  </CardHeader>
                  <CardContent>
                    {roleBindingDetails && (
                      <ul className="grid gap-3">
                        <li className="flex items-center justify-between">
                          <span className="text-muted-foreground">Kind:</span>
                          <span>{roleBindingDetails.roleRef.kind}</span>
                        </li>
                        <li className="flex items-center justify-between">
                          <span className="text-muted-foreground">Name:</span>
                          <span>{roleBindingDetails.roleRef.name}</span>
                        </li>
                        <li className="flex items-center justify-between">
                          <span className="text-muted-foreground">
                            API Group:
                          </span>
                          <span>{roleBindingDetails.roleRef.apiGroup}</span>
                        </li>
                      </ul>
                    )}
                  </CardContent>
                </Card>
              </TabsContent>
              <TabsContent value="graph">

              </TabsContent>
            </Tabs>
          </div>
        </main>
      </div>
    </div>
  );
};

export default RoleBindingDetailsPage;
