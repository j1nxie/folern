"use client";

import { useAtom } from "jotai";
import { redirect } from "next/navigation";
import { type ComponentType, useEffect } from "react";

import { isLoggedInAtom } from "@/atoms/auth";

export default function isAuth<P extends object>(Component: ComponentType<P>): ComponentType<P> {
	return function IsAuth(props: P): React.ReactElement | null {
		const [isLoggedIn, _] = useAtom(isLoggedInAtom);

		useEffect(() => {
			if (!isLoggedIn) {
				return redirect("/");
			}
		}, [isLoggedIn]);

		if (!isLoggedIn) {
			return null;
		}

		return <Component {...props} />;
	};
}
