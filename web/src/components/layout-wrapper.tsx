"use client";

import { useEffect, useState } from "react";

export function LayoutWrapper({
	children,
}: {
	children: React.ReactNode;
}): React.JSX.Element {
	const [isLoading, setIsLoading] = useState(true);

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
