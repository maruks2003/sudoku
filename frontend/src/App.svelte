<script lang="ts">
    import { onMount } from "svelte";
    import Board from "./lib/Board.svelte";
    import { ensureSession } from "./lib/api/client";
    import {
        addCellNote,
        clearCellValue,
        getGame,
        removeCellNote,
        setCellValue,
        type Game,
    } from "./lib/api/game";

    // The ID the backend's stub handler uses; replaced once there is a game list.
    const gameId = "11111111-1111-4111-8111-111111111111";

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
            return getGame(gameId);
        });
    });

    function handleDigit(index: number, digit: number) {
        const writingNote = noteMode; // the mode at the moment of the keypress
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
