"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Copy } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { format } from "date-fns";
import { useEffect, useState } from "react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Plus } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import yaml from "js-yaml";
interface Resources {
  [key: string]: string[];
}

interface Rules {
  apiGroups: string[];
  resources: string[];
  resourceNames: string[];
  verbs: string[];
}

interface RoleDetails {
  role: {
    metadata: {
      creationTimestamp: string;
      managedFields: [
        {
          apiVersion: string;
          fieldsType: string;
          fieldsV1: {
            f: {
              rules: [
                {
                  apiGroups: string[];
                  resources: string[];
                  resourceNames: string[];
                  verbs: string[];
                }
              ];
            };
          };
          manager: string;
          operation: string;
          time: string;
        }
      ];
      name: string;
      namespace: string;
      resourceVersion: string;
      uid: string;
    };
    rules: [
      {
        apiGroups: string[];
        resources: string[];
        resourceNames: string[];
        verbs: string[];
      }
    ];
  };
}

const newRuleSchema = z.object({
  resources: z.string().min(1, "You must select a resource"),
  verbs: z.array(z.string()).refine((value) => value.some((item) => item), {
    message: "You have to select at least one item.",
  }),
});

type NewRuleFormValues = z.infer<typeof newRuleSchema>;

const verbs = [
  { name: "Create" },
  { name: "Delete" },
  { name: "Get" },
  { name: "List" },
  { name: "Patch" },
  { name: "Update" },
  { name: "Watch" },
];

const fetchRoleDetails = async (namespace: string, name: string) => {
  const URL = `http://localhost:8080/api/roles/details?roleName=${name}&namespace=${namespace}`;
  const response = await fetch(URL, {
    method: "GET",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  });

  const data = await response.json();
  console.log(data, " DATA ****");
  return data;
};

const fetchResources = async () => {
  const URL = `http://localhost:8080/api/resources`;
  const response = await fetch(URL, {
    method: "GET",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  });

  const data = await response.json();
  return data;
};

