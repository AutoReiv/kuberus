
import { Resources } from "@/interfaces/resources";
import { createResourceHooks } from "./use-resource";

export const {
  useResourceList: useResources,
  useCreateResource: useCreateResource,
  useDeleteResource: useDeleteResource,
} = createResourceHooks<Resources[]>("Resources");