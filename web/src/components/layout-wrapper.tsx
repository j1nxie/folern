"use client";

import { useAtom, useAtomValue } from "jotai";
import { useEffect, useState } from "react";

import { authVerificationAtom, isLoggedInAtom, logoutAtom } from "@/atoms/auth";

export function LayoutWrapper({
	children,
}: {
	children: React.ReactNode;
}): React.JSX.Element {
	const [isLoading, setIsLoading] = useState(true);
	const [isLoggedIn, setIsLoggedIn] = useAtom(isLoggedInAtom);
	const { error } = useAtomValue(authVerificationAtom);
	const { mutate: logout } = useAtomValue(logoutAtom);

	useEffect(() => {
		if (error && isLoggedIn) {
			setIsLoggedIn(false);
			logout();
		}
	}, [error, isLoggedIn, setIsLoggedIn, logout]);

	useEffect(() => {
		setIsLoading(false);
	}, []);

	if (isLoading) {
		return (
			<div className="flex grow">
				<div className="my-auto flex grow flex-col items-center justify-center gap-2">
					<div className="size-[200px] animate-pulse rounded-xl bg-muted" />
				</div>
			</div>
		);
	}

	return <>{children}</>;
}
