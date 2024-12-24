import Link from "next/link";

import { Button } from "@/components/ui/button";

export default function NotFound(): React.JSX.Element {
	return (
		<div className="my-auto flex grow flex-col items-center justify-center gap-2">
			<h1 className="text-2xl font-bold">Not Found</h1>
			<p>The requested page could not be found.</p>
			<Button variant="link" asChild>
				<Link href="/" className="text-rose-400">Go back home.</Link>
			</Button>
		</div>
	);
}
