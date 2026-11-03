<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { listGames } from './lib/api';
  import type { GameSummary } from './lib/types';

  export let selectedSlot: string | null;

  const dispatch = createEventDispatcher<{ select: string }>();

  let games: GameSummary[] = [];
  let error = '';

  onMount(async () => {
    try {
      games = await listGames();
    } catch (e: any) {
      error = e.message;
    }
  });
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#each games as game}
  <button
    class="game-card"
    class:selected={selectedSlot === game.slot}
    on:click={() => dispatch('select', game.slot)}
  >
    <div class="slot">
      {game.slot.replace('_', ' ')}
      {#if game.mission}
        <span class="mission-badge" title="Saved during a mission (turn {game.mission.turn})">In mission</span>
      {/if}
    </div>
    <div class="title">{game.title}</div>
    <div class="meta">
      <span>{game.time}</span>
    </div>
    <div class="counts">
      {game.soldierCount} soldiers &middot;
      {game.baseCount} bases &middot;
      {game.craftCount} craft
    </div>
  </button>
{/each}

<style>
  /* savegame slots look like the buttons of the Geoscape */
  .game-card {
    display: block;
    width: 100%;
    text-align: left;
    text-transform: none;
    color: var(--text);
    background: var(--row-4);
    border-color: var(--row-1) var(--bg) var(--bg) var(--row-1);
    padding: 4px 8px;
    margin-bottom: 6px;
  }
  .game-card:hover {
    background: var(--row-3);
  }
  .game-card.selected {
    background: var(--button);
    color: var(--button-text);
    border-color: var(--button-light) var(--button-dark) var(--button-dark) var(--button-light);
  }
  .slot {
    color: var(--muted);
    text-transform: uppercase;
  }
  .selected .slot,
  .selected .meta,
  .selected .counts {
    color: var(--button-dark);
  }
  .mission-badge {
    float: right;
    color: var(--danger);
    background: var(--bg);
    border: 1px solid var(--danger);
    padding: 0 4px;
    text-transform: uppercase;
  }
  .title {
    color: var(--value);
  }
  .selected .title {
    color: var(--button-text);
  }
  .title,
  .meta,
  .counts {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .meta,
  .counts {
    color: var(--muted);
  }
  .error {
    color: var(--danger);
  }
</style>
