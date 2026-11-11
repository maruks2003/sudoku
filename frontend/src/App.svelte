<script lang="ts">
    import { onMount } from "svelte";
    import Board from "./lib/Board.svelte";
    import { ensureSession } from "./lib/api/client";
    import {
        addCellNote,
        clearCellValue,
        createGame,
        getGame,
        listGames,
        removeCellNote,
        setCellValue,
        type Game,
    } from "./lib/api/game";

    let game = $state<Game | null>(null);
    let noteMode = $state(false);
    let errorMessage = $state("");

    // Requests run one after another, so a slow response can't overwrite a newer one.
    let pending: Promise<void> = Promise.resolve();

    function run(action: () => Promise<Game>) {
        pending = pending.then(async () => {
            try {
                game = await action();
                errorMessage = "";
            } catch (caught) {
                errorMessage = String(caught);
            }
        });
    }

    onMount(() => {
        run(async () => {
            await ensureSession();
            let games = await listGames();
            let id = games.at(0)?.id;
            if (games.length > 1 && id !== undefined) {
                return getGame(id);
            }
            return createGame("easy");
        });
    });

    function handleDigit(index: number, digit: number) {
        const writingNote = noteMode; // the mode at the moment of the keypress
        if (game == null) {
            return;
        }
        const gameId = game.id;
        run(() => {
            if (!writingNote) {
                return setCellValue(gameId, index, digit);
            }
            const hasNote =
                ((game?.notes[index] ?? 0) & (1 << (digit - 1))) !== 0;
            return hasNote
                ? removeCellNote(gameId, index, digit)
                : addCellNote(gameId, index, digit);
        });
    }

    function handleClear(index: number) {
        if (game == null) {
            return;
        }
        const gameId = game.id;
        run(() => clearCellValue(gameId, index));
    }

    function toggleNoteMode() {
        noteMode = !noteMode;
    }
</script>

<main>
    <h1>Sudoku</h1>
    {#if errorMessage}
        <p>{errorMessage}</p>
    {/if}
    {#if game}
        <Board
            {game}
            {noteMode}
            onDigit={handleDigit}
            onClear={handleClear}
            onToggleNoteMode={toggleNoteMode}
        />
        <button type="button" onclick={toggleNoteMode}>
            Mode: {noteMode ? "notes" : "values"} (press N)
        </button>
    {:else if !errorMessage}
        <p>Loading…</p>
    {/if}
</main>
