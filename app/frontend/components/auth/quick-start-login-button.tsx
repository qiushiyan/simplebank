"use client";

import { TestUsername, testUsers } from "@/lib/user";
import { signIn } from "next-auth/react";
import { Button } from "../ui/button";

type Props = {
	name: TestUsername;
};

export const QuickStartLoginButton = ({ name }: Props) => {
	return (
		<Button
			onClick={() => {
				signIn(
					"credentials",
					{
						callbackUrl: "/",
					},
					{
						username: name,
					},
				);
			}}
		>
			Log in as {name}
		</Button>
	);
};
