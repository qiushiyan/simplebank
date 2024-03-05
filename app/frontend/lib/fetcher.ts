import { Session } from "next-auth";
import { env } from "./env.mjs";

export const fetcher = async (
	endpoint: string,
	method: "GET" | "POST",
	user: Session["user"],
	body?: any,
) => {
	const response = await fetch(`${env.API_URL}/${endpoint}`, {
		method,
		headers: {
			"Content-Type": "application/json",
			Authorization: `Bearer ${user.access_token}`,
		},
		body: body ? JSON.stringify(body) : undefined,
	});

	if (!response.ok) {
		return null;
	}

	return response.json();
};
