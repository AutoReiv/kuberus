"use client";
import React, { useCallback } from "react";
import ForceGraph3D from "react-force-graph-3d";
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
import SpriteText from "three-spritetext";
import { useQuery } from "@tanstack/react-query";
import { useTheme } from "next-themes";
import { apiClient } from "@/lib/apiClient";

interface RoleBindingDetail {
  metadata: {
    name: string;
    namespace: string;
    creationTimestamp: string;
    resourceVersion: string;
  };
  subjects: {
    kind: string;
    name: string;
    namespace?: string;
  }[];
  roleRef: {
    kind: string;
    name: string;
    apiGroup: string;
  };
}

// const fetchRoleBindingDetails = async (namespace, name) => {
//   const URL = `http://localhost:8080/api/rolebinding/details?name=${name}&namespace=${namespace}`;
//   const response = await fetch(URL, {
//     method: "GET",
//     headers: {
//       Accept: "application/json",
//       "Content-Type": "application/json",
//     },
//   });

//   const data = await response.json();
//   return data;
// };

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
  } = useQuery<RoleBindingDetail, Error>({
    queryKey: ["roleDetails", namespace, name],
    queryFn: () => apiClient.getRoleBindingDetails(namespace, name),
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  interface RoleBindingForceGraphProps {
    roleBinding: RoleBindingDetail;
  }

  const RoleBindingForceGraph: React.FC<RoleBindingForceGraphProps> = ({
    roleBinding,
  }) => {
    const { theme } = useTheme();

    const getNodeColor = useCallback(
      (node) => {
        const colors = {
          light: {
            1: "#0ea5e9", // sky-500
            2: "#22c55e", // green-500
            3: "#f59e0b", // amber-500
          },
          dark: {
            1: "#0284c7", // sky-600
            2: "#16a34a", // green-600
            3: "#d97706", // amber-600
          },
        };
        return colors[theme][node.group];
      },
      [theme]
    );

    const graphData = {
      nodes: [
        {
          id: "roleBinding",
          group: 1,
          label: `RoleBinding: ${roleBinding.metadata.name}`,
        },
        {
          id: "role",
          group: 2,
          label: `${roleBinding.roleRef.kind}: ${roleBinding.roleRef.name}`,
        },
        ...roleBinding.subjects.map((subject, index) => ({
          id: `subject-${index}`,
          group: 3,
          label: `${subject.kind}: ${subject.name}`,
        })),
      ],
      links: [
        { source: "roleBinding", target: "role" },
        ...roleBinding.subjects.map((_, index) => ({
          source: "roleBinding",
          target: `subject-${index}`,
        })),
      ],
    };

    return (
      <ForceGraph3D
        width={800}
        height={600}
        graphData={graphData}
        nodeLabel="label"
        nodeThreeObjectExtend={true}
        nodeThreeObject={(node) => {
          const sprite = new SpriteText(node.label);
          sprite.color = getNodeColor(node);
          sprite.textHeight = 8;
          return sprite;
        }}
        nodeAutoColorBy={getNodeColor}
        linkColor={() => theme === "dark" ? "var(--secondary)" : "var(--muted)"} // gray-600 for dark, gray-400 for light
        backgroundColor={theme === "dark" ? "#1f2937" : "#f3f4f6"} // gray-800 for dark, gray-100 for light
        linkDirectionalParticles={7}
        linkDirectionalParticleSpeed={0.001}
        d3VelocityDecay={0.5}
      />
    );
  };

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
                {/* <RoleBindingForceGraph roleBinding={roleBindingDetails} /> */}

              </TabsContent>
            </Tabs>
          </div>
        </main>
      </div>
    </div>
  );
};

export default RoleBindingDetailsPage;
