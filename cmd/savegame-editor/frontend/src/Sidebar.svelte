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
  .game-card {
    display: block;
    width: 100%;
    text-align: left;
    background: #1a1a2e;
    border: 1px solid #0f3460;
    border-radius: 6px;
    padding: 10px 12px;
    margin-bottom: 8px;
    cursor: pointer;
    color: #e0e0e0;
    font-family: inherit;
    font-size: 13px;
    transition: border-color 0.15s;
  }
  .game-card:hover {
    border-color: #00ff41;
  }
  .game-card.selected {
    border-color: #00ff41;
    background: #0f3460;
  }
  .slot {
    font-size: 11px;
    color: #888;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  .mission-badge {
    float: right;
    background: #5c0a0a;
    border: 1px solid #ff4444;
    color: #ff8888;
    border-radius: 3px;
    padding: 0 5px;
    font-size: 10px;
    letter-spacing: 0.5px;
  }
  .title {
    font-weight: bold;
    color: #00ff41;
    margin: 2px 0;
  }
  .meta {
    font-size: 12px;
    color: #888;
  }
  .counts {
    font-size: 11px;
    color: #666;
    margin-top: 4px;
  }
  .error {
    color: #ff4444;
    padding: 8px;
    font-size: 13px;
  }
</style>
