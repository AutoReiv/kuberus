   # K-RBAC

   **Transform Your Kubernetes Management with K-RBAC!**

   Dive into seamless Kubernetes Role-Based Access Control with K-RBAC, the ultimate tool for managing your cluster's security with ease and precision. Whether you're a seasoned Kubernetes admin or just starting out, K-RBAC offers an intuitive interface that simplifies the complex world of roles, bindings, and permissions.

   **Why Choose K-RBAC?**

   - **Simplified Authorization**: Streamline Kubernetes authorization with an intuitive interface that makes managing roles and permissions straightforward.
   - **Effortless Management**: Navigate and manage users, groups, and namespaces with just a few clicks.
   - **Comprehensive Control**: Create, update, and delete roles and bindings without breaking a sweat.
   - **Secure and Scalable**: Built to handle the demands of modern Kubernetes environments.

   ## Getting Started

   To get started with K-RBAC, follow these simple steps:

   ### Prerequisites

   - Ensure Docker is installed on your machine.

   ### Running the Application

   1. **Pull the Docker Image**:
      ```bash
      docker pull ghcr.io/autoreiv/k-rbac:latest

   2.  **Pull the Docker Image**:
      ```bash
      docker run -p 80:80 ghcr.io/autoreiv/k-rbac:latest

   3. **Access the Application**: 
      Open your web browser and navigate to http://localhost:80

   ### Configuration
   
   To connect K-RBAC to your Kubernetes cluster, ensure your kubeconfig file is located at ~/.kube/config or set the KUBECONFIG environment variable to point to your kubeconfig file.

   ### Contributing
   
   We welcome contributions! Please feel free to submit a Pull Request.

   ### License

   This project is licensed under the MIT License - see the LICENSE file for details.

