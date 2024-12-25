"use client";

import { useAtom } from "jotai";
import { useRouter } from "next/navigation";
import { use, useEffect } from "react";

import { authKamaitachiCallback } from "@/atoms/auth";

export default function KamaitachiAuthCallback({ searchParams }: { searchParams: Promise<Record<string, string | undefined>> }): React.JSX.Element {
	const { code } = use(searchParams);
	const router = useRouter();
	const [{ mutate: processKamaitachiAuth, status }] = useAtom(authKamaitachiCallback);

	useEffect(() => {
		if (code) {
			processKamaitachiAuth({ code }, {
				onSuccess: () => {
					router.replace("/");
				},
			});
		}
	}, [code, router, processKamaitachiAuth]);

	if (status === "error") {
		throw new Error("Something wrong happened while logging you in.");
	}

	if (status === "pending") {
		return (
			<div className="my-auto flex grow flex-col items-center justify-center gap-2">
				<p>Logging you in through Kamaitachi...</p>
			</div>
		);
	}

	return <></>;
}