const RoleDetailsPage = ({
  params,
}: {
  params: { namespace: string; name: string };
}) => {
  const { namespace, name } = params;
  const queryClient = useQueryClient();
  const {
    data: roleDetails,
    isLoading,
    error,
    refetch: refetchRoleDetails,
  } = useQuery<RoleDetails, Error>({
    queryKey: ["roleDetails", namespace, name],
    queryFn: () => fetchRoleDetails(namespace, name),
  });

  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [newRules, setNewRules] = useState<Rules[]>([]);

  const { data: resources } = useQuery<Resources, Error>({
    queryKey: ["resources"],
    queryFn: () => fetchResources(),
  });

  const form = useForm<NewRuleFormValues>({
    resolver: zodResolver(newRuleSchema),
    defaultValues: {
      resources: "",
      verbs: ["Create"],
    },
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  const onSubmit = (data: NewRuleFormValues) => {
    const newRule: Rules = {
      apiGroups: [""],
      resources: [data.resources],
      resourceNames: [],
      verbs: data.verbs,
    };
    setNewRules([...newRules, newRule]);
    setIsDialogOpen(false);
    form.reset();
  };

  const confirmRules = async () => {
    const updatedRules = [...roleDetails.role.rules, ...newRules];
    const roleData = {
      metadata: {
        name: name,
        namespace: namespace,
      },
      rules: updatedRules,
    };

    const URL = `http://localhost:8080/api/roles?namespace=${namespace}&name=${name}`;
    await fetch(URL, {
      method: "PUT",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(roleData),
    });

    queryClient.setQueryData(
      ["roleDetails", namespace, name],
      (oldData: any) => ({
        ...oldData,
        role: {
          ...oldData.role,
          rules: updatedRules,
        },
      })
    );
    // Clear new rules
    setNewRules([]);
    // Refetch role details
    refetchRoleDetails();
  };

  const exportRoleDetails = (roleDetails: RoleDetails) => {
    const yamlData = yaml.dump(roleDetails.role);
    const blob = new Blob([yamlData], { type: "text/yaml" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `${roleDetails.role.metadata.name}-role.yaml`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const PermissionImpactAnalysis = ({ rules }: { rules: Rules[] }) => {
    const analyzeImpact = (rules: Rules[]) => {
      let impact = {
        highRisk: 0,
        mediumRisk: 0,
        lowRisk: 0,
      };

      rules.forEach((rule) => {
        if (rule.verbs.includes("Delete") || rule.verbs.includes("Update")) {
          impact.highRisk++;
        } else if (rule.verbs.includes("Create")) {
          impact.mediumRisk++;
        } else {
          impact.lowRisk++;
        }
      });

      return impact;
    };

    const impact = analyzeImpact(rules);

    return (
      <Card>
        <CardHeader>
          <CardTitle>Permission Impact Analysis</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <div className="flex justify-between">
              <span>High Risk Permissions:</span>
              <Badge variant="destructive">{impact.highRisk}</Badge>
            </div>
            <div className="flex justify-between">
              <span>Medium Risk Permissions:</span>
              <Badge variant="destructive">{impact.mediumRisk}</Badge>
            </div>
            <div className="flex justify-between">
              <span>Low Risk Permissions:</span>
              <Badge variant="secondary">{impact.lowRisk}</Badge>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  };

  return (
    <div className="flex min-h-screen w-full flex-col bg-muted/40">
      <div className="flex flex-col sm:gap-4">
        <main className="grid flex-1 items-start gap-4 p-4 sm:px-6 md:gap-8 lg:grid-cols-3 xl:grid-cols-3">
          <div>
            <Card className="overflow-hidden" x-chunk="dashboard-05-chunk-4">
              <CardHeader className="flex flex-row items-start bg-muted/50">
                <div className="grid gap-0.5">
                  <CardTitle className="group flex items-center gap-2 text-lg">
                    {roleDetails.role.metadata.name}
                    <Button
                      size="icon"
                      variant="outline"
                      className="h-6 w-6 opacity-0 transition-opacity group-hover:opacity-100"
                    >
                      <Copy className="h-3 w-3" />
                      <span className="sr-only">Copy Order ID</span>
                    </Button>
                  </CardTitle>
                  <CardDescription>
                    UID: {roleDetails.role.metadata.uid}
                  </CardDescription>
                </div>
              </CardHeader>
              <CardContent className="p-6 text-sm">
                <div className="grid gap-3">
                  <div className="font-semibold">Role Details</div>
                  <ul className="grid gap-3">
                    <li className="flex items-center justify-between">
                      <span className="text-muted-foreground">Name:</span>
                      <span>{roleDetails.role.metadata.name}</span>
                    </li>
                    <li className="flex items-center justify-between">
                      <span className="text-muted-foreground">Namespace:</span>
                      <span>{roleDetails.role.metadata.namespace}</span>
                    </li>
                    <li className="flex items-center justify-between">
                      <span className="text-muted-foreground">
                        Creation Date:
                      </span>
                      <span>
                        {format(
                          new Date(roleDetails.role.metadata.creationTimestamp),
                          "MM/dd - hh:mm:ss a"
                        )}
                      </span>
                    </li>
                    <li className="flex items-center justify-between">
                      <span className="text-muted-foreground">
                        Resource Version:
                      </span>
                      <span>{roleDetails.role.metadata.resourceVersion}</span>
                    </li>
                  </ul>
                  <Separator></Separator>
                  <PermissionImpactAnalysis rules={roleDetails.role.rules} />
                  <Separator></Separator>
                  <Button
                    onClick={() => exportRoleDetails(roleDetails)}
                    variant="outline"
                    size="sm"
                  >
                    Export YAML
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
          <div className="grid auto-rows-max items-start gap-4 md:gap-8 lg:col-span-2">
            <Tabs defaultValue="rules">
              <div className="flex items-center">
                <TabsList>
                  <TabsTrigger value="rules">Rules</TabsTrigger>
                  <TabsTrigger value="roleBindings">Role Bindings</TabsTrigger>
                </TabsList>
              </div>
              <TabsContent value="rules">
                <Card>
                  <CardHeader className="px-7 flex-row items-center justify-between">
                    <div>
                      <CardTitle>Rules</CardTitle>
                      <CardDescription>
                        Here are the rules for this role.
                      </CardDescription>
                    </div>
                    <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                      <DialogTrigger asChild>
                        <Button variant="outline" size="sm">
                          <Plus className="h-4 w-4 mr-2" />
                          Add Rule
                        </Button>
                      </DialogTrigger>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>Add New Rule</DialogTitle>
                        </DialogHeader>
                        <Form {...form}>
                          <form
                            onSubmit={form.handleSubmit(onSubmit)}
                            className="space-y-8"
                          >
                            <FormField
                              control={form.control}
                              name="resources"
                              render={({ field }) => (
                                <FormItem>
                                  <FormLabel>Resources</FormLabel>
                                  <Select
                                    onValueChange={field.onChange}
                                    defaultValue={field.value}
                                  >
                                    <FormControl>
                                      <SelectTrigger>
                                        <SelectValue placeholder="Select a resource" />
                                      </SelectTrigger>
                                    </FormControl>
                                    <SelectContent>
                                      {resources.resources.map((resource) => (
                                        <SelectItem
                                          key={resource}
                                          value={resource}
                                        >
                                          {resource}
                                        </SelectItem>
                                      ))}
                                    </SelectContent>
                                  </Select>
                                  <FormDescription>
                                    Choose the resource for this rule.
                                  </FormDescription>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />
                            <FormField
                              control={form.control}
                              name="verbs"
                              render={({ field }) => (
                                <FormItem className="columns-2">
                                  {verbs.map((badge) => (
                                    <FormField
                                      key={badge.name}
                                      control={form.control}
                                      name="verbs"
                                      render={({ field }) => {
                                        return (
                                          <FormItem
                                            key={badge.name}
                                            className="flex flex-row items-start space-x-3 space-y-0"
                                          >
                                            <FormControl>
                                              <Checkbox
                                                checked={field.value?.includes(
                                                  badge.name
                                                )}
                                                onCheckedChange={(checked) => {
                                                  return checked
                                                    ? field.onChange([
                                                        ...field.value,
                                                        badge.name,
                                                      ])
                                                    : field.onChange(
                                                        field.value?.filter(
                                                          (value) =>
                                                            value !== badge.name
                                                        )
                                                      );
                                                }}
                                              />
                                            </FormControl>
                                            <FormLabel className="text-sm font-normal">
                                              {badge.name}
                                            </FormLabel>
                                          </FormItem>
                                        );
                                      }}
                                    />
                                  ))}
                                  <FormMessage />
                                </FormItem>
                              )}
                            />
                            <Button type="submit">Submit</Button>
                          </form>
                        </Form>
                      </DialogContent>
                      {newRules.length > 0 && (
                        <Button
                          onClick={confirmRules}
                          variant="outline"
                          size="sm"
                          className="ml-2 bg-green-600 text-white"
                        >
                          Confirm Rules
                        </Button>
                      )}
                    </Dialog>
                  </CardHeader>
                  <CardContent>
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead className="w-[calc(10%)]">#</TableHead>
                          <TableHead className="hidden sm:table-cell w-1/3">
                            Resources
                          </TableHead>
                          <TableHead className="hidden sm:table-cell w-1/3">
                            Verbs
                          </TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {[...roleDetails.role.rules, ...newRules].map(
                          (rule, index) => (
                            <TableRow key={index}>
                              <TableCell className="font-medium">
                                {index + 1}
                              </TableCell>
                              <TableCell>
                                {rule.resources.map((resource) => (
                                  <Badge key={resource} variant="default">
                                    {resource}
                                  </Badge>
                                ))}
                              </TableCell>
                              <TableCell className="flex gap-2 flex-wrap">
                                {rule.verbs.map((verb) => (
                                  <Badge key={verb} variant="default">
                                    {verb}
                                  </Badge>
                                ))}
                              </TableCell>
                            </TableRow>
                          )
                        )}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>
              </TabsContent>
              <TabsContent value="roleBindings">
               
              </TabsContent>
            </Tabs>
          </div>
        </main>
      </div>
    </div>
  );
};

export default RoleDetailsPage;
