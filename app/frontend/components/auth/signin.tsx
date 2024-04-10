"use client";

import { signIn } from "next-auth/react";
import { useFormStatus } from "react-dom";
import { Button } from "../ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "../ui/card";
import { FloatingInput, FloatingLabel } from "../ui/floating-label-input";
import { Spinner } from "../ui/spinner";

const SubmitButton = () => {
	const { pending } = useFormStatus();
	return (
		<footer className="mt-4 flex justify-end">
			<Button disabled={pending}>
				{pending && <Spinner className="size-4" />} Sign in
			</Button>
		</footer>
	);
};

export const SignIn = () => {
	const onSubmit = (data: FormData) => {
		const username = String(data.get("username"));
		const password = String(data.get("password"));
		signIn(
			"credentials",
			{
				callbackUrl: "/",
			},
			{
				username,
				password,
				endpoint: "signin",
			},
		);
	};

	return (
		<Card>
			<CardHeader>
				<CardTitle>Sign in</CardTitle>
				<CardDescription className="space-y-2">
					<p>Welcome back! Log in to your account.</p>
					<p>Skip this by using accounts listed on the right.</p>
				</CardDescription>
			</CardHeader>
			<CardContent>
				<form
					className="space-y-3"
					onSubmit={(e) => {
						e.preventDefault();
						onSubmit(new FormData(e.currentTarget));
					}}
				>
					<div className="space-y-1 relative">
						<FloatingInput
							id="username"
							defaultValue=""
							name="username"
							required
						/>
						<FloatingLabel htmlFor="username">Username</FloatingLabel>
					</div>
					<div className="space-y-1 relative">
						<FloatingInput
							id="password"
							defaultValue=""
							name="password"
							type="password"
							required
						/>
						<FloatingLabel htmlFor="password">Password</FloatingLabel>
					</div>
					<SubmitButton />
				</form>
			</CardContent>
		</Card>
	);
};
