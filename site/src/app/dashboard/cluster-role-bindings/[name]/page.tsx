"use client";
import { useParams } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useQuery } from "@tanstack/react-query";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { apiClient } from "@/lib/apiClient";

interface ClusterRoleBindingDetail {
  clusterRole: {
    metadata: {
      name: string;
      uid: string;
      resourceVersion: string;
      creationTimestamp: string;
      labels: {
        [key: string]: string;
      };
      annotations: {
        [key: string]: string;
      };
      managedFields: {
        manager: string;
        operation: string;
        apiVersion: string;
        time: string;
        fieldsType: string;
        fieldsV1: {
          [key: string]: any;
        };
      }[];
    };
    rules: {
      verbs: string[];
      apiGroups: string[];
      resources: string[];
      nonResourceURLs?: string[];
    }[];
  };
  clusterRoleBindings: {
    metadata: {
      name: string;
      uid: string;
      resourceVersion: string;
      creationTimestamp: string;
      labels?: {
        [key: string]: string;
      };
      annotations?: {
        [key: string]: string;
      };
      managedFields: {
        manager: string;
        operation: string;
        apiVersion: string;
        time: string;
        fieldsType: string;
        fieldsV1: {
          [key: string]: any;
        };
      }[];
    };
    subjects: {
      kind: string;
      apiGroup?: string;
      name: string;
    }[];
    roleRef: {
      apiGroup: string;
      kind: string;
      name: string;
    };
  }[];
}

const ClusterRoleBindingDetailsPage = (params: {
  params: { name: string };
}) => {
  const { name } = useParams();

  const {
    data: clusterRoleBindingDetails,
    isLoading,
    error,
  } = useQuery<ClusterRoleBindingDetail>({
    queryKey: ["clusterRoleDetails", name],
    queryFn: () => apiClient.getClusterRoleBindingDetails(name as string),
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  if (!clusterRoleBindingDetails) return <div>No data available</div>;

  const { clusterRole, clusterRoleBindings } = clusterRoleBindingDetails;

  return (
    <div className="flex min-h-screen w-full flex-col bg-muted/40 p-4 space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Cluster Role: {clusterRole.metadata.name}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <h3 className="font-semibold">Metadata</h3>
              <ul className="list-disc list-inside">
                <li>UID: {clusterRole.metadata.uid}</li>
                <li>
                  Created:{" "}
                  {new Date(
                    clusterRole.metadata.creationTimestamp
                  ).toLocaleString()}
                </li>
                <li>
                  Resource Version: {clusterRole.metadata.resourceVersion}
                </li>
              </ul>
            </div>
            <div>
              <h3 className="font-semibold">Labels</h3>
              {Object.entries(clusterRole.metadata.labels || {}).map(
                ([key, value]) => (
                  <div key={key}>
                    {key}: {value}
                  </div>
                )
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Rules</CardTitle>
        </CardHeader>
        <CardContent>
          <Accordion type="single" collapsible className="w-full">
            {clusterRole.rules.map((rule, index) => (
              <AccordionItem value={`item-${index}`} key={index}>
                <AccordionTrigger>Rule {index + 1}</AccordionTrigger>
                <AccordionContent>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <h4 className="font-semibold">API Groups</h4>
                      <ul className="list-disc list-inside">
                        {rule.apiGroups &&
                          rule.apiGroups.map((group, i) => (
                            <li key={i}>{group}</li>
                          ))}
                      </ul>
                    </div>
                    <div>
                      <h4 className="font-semibold">Resources</h4>
                      <ul className="list-disc list-inside">
                        {rule.resources &&
                          rule.resources.map((resource, i) => (
                            <li key={i}>{resource}</li>
                          ))}
                      </ul>
                    </div>
                    <div>
                      <h4 className="font-semibold">Verbs</h4>
                      <ul className="list-disc list-inside">
                        {rule.verbs &&
                          rule.verbs.map((verb, i) => <li key={i}>{verb}</li>)}
                      </ul>
                    </div>
                    {rule.nonResourceURLs && (
                      <div>
                        <h4 className="font-semibold">Non-Resource URLs</h4>
                        <ul className="list-disc list-inside">
                          {rule.nonResourceURLs.map((url, i) => (
                            <li key={i}>{url}</li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                </AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Cluster Role Bindings</CardTitle>
        </CardHeader>
        <CardContent>
          <Accordion type="single" collapsible className="w-full">
            {clusterRoleBindings.map((binding, index) => (
              <AccordionItem value={`binding-${index}`} key={index}>
                <AccordionTrigger>{binding.metadata.name}</AccordionTrigger>
                <AccordionContent>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <h4 className="font-semibold">Metadata</h4>
                      <ul className="list-disc list-inside">
                        <li>UID: {binding.metadata.uid}</li>
                        <li>
                          Created:{" "}
                          {new Date(
                            binding.metadata.creationTimestamp
                          ).toLocaleString()}
                        </li>
                        <li>
                          Resource Version: {binding.metadata.resourceVersion}
                        </li>
                      </ul>
                    </div>
                    <div>
                      <h4 className="font-semibold">Role Reference</h4>
                      <ul className="list-disc list-inside">
                        <li>Kind: {binding.roleRef.kind}</li>
                        <li>Name: {binding.roleRef.name}</li>
                        <li>API Group: {binding.roleRef.apiGroup}</li>
                      </ul>
                    </div>
                    <div className="col-span-2">
                      <h4 className="font-semibold">Subjects</h4>
                      <table className="w-full">
                        <thead>
                          <tr>
                            <th className="text-left">Kind</th>
                            <th className="text-left">Name</th>
                            <th className="text-left">API Group</th>
                          </tr>
                        </thead>
                        <tbody>
                          {binding.subjects.map((subject, subIndex) => (
                            <tr key={subIndex}>
                              <td>{subject.kind}</td>
                              <td>{subject.name}</td>
                              <td>{subject.apiGroup}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                </AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </CardContent>
      </Card>
    </div>
  );
};

export default ClusterRoleBindingDetailsPage;
