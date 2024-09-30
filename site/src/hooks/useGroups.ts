
import { Group } from "@/interfaces/group";
import { createResourceHooks } from "./use-resource";

export const {
  useResourceList: useGroups,
  useCreateResource: useCreateGroup,
  useDeleteResource: useDeleteGroup,
} = createResourceHooks<Group[]>("Groups");
