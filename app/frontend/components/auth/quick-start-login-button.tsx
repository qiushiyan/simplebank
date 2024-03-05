"use client";

import { TestUsername, testUsers } from "@/lib/user";
import { signIn } from "next-auth/react";
import { useState, useTransition } from "react";
import { Button } from "../ui/button";
import { LoadingButton } from "../ui/loading-button";
import { Spinner } from "../ui/spinner";

type Props = {
	name: TestUsername;
};

export const QuickStartLoginButton = ({ name }: Props) => {
	const [pending, startTransition] = useTransition();
	return (
		<LoadingButton
			loading={pending}
			onClick={async () => {
				startTransition(async () => {
					await signIn(
						"credentials",
						{
							callbackUrl: "/",
						},
						{
							username: name,
						},
					);
				});
			}}
		>
			Log in as {name}
		</LoadingButton>
	);
};
