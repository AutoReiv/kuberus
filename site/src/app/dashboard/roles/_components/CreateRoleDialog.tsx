import React, { useEffect, useState } from "react";
import { PlusCircle, Trash2 } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { motion, AnimatePresence } from "framer-motion";
import { z } from "zod";
import { useFieldArray, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/apiClient";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ScrollArea } from "@/components/ui/scroll-area";

const dnsNameRegex =
  /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/;

const resourceSchema = z.object({
  name: z.string().min(1, "Resource name is required"),
  verbs: z.array(z.string()).min(1, "At least one verb is required"),
});

const roleSchema = z
  .object({
    roleName: z
      .string()
      .min(1, "Role name is required")
      .regex(dnsNameRegex, "Role name must be DNS compliant"),
    namespace: z
      .string()
      .min(1, "Namespace is required")
      .regex(dnsNameRegex, "Role name must be DNS compliant"),
    resources: z
      .array(resourceSchema)
      .min(1, "At least one resource is required"),
    apiGroup: z.string().optional(),
  })
  .superRefine((data, ctx) => {
    if (data.resources.some((resource) => resource.verbs.length === 0)) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Each resource must have at least one verb selected",
        path: ["resources"],
      });
    }
  });

const CreateRoleDialog = ({ onSubmit }) => {
  const [showApiGroup, setShowApiGroup] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    control,
    getValues,
    setValue,
    setError,
    clearErrors,
  } = useForm({
    resolver: zodResolver(roleSchema),
    defaultValues: {
      roleName: "",
      namespace: "",
      resources: [{ name: "", verbs: [] }],
      apiGroup: "",
    },
    mode: "onChange",
  });

  // Get Namespaces
  const { data: namespaces } = useQuery({
    queryKey: ["namespaces"],
    queryFn: () => apiClient.getNamespaces(),
  });


  console.log(namespaces, 'namespaces')

  const { data: resources } = useQuery({
    queryKey: ["resources"],
    queryFn: () => apiClient.getResources(),
  });

  // Inside the component:
  const { fields, append, remove } = useFieldArray({
    control,
    name: "resources",
  });

  useEffect(() => {
    if (fields.length > 0) {
      clearErrors("resources");
    }
  }, [fields, clearErrors]);

  const availableVerbs = [
    "Get",
    "List",
    "Watch",
    "Create",
    "Update",
    "Patch",
    "Delete",
  ];

  const addResource = () => {
    append({ name: "", verbs: [] });
  };

  const removeResource = (index) => {
    remove(index);
  };

  const toggleVerb = (index: number, verb: string) => {
    const currentVerbs = getValues(`resources.${index}.verbs`);
    const updatedVerbs = currentVerbs.includes(verb)
      ? currentVerbs.filter((v) => v !== verb)
      : [...currentVerbs, verb];

    setValue(`resources.${index}.verbs`, updatedVerbs);

    if (updatedVerbs.length === 0) {
      setError(`resources.${index}.verbs`, {
        type: "manual",
        message: "At least one verb is required",
      });
    } else {
      clearErrors(`resources.${index}.verbs`);
    }
  };

  const renderError = (err) => {
    if (typeof err === "string") return <p>{err}</p>;
    if (Array.isArray(err)) return err.map((e, i) => <p key={i}>{e}</p>);
    if (typeof err === "object") return <p>{JSON.stringify(err.message)}</p>;
    return null;
  };

  const LabelWithTooltip = ({ htmlFor, children, error }) => (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Label
            htmlFor={htmlFor}
            className={`whitespace-nowrap ${error ? "text-red-500" : ""}`}
          >
            {children}
          </Label>
        </TooltipTrigger>
        {error &&
          (console.log(error),
          (<TooltipContent>{renderError(error)}</TooltipContent>))}
      </Tooltip>
    </TooltipProvider>
  );



  return (
    <div className="w-full mx-auto">
      <CardHeader className="py-2 px-0">
        <CardTitle>Create RBAC Role</CardTitle>
      </CardHeader>
      <CardContent className="px-0">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="flex items-center space-x-2">
            <div className="flex items-center gap-4 flex-1">
              <LabelWithTooltip
                htmlFor="roleName"
                error={errors.roleName?.message}
              >
                Role Name *
              </LabelWithTooltip>
              <Input
                id="roleName"
                {...register("roleName")}
                placeholder="Enter role name"
              />
            </div>

            <div className="flex items-center gap-4 flex-1">
              <LabelWithTooltip
                htmlFor="namespace"
                error={errors.namespace?.message}
              >
                Namespace *
              </LabelWithTooltip>
              <Select onValueChange={(value) => setValue("namespace", value)}>
                <SelectTrigger>
                  <SelectValue placeholder="Select namespace" />
                </SelectTrigger>
                <SelectContent>
                  {namespaces?.items.map((namespace) => (
                    <SelectItem
                      key={namespace.metadata.name}
                      value={namespace.metadata.name}
                    >
                      {namespace.metadata.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <ScrollArea className="h-auto max-h-[40vh] overflow-y-auto">
            <AnimatePresence>
              {fields.map((field, index) => (
                <motion.div
                  key={field.id}
                  initial={{ opacity: 0, y: -20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  transition={{ duration: 0.3 }}
                  className="mb-4 last-of-type:mb-0"
                >
                  <Card className="p-4">
                    <div className="flex items-center gap-4">
                      <div className="flex items-center gap-4 flex-1">
                        <LabelWithTooltip
                          htmlFor={`resource-${index}`}
                          error={errors.resources?.[index]?.name?.message}
                        >
                          Resource *
                        </LabelWithTooltip>
                        <Select
                          onValueChange={(value) =>
                            setValue(`resources.${index}.name`, value)
                          }
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Select resource" />
                          </SelectTrigger>
                          <SelectContent>
                            {resources?.resources.map((resource) => (
                              <SelectItem key={resource} value={resource}>
                                {resource}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="flex items-center gap-4 flex-1">
                        <Label htmlFor="apiGroup" className="whitespace-nowrap">
                          API Group
                        </Label>
                        <Input
                          id="apiGroup"
                          {...register("apiGroup")}
                          placeholder="Enter API group - (Optional)"
                        />
                      </div>
                    </div>
                    <div className="mt-2 flex justify-between items-center gap-4">
                      <div className="flex items-center gap-4">
                        <LabelWithTooltip
                          htmlFor="roleName"
                          error={errors.resources?.[index]?.verbs}
                        >
                          Verbs *
                        </LabelWithTooltip>
                        <Label className="whitespace-nowrap mr-6"></Label>
                        <div className="flex flex-wrap gap-2 mt-1">
                          {availableVerbs.map((verb) => (
                            <Badge
                              key={verb}
                              variant={
                                getValues(`resources.${index}.verbs`).includes(
                                  verb
                                )
                                  ? "default"
                                  : "outline"
                              }
                              className={`cursor-pointer hover:bg-green-500 hover:text-white hover:border-green-500 ${
                                getValues(`resources.${index}.verbs`).includes(
                                  verb
                                )
                                  ? "bg-green-500 text-white"
                                  : ""
                              }`}
                              onClick={() => toggleVerb(index, verb)}
                            >
                              {verb}
                            </Badge>
                          ))}
                        </div>
                      </div>
                      {index > 0 && (
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="h-6"
                          onClick={() => removeResource(index)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      )}
                    </div>
                  </Card>
                </motion.div>
              ))}
            </AnimatePresence>
          </ScrollArea>
          <Button
            type="button"
            variant="outline"
            onClick={addResource}
            className="w-full"
          >
            <PlusCircle className="mr-2 h-4 w-4" /> Add Resource
          </Button>

          <AnimatePresence>
            {showApiGroup && (
              <motion.div
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: "auto" }}
                exit={{ opacity: 0, height: 0 }}
                transition={{ duration: 0.3 }}
              ></motion.div>
            )}
          </AnimatePresence>

          <Button type="submit" className="w-full">
            Create Role
          </Button>
        </form>
      </CardContent>
    </div>
  );
};

export default CreateRoleDialog;
