"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { CheckCircle2, Copy, XCircle } from "lucide-react";
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
import { useState } from "react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogFooter,
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
import "@xyflow/react/dist/style.css";
import yaml from "js-yaml";
import LiveYAMLViewer from "./_components/LiveYAMLViewer";
import { toast } from "sonner";
import { Label } from "@/components/ui/label";
import { useRouter } from "next/navigation";
import { apiClient } from "@/lib/apiClient";

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

const RoleDetailsPage = ({
  params,
}: {
  params: { namespace: string; name: string };
}) => {
  const { namespace, name } = params;
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [newRules, setNewRules] = useState<Rules[]>([]);
  const [activatedVerbsByRule, setActivatedVerbsByRule] = useState<{[key: number]: string[];}>({});
  const [isDuplicateDialogOpen, setIsDuplicateDialogOpen] = useState(false);
  const router = useRouter();
  const queryClient = useQueryClient();

  const { data: roleDetails, isLoading, error, refetch: refetchRoleDetails } = useQuery<RoleDetails, Error>({
    queryKey: ["roleDetails", namespace, name],
    queryFn: () => apiClient.getRoleDetails(namespace, name),
  });

  const { data: resources } = useQuery<Resources, Error>({
    queryKey: ["resources"],
    queryFn: () => apiClient.getResources(),
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

  const toggleVerb = async (ruleIndex: number, verb: string) => {
    const ruleVerbs = roleDetails.role.rules[ruleIndex].verbs;
    const updatedRuleVerbs = ruleVerbs.includes(verb)
      ? ruleVerbs.filter((v) => v !== verb)
      : [...ruleVerbs, verb];

    try {
      const updatedRules = await updateRoleRule(ruleIndex, updatedRuleVerbs);
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

      // Update the activatedVerbsByRule state
      setActivatedVerbsByRule((prev) => ({
        ...prev,
        [ruleIndex]: updatedRuleVerbs,
      }));
    } catch (error) {
      console.error("Failed to update role:", error);
    }
  };

  const onSubmit = (data: NewRuleFormValues) => {
    const newRule: Rules = {
      apiGroups: [""],
      resources: [data.resources],
      resourceNames: [],
      verbs: data.verbs,
    };

    // Check if the rule already exists
    const isDuplicate = [...roleDetails.role.rules, ...newRules].some(
      (rule) =>
        rule.resources.join(",") === newRule.resources.join(",") &&
        rule.verbs.join(",") === newRule.verbs.join(",")
    );

    if (!isDuplicate) {
      setNewRules([...newRules, newRule]);
      setIsDialogOpen(false);
      form.reset();
    } else {
      // Show a toast notification for duplicate rule
      toast(
        <div className="flex items-center justify-start gap-4">
          <XCircle className="text-red-500" />
          <span>{`Cannot add in a duplicate rule.`}</span>
        </div>
      );
    }
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

  const updateRoleRule = async (ruleIndex: number, updatedVerbs: string[]) => {
    const updatedRules = [...roleDetails.role.rules];
    updatedRules[ruleIndex].verbs = updatedVerbs;

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

    return updatedRules;
  };

  const duplicateRole = async (newNamespace: string, newName: string) => {
    const newRoleData = {
      metadata: {
        ...(() => {
          const { resourceVersion, ...rest } = roleDetails.role.metadata;
          return rest;
        })(),
        name: newName,
        namespace: newNamespace,
      },
      rules: roleDetails.role.rules,
    };

    const URL = `http://localhost:8080/api/roles?namespace=${newNamespace}&name=${newName}`;
    try {
      const response = await fetch(URL, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newRoleData),
      });

      if (response.ok) {
        toast(
          <div className="flex items-center justify-start gap-4">
            <CheckCircle2 className="text-green-500" />
            <span>{`Successfully duplicated ${newName} in ${newNamespace}`}</span>
          </div>
        );

        router.push("/dashboard/roles");
      } else {
        throw new Error("Failed to duplicate role");
      }
    } catch (error) {
      toast(
        <div className="flex items-center justify-start gap-4">
          <XCircle className="text-red-500" />
          <span>{`${error}`}</span>
        </div>
      );
    }
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

  const DuplicateRoleDialog = ({ isOpen, onClose, onDuplicate }) => {
    const [newNamespace, setNewNamespace] = useState("");
    const [newName, setNewName] = useState("");

    const handleDuplicate = () => {
      onDuplicate(newNamespace, newName);
      onClose();
    };

    return (
      <Dialog open={isOpen} onOpenChange={onClose}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Duplicate Role</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="new-namespace" className="text-right">
                New Namespace
              </Label>
              <Input
                id="new-namespace"
                value={newNamespace}
                onChange={(e) => setNewNamespace(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="new-name" className="text-right">
                New Name
              </Label>
              <Input
                id="new-name"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                className="col-span-3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button onClick={handleDuplicate}>Duplicate</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  };

  const handleYamlUpdate = async (updatedRules: any, updatedMetadata: any) => {
    const updatedRoleData = {
      metadata: updatedMetadata,
      rules: updatedRules,
    };

    try {
      const URL = `http://localhost:8080/api/roles?namespace=${updatedMetadata.namespace}&name=${updatedMetadata.name}`;
      const response = await fetch(URL, {
        method: "PUT",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedRoleData),
      });

      if (response.ok) {
        refetchRoleDetails();
        toast(
          <div className="flex items-center justify-start gap-4">
            <CheckCircle2 className="text-green-500" />
            <span>Role updated successfully.</span>
          </div>
        );
      } else {
        throw new Error("Failed to update role");
      }
    } catch (error) {
      toast(
        <div className="flex items-center justify-start gap-4">
          <XCircle className="text-red-500" />
          <span>{`Failed to update role: ${error}`}</span>
        </div>
      );
    }
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
                    Role Details
                    <Button
                      onClick={() => setIsDuplicateDialogOpen(true)}
                      variant="outline"
                      size="sm"
                    >
                      <Copy className="h-4 w-4 mr-2" />
                      Duplicate Role
                    </Button>
                  </div>
                  <DuplicateRoleDialog
                    isOpen={isDuplicateDialogOpen}
                    onClose={() => setIsDuplicateDialogOpen(false)}
                    onDuplicate={duplicateRole}
                  />
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
                  {/* <Separator></Separator> */}
                  {/* <PermissionImpactAnalysis rules={roleDetails.role.rules} /> */}
                  <Separator></Separator>
                  <LiveYAMLViewer
                    rules={[...roleDetails.role.rules, ...newRules]}
                    metadata={{
                      name: roleDetails.role.metadata.name,
                      namespace: roleDetails.role.metadata.namespace,
                    }}
                    onUpdate={handleYamlUpdate}
                  />
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
                  <TabsTrigger value="roleDiagram">Role Diagram</TabsTrigger>
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
                                {verbs.map((verb) => (
                                  <Badge
                                    key={verb.name}
                                    variant={
                                      rule.verbs.includes(verb.name)
                                        ? "success"
                                        : "secondary"
                                    }
                                    className={`cursor-pointer ${
                                      rule.verbs.includes(verb.name)
                                        ? ""
                                        : "opacity-50"
                                    }`}
                                    onClick={() => toggleVerb(index, verb.name)}
                                  >
                                    {verb.name}
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
              <TabsContent value="roleDiagram">
                <Card>
                  <CardHeader>
                    <CardTitle>Role Flow Diagram</CardTitle>
                    <CardDescription>
                      Visual representation of role permissions
                    </CardDescription>
                  </CardHeader>
                  <CardContent></CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        </main>
      </div>
    </div>
  );
};

export default RoleDetailsPage;
