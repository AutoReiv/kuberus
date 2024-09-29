
import { ClusterRole } from "@/interfaces/clusterRole";
import { createResourceHooks } from "./use-resource";


export const {
  useResourceList: useClusterRoles,
  useCreateResource: useCreateClusterRole,
  useDeleteResource: useDeleteClusterRoles,
} = createResourceHooks<ClusterRole>("ClusterRoles");

