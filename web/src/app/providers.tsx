"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Provider } from "jotai";
import { ThemeProvider } from "next-themes";
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
			<ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
				<QueryClientProvider client={queryClient}>
					{children}
				</QueryClientProvider>
			</ThemeProvider>
		</Provider>
	);
}
