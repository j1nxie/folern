"use client";

import { useAtom } from "jotai";
import { useRouter } from "next/navigation";
import { use, useEffect } from "react";

import { authDiscordCallback, isLoggedInAtom } from "@/atoms/auth";

export default function AuthCallback({ searchParams }: { searchParams: Promise<Record<string, string | undefined>> }): React.JSX.Element {
	const { code, state } = use(searchParams);
	const router = useRouter();
	const [{ mutate: processAuth, status }] = useAtom(authDiscordCallback);
	const [_, setIsLoggedIn] = useAtom(isLoggedInAtom);

	useEffect(() => {
		if (code && state && status === "idle") {
			processAuth({ code, state }, {
				onSuccess: () => {
					setIsLoggedIn(true);
					router.replace("/");
				},
			});
		}
	}, [code, state, router, processAuth, status, setIsLoggedIn]);

	if (status === "error") {
		throw new Error("Something wrong happened while logging you in.");
	}

	if (status === "pending") {
		return (
			<div className="my-auto flex grow flex-col items-center justify-center gap-2">
				<p>Logging you in through Discord...</p>
			</div>
		);
	}

	return <></>;
}
