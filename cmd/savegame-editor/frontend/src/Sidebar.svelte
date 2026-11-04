<script lang="ts">
  import { createEventDispatcher, onDestroy } from 'svelte';
  import { listGames } from './lib/api';
  import type { GameSummary } from './lib/types';

  export let selectedSlot: string | null;
  // incremented by the parent whenever a savegame has changed (e.g. after
  // a reload or save), which reloads the list of savegames
  export let refresh = 0;

  // how often the list of savegames is checked for new and changed savegames
  const pollInterval = 3000;

  const dispatch = createEventDispatcher<{ select: string; update: GameSummary[] }>();

  let games: GameSummary[] = [];
  let error = '';

  async function load() {
    try {
      games = await listGames();
      error = '';
      dispatch('update', games);
    } catch (e: any) {
      error = e.message;
    }
  }

  $: refresh, load();

  const timer = setInterval(load, pollInterval);
  onDestroy(() => clearInterval(timer));
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
      <span class="badges">
        {#if game.stale}
          <span class="badge stale-badge" title="The savegame was changed on disk (e.g. saved in the game) and needs to be reloaded">Changed</span>
        {/if}
        {#if game.mission}
          <span class="badge mission-badge" title="Saved during a mission (turn {game.mission.turn})">In mission</span>
        {/if}
      </span>
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
  .badges {
    float: right;
    display: inline-flex;
    gap: 4px;
  }
  .badge {
    background: var(--bg);
    border: 1px solid currentColor;
    padding: 0 4px;
    text-transform: uppercase;
  }
  .mission-badge {
    color: var(--danger);
  }
  .stale-badge {
    color: var(--warning);
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
