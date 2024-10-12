import { RoleDetails } from "@/interfaces/roleDetails";
import { createResourceHooks } from "./use-resource";

export const {
  useResourceList: useRoleDetails,
  useCreateResource: useCreateRoleDetails,
  useDeleteResource: useDeleteRoleDetails,
} = createResourceHooks<RoleDetails>("RoleDetails");

// Add a new custom hook for role details
export const useRoleDetailsWithParams = (namespace: string, name: string) => {
  return useRoleDetails(namespace, name);
};
