import { Session } from "next-auth";
import { BaseUser } from "..";
import { env } from "./env.mjs";

export const ping = async () => {
	try {
		await fetch(`${env.API_URL}/liveness`);
		return true;
	} catch (error) {
		return false;
	}
};

export const fetcher = async (
	endpoint: string,
	method: "GET" | "POST" | "PATCH",
	user?: BaseUser,
	body?: any,
) => {
	const response = await fetch(`${env.API_URL}/${endpoint}`, {
		method,
		headers: {
			"Content-Type": "application/json",
			Authorization: `Bearer ${user?.access_token}`,
		},
		body: body ? JSON.stringify(body) : undefined,
	});

	return response.json();
};
