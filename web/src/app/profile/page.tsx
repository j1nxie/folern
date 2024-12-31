"use client";

import { useAtom } from "jotai";

import { userDataAtom } from "@/atoms/user";
import isAuth from "@/components/is-auth";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";

function ProfilePage(): React.JSX.Element {
	const [{ data, isLoading, isPending }] = useAtom(userDataAtom);

	return (
		<div className="mx-auto flex w-full max-w-screen-xl flex-col gap-8 p-8 md:flex-row">
			<div className="flex flex-col gap-3 pr-4">
				<p className="text-3xl font-bold">
					{data?.username}
				</p>
				<p>Last updated: Today</p>
				<Button>Sync</Button>
			</div>

			<div className="flex flex-1 flex-col">
				<Select defaultValue="genre">
					<SelectTrigger className="mb-4 w-full">
						<SelectValue placeholder="Select sorting category" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="genre">Genre</SelectItem>
						<SelectItem value="version">Version</SelectItem>
					</SelectContent>
				</Select>

				<div className="mb-4 w-full">
					<h1 className="text-xl font-bold">ALL</h1>
					<Separator className="my-2" />
					Data here
				</div>

				<div className="mb-4 w-full">
					<h1 className="text-xl font-bold">ALL</h1>
					<Separator className="my-2" />
					Data here
				</div>

				<div className="mb-4 w-full">
					<h1 className="text-xl font-bold">ALL</h1>
					<Separator className="my-2" />
					Data here
				</div>

				<div className="mb-4 w-full">
					<h1 className="text-xl font-bold">ALL</h1>
					<Separator className="my-2" />
					Data here
				</div>
			</div>
		</div>
	);
}

export default isAuth(ProfilePage);
