"use client";

import { useQuery } from "@tanstack/react-query";
import { BentoGrid, BentoGridItem } from "./_components/bento-grid";

interface RoleDetails {
  role: {
    metadata: {
      creationTimestamp: string;
      managedFields: [
        {
          apiVersion: string;
          fieldsType: string;
          fieldsV1: {
            f: {
              rules: [
                {
                  apiGroups: string[];
                  resources: string[];
                  resourceNames: string[];
                  verbs: string[];
                }
              ];
            };
          };
          manager: string;
          operation: string;
          time: string;
        }
      ];
      name: string;
      namespace: string;
      resourceVersion: string;
      uid: string;
    };
    rules: [
      {
        apiGroups: string[];
        resources: string[];
        resourceNames: string[];
        verbs: string[];
      }
    ];
  };
}

const fetchRoleDetails = async (namespace: string, name: string) => {
  const URL = `http://localhost:8080/api/roles/details?roleName=${name}&namespace=${namespace}`;
  const response = await fetch(URL, {
    method: "GET",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  });
  if (!response.ok) {
    throw new Error("Failed to fetch role details");
  }
  const data = await response.json();
  console.log(data, " DATA ****");
  return data;
};

const RoleDetailsPage = ({
  params,
}: {
  params: { namespace: string; name: string };
}) => {
  const { namespace, name } = params;

  const {
    data: roleDetails,
    isLoading,
    error,
  } = useQuery<RoleDetails, Error>({
    queryKey: ["roleDetails", namespace, name],
    queryFn: () => fetchRoleDetails(namespace, name),
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  const Skeleton = () => (
    <div className="flex flex-1 w-full h-full min-h-[6rem] rounded-xl bg-gradient-to-br from-neutral-200 dark:from-neutral-900 dark:to-neutral-800 to-neutral-100"></div>
  );

  const items = [
    {
      title: "Metadata",
      description: (
        <div>
          <ul className="space-y-2">
            <li className="flex items-center space-x-2">
              <span className="font-bold">Name:</span>
              <span>{roleDetails?.role.metadata.name}</span>
            </li>
            <li className="flex items-center space-x-2">
              <span className="font-bold">Namespace:</span>
              <span>{roleDetails?.role.metadata.namespace}</span>
            </li>
            <li className="flex items-center space-x-2">
              <span className="font-bold">Created At:</span>
              <span>{new Date(roleDetails?.role.metadata.creationTimestamp).toLocaleString()}</span>
            </li>
            {/* <li className="flex items-center space-x-2">
              <span className="font-bold">Updated At:</span>
              <span>{new Date(roleDetails?.role.metadata.updatedAt).toLocaleString()}</span>
            </li> */}
            <li className="flex items-center space-x-2">
              <span className="font-bold">Description:</span>
              {/* <span>{roleDetails?.description || "N/A"}</span> */}
            </li>
          </ul>
        </div>
      ),
      // header: <Skeleton />,
      // icon: <IconClipboardCopy className="h-4 w-4 text-neutral-500" />,
    },
    {
      title: 'Rules',
      description: (
        <div>
          <ul className="space-y-2">
            {/* {roleDetails?.role.rules.map((rule, index) => (
              <li key={index} className="flex items-center space-x-2">
                <span className="font-bold">{}</span>
                <span>
                  {rule.apiGroups.join(", ")} {rule.resources.join(", ")} {rule.verbs.join(", ")}
                </span>
              </li>
            ))} */}
          </ul>
        </div>
      )
      // header: <Skeleton />,
      // icon: <IconFileBroken className="h-4 w-4 text-neutral-500" />,
    },
    {
      title: "The Art of Design",
      description: "Discover the beauty of thoughtful and functional design.",
      header: <Skeleton />,
      // icon: <IconSignature className="h-4 w-4 text-neutral-500" />,
    },
    {
      title: "The Power of Communication",
      description:
        "Understand the impact of effective communication in our lives.",
      header: <Skeleton />,
      // icon: <IconTableColumn className="h-4 w-4 text-neutral-500" />,
    },
    {
      title: "The Pursuit of Knowledge",
      description: "Join the quest for understanding and enlightenment.",
      header: <Skeleton />,
      // icon: <IconArrowWaveRightUp className="h-4 w-4 text-neutral-500" />,
    },
    {
      title: "The Joy of Creation",
      description: "Experience the thrill of bringing ideas to life.",
      header: <Skeleton />,
      // icon: <IconBoxAlignTopLeft className="h-4 w-4 text-neutral-500" />,
    },
    {
      title: "The Spirit of Adventure",
      description: "Embark on exciting journeys and thrilling discoveries.",
      header: <Skeleton />,
      // icon: <IconBoxAlignRightFilled className="h-4 w-4 text-neutral-500" />,
    },
  ];

  return (
    <BentoGrid className="max-w-4xl mx-auto">
      {items.map((item, i) => (
        <BentoGridItem
          key={i}
          title={item.title}
          description={item.description}
          header={item.header}
          // icon={item.icon}
          className={i === 3 || i === 6 ? "md:col-span-2" : ""}
        />
      ))}
    </BentoGrid>
  );
};

export default RoleDetailsPage;
