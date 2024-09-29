import { Role } from "@/interfaces/role";
import { createResourceHooks } from "./use-resource";

export const {
  useResourceList: useRoles,
  useCreateResource: useCreateRole,
  useDeleteResource: useDeleteRoles,
} = createResourceHooks<Role[]>("Roles");