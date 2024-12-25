"use client";

import { SiDiscord } from "@icons-pack/react-simple-icons";
import { useAtom } from "jotai";
import Image from "next/image";
import Link from "next/link";

import { isLoggedInAtom } from "@/atoms/auth";
import { Button } from "@/components/ui/button";
import { getAuthURL } from "@/lib/api/auth/get-url";

export default function Home(): React.JSX.Element {
	const [isLoggedIn, _] = useAtom(isLoggedInAtom);

	return (
		<div className="my-auto flex grow flex-col items-center justify-center gap-2">
			<Image src="/folern.png" alt="folern" width={200} height={200} className="rounded-xl" priority />

			<div className="my-4 flex flex-col items-center gap-2 px-4 text-center">
				<h1 className="text-2xl font-bold">folern</h1>
				<p>An unofficial tool for tracking and calculating OVER POWER stats for CHUNITHM.</p>
				<p>Powered by data from <Link href="https://kamai.tachi.ac" className="text-rose-400 underline-offset-4 hover:underline">Kamaitachi</Link>.</p>
			</div>

			{!isLoggedIn && (
				<Button onClick={async () => {
					const { url } = await getAuthURL();

					window.location.href = url;
				}}
				>
					<SiDiscord />
					<p>Login with <span className="font-bold">Discord</span></p>
				</Button>
			)}
		</div>
	);
}
