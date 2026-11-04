<script lang="ts">
  import { onMount } from 'svelte';
  import { listSoldiers, resurrectSoldier } from './lib/api';
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

  async function resurrect(idx: number) {
    try {
      await resurrectSoldier(gameSlot, idx);
      await load();
    } catch (e: any) {
      error = e.message;
    }
  }

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
          {#if s.isDead}
            <td title="Estimated rank before being killed in action">{s.lastRank}?</td>
          {:else}
            <td>{s.rank}</td>
          {/if}
          <td>{s.baseName}</td>
          <td>{s.craftName}</td>
          <td class="num">{s.missions}</td>
          <td class="num">{s.kills}</td>
          <td>
            {#if s.isDead}
              <span class="badge dead">DEAD</span>
              <button class="resurrect" title="Comes back as {s.lastRank}" on:click|stopPropagation={() => resurrect(s.index)}>Resurrect</button>
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
  tbody tr { cursor: pointer; }
  tr.dead td { color: var(--dim); }
  tr.wounded td:not(.num) { color: var(--warning); }
  .badge.dead { color: var(--danger); }
  .badge.wounded { color: var(--warning); }
  .badge.ok { color: var(--ok); }
  .back { margin-bottom: 12px; }
  .resurrect { margin-left: 8px; }
</style>
