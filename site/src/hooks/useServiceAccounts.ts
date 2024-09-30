import { ServiceAccount } from "@/interfaces/serviceAccount";
import { createResourceHooks } from "./use-resource";

export const {
  useResourceList: useServiceAccounts,
  useCreateResource: useCreateServiceAccount,
  useDeleteResource: useDeleteServiceAccount,
} = createResourceHooks<ServiceAccount[]>("ServiceAccounts");
