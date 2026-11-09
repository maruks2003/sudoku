import createClient from "openapi-fetch";
import type { paths } from "./schema";

export const client = createClient<paths>({ baseUrl: "" });

// Attach the bearer token to every request.
client.use({
    onRequest({ request }) {
        const token = localStorage.getItem("token");
        if (token) {
            request.headers.set("Authorization", `Bearer ${token}`);
        }
        return request;
    }
})

// Create an anonymous player on first visit.
export async function ensureSession(): Promise<void> {
    if (localStorage.getItem("token")) return;
    const { data, error } = await client.POST("/session");
    if (error || !data) throw new Error("Could not create session");
    localStorage.setItem("token", data.token);
}
