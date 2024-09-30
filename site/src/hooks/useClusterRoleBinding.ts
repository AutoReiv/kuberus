import { ClusterRoleBinding } from "@/interfaces/clusterRoleBinding";
import { createResourceHooks } from "./use-resource";

export const {
  useResourceList: useClusterRoleBindings,
  useCreateResource: useCreateClusterRoleBinding,
  useDeleteResource: useDeleteClusterRoleBinding,
} = createResourceHooks<ClusterRoleBinding[]>("ClusterRoleBindings");
