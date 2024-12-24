import type { NextConfig } from "next";

const nextConfig: NextConfig = {
	experimental: {
		serverActions: {
			allowedOrigins: ["localhost:8081"],
		},
	},
};

export default nextConfig;
