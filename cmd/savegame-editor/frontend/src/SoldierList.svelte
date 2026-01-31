<script lang="ts">
  import { onMount } from 'svelte';
  import { listSoldiers } from './lib/api';
  import type { SoldierSummary } from './lib/types';
  import SoldierEdit from './SoldierEdit.svelte';

  export let gameSlot: string;

  let soldiers: SoldierSummary[] = [];
  let selectedIdx: number | null = null;
  let error = '';

  async function load() {
    try {
      soldiers = await listSoldiers(gameSlot);
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);
  $: if (gameSlot) { selectedIdx = null; load(); }

  function handleBack() {
    selectedIdx = null;
    load();
  }
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#if selectedIdx !== null}
  <button class="back" on:click={handleBack}>&larr; Back to list</button>
  <SoldierEdit {gameSlot} soldierIdx={selectedIdx} on:saved={handleBack} />
{:else}
  <table>
    <thead>
      <tr>
        <th>Name</th>
        <th>Rank</th>
        <th>Base</th>
        <th>Craft</th>
        <th>Missions</th>
        <th>Kills</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody>
      {#each soldiers as s}
        <tr
          class:dead={s.isDead}
          class:wounded={s.isWounded}
          on:click={() => selectedIdx = s.index}
        >
          <td class="name">{s.name}</td>
          <td>{s.rank}</td>
          <td>{s.baseName}</td>
          <td>{s.craftName}</td>
          <td class="num">{s.missions}</td>
          <td class="num">{s.kills}</td>
          <td>
            {#if s.isDead}<span class="badge dead">DEAD</span>
            {:else if s.isWounded}<span class="badge wounded">WOUNDED</span>
            {:else}<span class="badge ok">OK</span>
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style>
  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 13px;
  }
  th {
    text-align: left;
    padding: 6px 10px;
    border-bottom: 1px solid #0f3460;
    color: #888;
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  td {
    padding: 6px 10px;
    border-bottom: 1px solid #16213e;
  }
  tr { cursor: pointer; transition: background 0.1s; }
  tr:hover { background: #16213e; }
  tr.dead { opacity: 0.5; }
  tr.wounded td { color: #ffaa44; }
  .name { font-weight: bold; }
  .num { text-align: right; font-variant-numeric: tabular-nums; }
  .badge {
    font-size: 10px;
    padding: 2px 6px;
    border-radius: 3px;
    font-weight: bold;
  }
  .badge.dead { background: #5c0a0a; color: #ff4444; }
  .badge.wounded { background: #5c3a0a; color: #ffaa44; }
  .badge.ok { background: #0a3a0a; color: #44ff44; }
  .back {
    background: none;
    border: none;
    color: #00ff41;
    cursor: pointer;
    font-family: inherit;
    font-size: 13px;
    padding: 4px 0;
    margin-bottom: 12px;
  }
  .back:hover { text-decoration: underline; }
  .error {
    color: #ff4444;
    padding: 8px;
    margin-bottom: 8px;
    background: rgba(255,68,68,0.1);
    border-radius: 4px;
  }
</style>
