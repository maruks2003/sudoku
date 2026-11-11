import createClient from "openapi-fetch";
import type { paths } from "./schema";

export const client = createClient<paths>({ baseUrl: "" });

const TOKEN_KEY = "token";
const RENEWED_AT_KEY = "sessionRenewedAt";

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

// On a 401, replace the session and reload the page.
let renewal: Promise<void> | null = null;

client.use({
    async onResponse({ request, response }) {
        const isSessionRequest = new URL(request.url).pathname === "/session";
        if (response.status !== 401 || isSessionRequest) return;

        if (renewal === null) {
            if (renewedRecently()) return; // a fresh token was rejected too: don't loop
            renewal = renewSession();
        }
        await renewal;
        location.reload();

        // Keep the original call pending so the UI doesn't flash an error
        // while the page reloads.
        await new Promise<never>(() => { });
    },
});

function renewedRecently(): boolean {
    const renewedAt = Number(sessionStorage.getItem(RENEWED_AT_KEY));
    return Date.now() - renewedAt < 10_000;
}

async function renewSession(): Promise<void> {
    localStorage.removeItem(TOKEN_KEY);
    sessionStorage.setItem(RENEWED_AT_KEY, String(Date.now()));
    await ensureSession();
}

// Create an anonymous player on first visit.
export async function ensureSession(): Promise<void> {
    if (localStorage.getItem("token")) return;
    const { data, error } = await client.POST("/session");
    if (error || !data) throw new Error("Could not create session");
    localStorage.setItem("token", data.token);
}
