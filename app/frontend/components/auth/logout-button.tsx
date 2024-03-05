"use client";

import { signOut } from "next-auth/react";
import { Button } from "../ui/button";

interface Props extends React.HTMLAttributes<HTMLButtonElement> {}

export const LogoutButton = (props: Props) => {
	return (
		<Button {...props} onClick={() => signOut()}>
			Log out
		</Button>
	);
};
