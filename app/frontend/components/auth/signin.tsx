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
				<form className="space-y-2">
					<div className="space-y-1">
						<Label htmlFor="username">Username</Label>
						<Input id="username" defaultValue="" />
					</div>
					<div className="space-y-1">
						<Label htmlFor="password">Password</Label>
						<Input id="password" defaultValue="" type="password" />
					</div>
					<footer className="mt-4 flex justify-end">
						<Button>Sign in</Button>
					</footer>
				</form>
			</CardContent>
		</Card>
	);
};
