
import { Namespace } from "@/interfaces/namespace";
import { createResourceHooks } from "./use-resource";

export const {
  useResourceList: useNamespaces,
  useCreateResource: useCreateNamespace,
  useDeleteResource: useDeleteNamespaces,
} = createResourceHooks<Namespace[]>("Namespaces");
