import { RoleBinding } from "@/interfaces/roleBinding";
import { createResourceHooks } from "./use-resource";

export const {
  useResourceList: useRoleBindings,
  useCreateResource: useCreateRoleBinding,
  useDeleteResource: useDeleteRoleBinding,
} = createResourceHooks<RoleBinding[]>("RoleBindings");