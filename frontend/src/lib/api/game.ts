import { client } from "./client";
import type { components } from "./schema";

export type Game = components["schemas"]["Game"];

function cellPosition(index: number) {
    return { row: Math.floor(index / 9), col: index % 9 };
}

function unwrap(result: { data?: Game; error?: unknown }): Game {
    if (result.data === undefined) {
        throw new Error(`request failed: ${JSON.stringify(result.error)}`);
    }
    return result.data;
}

export async function getGame(gameId: string): Promise<Game> {
    return unwrap(
        await client.GET("/games/{gameId}", { params: { path: { gameId } } }),
    );
}


export async function setCellValue(gameId: string, index: number, number: number): Promise<Game> {
    return unwrap(
        await client.PUT("/games/{gameId}/cells/{row}/{col}/value", {
            params: { path: { gameId, ...cellPosition(index) } },
            body: { number },
        }),
    );
}

export async function clearCellValue(gameId: string, index: number): Promise<Game> {
    return unwrap(
        await client.DELETE("/games/{gameId}/cells/{row}/{col}/value", {
            params: { path: { gameId, ...cellPosition(index) } },
        }),
    );
}

export async function addCellNote(gameId: string, index: number, number: number): Promise<Game> {
    return unwrap(
        await client.PUT("/games/{gameId}/cells/{row}/{col}/notes/{number}", {
            params: { path: { gameId, ...cellPosition(index), number } },
        }),
    );
}

export async function removeCellNote(gameId: string, index: number, number: number): Promise<Game> {
    return unwrap(
        await client.DELETE("/games/{gameId}/cells/{row}/{col}/notes/{number}", {
            params: { path: { gameId, ...cellPosition(index), number } },
        }),
    );
}
