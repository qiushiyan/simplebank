"use client";

import { Button } from "../ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "../ui/card";
import { FloatingInput, FloatingLabel } from "../ui/floating-label-input";

export const SignIn = () => {
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
				<form className="space-y-3">
					<div className="space-y-1 relative">
						<FloatingInput id="username" defaultValue="" required />
						<FloatingLabel htmlFor="username">Username</FloatingLabel>
					</div>
					<div className="space-y-1 relative">
						<FloatingInput
							id="password"
							defaultValue=""
							type="password"
							required
						/>
						<FloatingLabel htmlFor="password">Password</FloatingLabel>
					</div>
					<footer className="mt-4 flex justify-end">
						<Button>Sign in</Button>
					</footer>
				</form>
			</CardContent>
		</Card>
	);
};
