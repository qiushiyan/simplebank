"use client";

import { signIn } from "next-auth/react";
import { useFormState, useFormStatus } from "react-dom";
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
				{pending && <Spinner className="size-4" />} Sign up
			</Button>
		</footer>
	);
};

export const SignUp = () => {
	const onSubmit = (data: FormData) => {
		const username = String(data.get("username"));
		const password = String(data.get("password"));
		const email = String(data.get("email"));

		signIn(
			"credentials",
			{
				callbackUrl: "/",
			},
			{
				username,
				password,
				email,
				endpoint: "signup",
			},
		);
	};

	return (
		<Card>
			<CardHeader>
				<CardTitle>Sign up</CardTitle>
				<CardDescription className="space-y-2">
					<p>Create a new account to get started.</p>
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
						<FloatingInput id="username" name="username" required />
						<FloatingLabel htmlFor="username">Username</FloatingLabel>
					</div>
					<div className="space-y-1 relative">
						<FloatingInput id="email" name="email" type="email" />
						<FloatingLabel htmlFor="email">Email (Optional)</FloatingLabel>
					</div>
					<div className="space-y-1 relative">
						<FloatingInput
							id="password"
							name="password"
							type="password"
							required
						/>
						<FloatingLabel htmlFor="password">Password</FloatingLabel>
					</div>
					<div className="space-y-1 relative">
						<FloatingInput
							name="confirm-password"
							type="password"
							required
							onChange={(e) => {
								const field = e.currentTarget;
								if (field.form) {
									const password = field.form.password;
									if (field.value !== password.value) {
										field.setCustomValidity("Passwords do not match");
									} else {
										field.setCustomValidity("");
									}
								}
							}}
						/>
						<FloatingLabel htmlFor="confirm-password">
							Confirm password
						</FloatingLabel>
					</div>
					<SubmitButton />
				</form>
			</CardContent>
		</Card>
	);
};
