/** @type {import('next').NextConfig} */
const nextConfig = {
  async redirects() {
    return [
      {
        source: "/",
        destination: "/dashboard/roles",
        permanent: true,
      },
      {
        source: "/dashboard",
        destination: "/dashboard/roles",
        permanent: true,
      },
    ];
  },
};

export default nextConfig;
