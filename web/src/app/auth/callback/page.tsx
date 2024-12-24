"use client";

import { useRouter } from "next/navigation";
import { use, useEffect, useState } from "react";

import { processCallback } from "@/lib/api/auth/process-callback";

export default function AuthCallback({ searchParams }: { searchParams: Promise<Record<string, string | undefined>> }): React.JSX.Element {
	const { code, state } = use(searchParams);
	const router = useRouter();
	const [isLoading, setIsLoading] = useState(true);

	useEffect(() => {
		async function login(): Promise<void> {
			try {
				if (!code || !state) {
					throw new Error("Invalid callback parameters.");
				}

				await processCallback({ code, state });

				localStorage.setItem("LOGGED_IN", "true");

				router.push("/");
			} catch (e: unknown) {
				throw new Error(`An error occurred when logging you in: ${e}`);
			} finally {
				setIsLoading(false);
			}
		}

		void login();
	}, [code, state, router]);

	if (isLoading) {
		return (
			<div className="my-auto flex grow flex-col items-center justify-center gap-2">
				<p>Logging you in through Discord...</p>
			</div>
		);
	}

	return <></>;
}
