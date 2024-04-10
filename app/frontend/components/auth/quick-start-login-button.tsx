"use client";

import { signIn } from "next-auth/react";
import { useTransition } from "react";
import { LoadingButton } from "../ui/loading-button";

type Props = {
	name: string;
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
							endpoint: "signin",
						},
					);
				});
			}}
		>
			Log in as {name}
		</LoadingButton>
	);
};
