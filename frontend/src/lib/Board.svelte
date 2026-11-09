<script lang="ts">
    import type { Game } from "./api/game";

    let {
        game,
        noteMode,
        onDigit,
        onClear,
        onToggleNoteMode,
    }: {
        game: Game;
        noteMode: boolean;
        onDigit: (index: number, digit: number) => void;
        onClear: (index: number) => void;
        onToggleNoteMode: () => void;
    } = $props();

    const cellIndexes = Array.from({ length: 81 }, (_, index) => index);
    const digits = [1, 2, 3, 4, 5, 6, 7, 8, 9];

    let hoveredIndex = $state<number | null>(null);
    const conflicts = $derived(new Set(game.conflicts));

    function handleKeydown(event: KeyboardEvent) {
        if (event.ctrlKey || event.metaKey || event.altKey) return;

        if (event.key === "n" || event.key === "N") {
            onToggleNoteMode();
            return;
        }
        if (hoveredIndex === null) return;

        if (/^[1-9]$/.test(event.key)) {
            onDigit(hoveredIndex, Number(event.key));
        } else if (["Backspace", "Delete", "0"].includes(event.key)) {
            onClear(hoveredIndex);
        }
    }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="board" class:note-mode={noteMode}>
    {#each cellIndexes as index (index)}
        {@const row = Math.floor(index / 9)}
        {@const col = index % 9}
        {@const isGiven = game.givens[index] !== "0"}
        {@const entry = game.entries[index]}
        <button
            type="button"
            tabindex="-1"
            class="cell"
            class:given={isGiven}
            class:conflict={conflicts.has(index)}
            class:hovered={hoveredIndex === index}
            class:box-right={col % 3 === 2 && col !== 8}
            class:box-bottom={row % 3 === 2 && row !== 8}
            onmouseenter={() => (hoveredIndex = index)}
            onmouseleave={() => (hoveredIndex = null)}
        >
            {#if isGiven}
                {game.givens[index]}
            {:else if entry !== "0"}
                <span class="entry">{entry}</span>
            {:else if game.notes[index] !== 0}
                <span class="notes">
                    {#each digits as digit}
                        <span
                            >{game.notes[index] & (1 << (digit - 1))
                                ? digit
                                : ""}</span
                        >
                    {/each}
                </span>
            {/if}
        </button>
    {/each}
</div>

<style>
    .board {
        display: grid;
        grid-template-columns: repeat(9, 3rem);
        grid-template-rows: repeat(9, 3rem);
        width: fit-content;
        border: 3px solid currentColor;
    }
    .cell {
        display: flex;
        align-items: center;
        justify-content: center;
        box-sizing: border-box;
        width: 3rem;
        height: 3rem;
        padding: 0;
        border: 1px solid rgb(128 128 128 / 0.35);
        border-radius: 0;
        background: transparent;
        color: inherit;
        font-size: 1.5rem;
        cursor: default;
    }
    .cell.box-right {
        border-right: 3px solid currentColor;
    }
    .cell.box-bottom {
        border-bottom: 3px solid currentColor;
    }
    .cell.given {
        font-weight: 700;
    }
    .cell.hovered {
        background: rgb(59 130 246 / 0.2);
    }
    .board.note-mode .cell.hovered {
        background: rgb(234 179 8 / 0.25);
    }
    .entry {
        color: #2563eb;
    }
    .cell.conflict .entry {
        color: #dc2626;
    }
    .notes {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        align-items: center;
        justify-items: center;
        width: 100%;
        height: 100%;
        font-size: 0.65rem;
        line-height: 1;
        opacity: 0.75;
    }
</style>
