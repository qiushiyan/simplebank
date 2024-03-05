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
import { Input } from "../ui/input";

export const SignUp = () => {
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
				<form className="space-y-2">
					<div className="space-y-1">
						<Label htmlFor="username">Username</Label>
						<Input id="username" name="username" required />
					</div>
					<div className="space-y-1">
						<Label htmlFor="email">Email (Optional)</Label>
						<Input id="email" name="email" type="email" />
					</div>
					<div className="space-y-1">
						<Label htmlFor="password">Password</Label>
						<Input id="password" name="password" type="password" required />
					</div>
					<div className="space-y-1">
						<Label htmlFor="confirm-password">Confirm password</Label>
						<Input
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
					</div>
					<footer className="mt-4 flex justify-end">
						<Button>Sign up</Button>
					</footer>
				</form>
			</CardContent>
		</Card>
	);
};
