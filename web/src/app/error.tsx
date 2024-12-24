"use client";

import Link from "next/link";

export default function ErrorPage({ error }: { error: Error & { digest?: string } }): React.JSX.Element {
	return (
		<div className="my-auto flex grow flex-col items-center justify-center gap-2">
			<h1 className="text-2xl font-bold">Error!</h1>
			<p>{error.message}</p>
			<Link href="/" className="text-rose-400">Go back home.</Link>
		</div>
	);
}
