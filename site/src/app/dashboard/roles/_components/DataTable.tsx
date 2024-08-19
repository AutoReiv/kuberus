"use client";

import { MoreHorizontal } from "lucide-react";
import YamlEditor from "@focus-reactive/react-yaml";

import { Button } from "@/components/ui/button";
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
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { create } from "zustand";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { useRouter } from "next/navigation";

const useStore = create((set: any) => ({
  namespaces: [],
  setNamespaces: (namespace) =>
    set(() => {
      namespace;
    }),
}));

export default function DataTable({ roles }) {
  const yaml = require("js-yaml");
  const yamlText = `
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: example-role123
  namespace: default
rules:
- apiGroups:
    - ""
  resources:
    - pods
  verbs:
    - get
    - watch
    - list
`;
  const [formText, setformText] = useState(yamlText);
  const [rolesArray, setRolesArray] = useState(roles);
  const router = useRouter();

  const handleChange = ({ text }) => {
    setformText(text);
  };

  const { data, isError } = useQuery({
    queryKey: ["namespace"],
    queryFn: async () => {
      const URL = "http://localhost:8080/api/namespaces";
      const response = await fetch(URL, {
        method: "GET",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      });
      const data = await response.json();
      return data;
    },
  });

  if (isError) {
    console.log("Error in getting all namespaces.");
  }

  const formSchema = z.object({
    username: z.string().min(2).max(50),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
    },
  });

  const createRole = (event) => {
    const jsonObject = yaml.load(formText);
    const jsonText = JSON.stringify(jsonObject, null, 2);
    const URL = "http://localhost:8080/api/roles";

    const rolePayload = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "Role",
      metadata: {
        name: "example-role",
        namespace: "your-namespace"
      },
      rules: [
        {
          apiGroups: [""],
          resources: ["pods"],
          verbs: ["get", "list", "watch"]
        }
      ]
    };

    fetch(URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: jsonText
    })
    .then(async (response) => {
      console.log(rolesArray)
      const newRole = await response.json();
      setRolesArray([...rolesArray, newRole]);
    });
  };

  const deleteRole = (id, namespace, name) => {
    const URL = `http://localhost:8080/api/roles?namespace=${namespace}&name=${name}`;
    fetch(URL, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    }).then(() => {
      const newArray = rolesArray.filter((role, index) => {
        return index !== id;
      });
      setRolesArray(newArray);
    });
  };

  const goToRolePage = async (role) => {
    await router.push(`/dashboard/roles/${role}`);
  };

  return (
    <Card className="h-full">
      <CardHeader>
        <div className="justify-between item-start flex">
          <div className="flex flex-col gap-4">
            <CardTitle className="font-bold">Roles</CardTitle>
            <CardDescription>Manage User Roles and Permissions</CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <Button>Delete</Button>
            <Dialog>
              <DialogTrigger type="button">Create</DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create a Role</DialogTitle>
                  <DialogDescription>
                    Type in your information for creating a form.
                  </DialogDescription>
                </DialogHeader>

                <Tabs defaultValue="form">
                  <TabsList className="w-full">
                    <TabsTrigger value="form" className="w-full">
                      FORM
                    </TabsTrigger>
                    <TabsTrigger value="yaml" className="w-full">
                      YAML
                    </TabsTrigger>
                  </TabsList>
                  <TabsContent value="form">
                    <Form {...form}>
                      <form
                        onSubmit={form.handleSubmit(createRole)}
                        className="space-y-8"
                      >
                        <FormField
                          control={form.control}
                          name="username"
                          render={({ field }) => (
                            <div>
                              <FormItem className="mb-4">
                                <FormLabel>Name</FormLabel>
                                <FormControl>
                                  <Input placeholder="shadcn" {...field} />
                                </FormControl>
                                <FormDescription>
                                  This is your public display name.
                                </FormDescription>
                                <FormMessage />
                              </FormItem>

                              <FormItem>
                                <FormLabel>Namespace</FormLabel>
                                <FormControl>
                                  <Select>
                                    <SelectTrigger>
                                      <SelectValue placeholder="Namespace" />
                                    </SelectTrigger>
                                    <SelectContent>
                                      {data.map((namespace) => (
                                        <SelectItem
                                          value={namespace.metadata.name}
                                          key={namespace.metadata.uid}
                                        >
                                          {namespace.metadata.name}
                                        </SelectItem>
                                      ))}
                                    </SelectContent>
                                  </Select>
                                </FormControl>
                                <FormDescription>
                                  Select the namespace.
                                </FormDescription>
                                <FormMessage />
                              </FormItem>

                              <FormItem>
                                <FormLabel>Verbs</FormLabel>
                                <FormControl>
                                  <div className="flex items-center justify-start gap-2 flex-wrap">
                                    <Badge variant="outline">Create</Badge>
                                    <Badge variant="outline">Delete</Badge>
                                    <Badge variant="outline">Get</Badge>
                                    <Badge variant="outline">List</Badge>
                                    <Badge variant="outline">Patch</Badge>
                                    <Badge variant="outline">Update</Badge>
                                    <Badge variant="outline">Watch</Badge>
                                  </div>
                                </FormControl>
                                <FormDescription>
                                  Select verb(s)
                                </FormDescription>
                                <FormMessage />
                              </FormItem>
                              <FormItem>
                                <FormLabel>Resources</FormLabel>
                                <FormControl>
                                  <Input placeholder="shadcn" {...field} />
                                </FormControl>
                                <FormDescription>
                                  Select the namespace.
                                </FormDescription>
                                <FormMessage />
                              </FormItem>
                            </div>
                          )}
                        />
                        <Button type="submit">Submit</Button>
                      </form>
                    </Form>
                  </TabsContent>
                  <TabsContent value="yaml">
                    <YamlEditor text={formText} onChange={handleChange} />
                    <Button onClick={() => createRole(event)}>
                      Create Role{" "}
                    </Button>
                  </TabsContent>
                </Tabs>
              </DialogContent>
            </Dialog>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <Table className="h-full">
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Namespace</TableHead>
              <TableHead className="hidden md:table-cell">Created at</TableHead>
              <TableHead>
                <span className="sr-only">Actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rolesArray.map((role, id) => (
              <TableRow key={role.metadata.uid}>
                <TableCell
                  className="font-medium cursor-pointer"
                  onClick={() => goToRolePage(role.metadata.name)}
                >
                  {role.metadata.name}
                </TableCell>
                <TableCell>{role.metadata.namespace}</TableCell>
                <TableCell>
                  {new Date(role.metadata.creationTimestamp).toLocaleDateString(
                    "en-US",
                    {
                      day: "2-digit",
                      month: "short",
                      year: "numeric",
                    }
                  )}
                </TableCell>
                <TableCell>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button aria-haspopup="true" size="icon" variant="ghost">
                        <MoreHorizontal className="h-4 w-4" />
                        <span className="sr-only">Toggle menu</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuLabel>Actions</DropdownMenuLabel>
                      <DropdownMenuItem>Edit</DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() =>
                          deleteRole(
                            id,
                            role.metadata.namespace,
                            role.metadata.name
                          )
                        }
                      >
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
      <CardFooter>
        <div className="text-xs text-muted-foreground">
          Showing <strong>1-10</strong> of <strong>32</strong> products
        </div>
      </CardFooter>
    </Card>
  );
}
