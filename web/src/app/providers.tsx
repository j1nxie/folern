"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Provider } from "jotai";
import { useState } from "react";

export function Providers({ children }: { children: React.ReactNode }): React.JSX.Element {
	const [queryClient] = useState(
		() => new QueryClient({
			defaultOptions: {
				mutations: {
					retry: 1,
				},
				queries: {
					retry: 1,
					refetchOnWindowFocus: true,
				},
			},
		}),
	);

	return (
		<Provider>
			<QueryClientProvider client={queryClient}>
				{children}
			</QueryClientProvider>
		</Provider>
	);
}
