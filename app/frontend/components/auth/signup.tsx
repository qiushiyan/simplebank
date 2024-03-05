"use client";

import { Label } from "@radix-ui/react-label";
import { Button } from "../ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "../ui/card";
import { FloatingInput, FloatingLabel } from "../ui/floating-label-input";
import { Input } from "../ui/input";

export const SignUp = () => (
	<Card>
		<CardHeader>
			<CardTitle>Sign up</CardTitle>
			<CardDescription className="space-y-2">
				<p>Create a new account to get started.</p>
				<p>Skip this by using accounts listed on the right.</p>
			</CardDescription>
		</CardHeader>
		<CardContent>
			<form className="space-y-3">
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
				<footer className="mt-4 flex justify-end">
					<Button>Sign up</Button>
				</footer>
			</form>
		</CardContent>
	</Card>
);
