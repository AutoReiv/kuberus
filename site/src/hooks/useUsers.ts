import { User } from "@/interfaces/user";
import { createResourceHooks } from "./use-resource";

export const {
  useResourceList: useUsers,
  useCreateResource: useCreateUser,
  useDeleteResource: useDeleteUser,
} = createResourceHooks<User[]>("Users");
