import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/apiClient";
import { toast } from "sonner";

const handleMutationResult = (
  resourceName: string,
  action: string,
  stage: string,
  isSuccess: boolean,
  result: any
) => {
  let message;
  switch (stage) {
    case "start":
      message = `${
        action.charAt(0).toUpperCase() + action.slice(1)
      }ing ${resourceName}...`;
      toast.info(message);
      break;
    case "processing":
      message = `Processing ${resourceName} ${action}...`;
      toast.info(message);
      break;
    case "complete":
      if (isSuccess) {
        switch (action) {
          case "create":
            message = `New ${resourceName} has been created successfully`;
            break;
          case "delete":
            message = `${resourceName} has been removed successfully`;
            break;
          case "update":
            message = `${resourceName} has been updated successfully`;
            break;
          default:
            message = `${resourceName} action completed successfully`;
        }
        toast.success(message);
      } else {
        message = `Error ${action}ing ${resourceName}: ${result.message}`;
        toast.error(message);
      }
      break;
    default:
      message = `Unknown stage for ${resourceName} ${action}`;
      toast.warning(message);
  }
};

export function createResourceHooks<T>(resourceName: string) {
  const useResourceList = () => {
    return useQuery({
      queryKey: [resourceName],
      queryFn: () => apiClient[`get${resourceName}`](),
    });
  };

  const createMutationHook = (
    action: string,
    mutationFn: (data: any) => Promise<any>
  ) => {
    return (options?: { onSuccess?: () => void }) => {
      const queryClient = useQueryClient();

      return useMutation({
        mutationFn,
        onMutate: () => {
          handleMutationResult(resourceName, action, "start", true, {});
        },
        onSuccess: () => {
          handleMutationResult(resourceName, action, "processing", true, {});
          options?.onSuccess?.();
        },
        onSettled: (data, error) => {
          queryClient.invalidateQueries({ queryKey: [resourceName] });
          handleMutationResult(
            resourceName,
            action,
            "complete",
            !error,
            error || data
          );
          queryClient.refetchQueries({ queryKey: [resourceName] });
        },
      });
    };
  };

  const useCreateResource = createMutationHook("create", (data: any) =>
    apiClient[`create${resourceName}`](data)
  );

  const useDeleteResource = createMutationHook(
    "delete",
    (resources: string | T[]) => {
      if (Array.isArray(resources)) {
        return Promise.all(
          resources.map((resource: any) => {
            apiClient[`delete${resourceName}`](
              resource.metadata?.namespace ||
                resource.namespace ||
                (resources as any).metadata?.name,
              resource.metadata?.name || resource.name
            );
          })
        );
      } else {
        return apiClient[`delete${resourceName}`](
          (resources as any).metadata?.namespace ||
            (resources as any).namespace ||
            (resources as any).metadata?.name,
          (resources as any).metadata?.name || (resources as any).name
        );
      }
    }
  );

  return { useResourceList, useCreateResource, useDeleteResource };
}
