<script lang="ts">
  import { onMount } from 'svelte';
  import { listCraft } from './lib/api';
  import type { CraftSummary } from './lib/types';

  export let gameSlot: string;

  const DAMAGED_HINT = 'The entry of this craft in CRAFT.DAT is empty (e.g. overwritten with zeros), ' +
    'so its values are meaningless. The craft still exists in the game (LOC.DAT refers to it).';

  let craft: CraftSummary[] = [];
  let error = '';

  async function load() {
    try {
      craft = await listCraft(gameSlot);
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);
  $: if (gameSlot) load();
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#if craft.length === 0 && !error}
  <div class="empty">No craft found.</div>
{:else}
  <table>
    <thead>
      <tr>
        <th>Name</th>
        <th>Type</th>
        <th>Base</th>
        <th>Status</th>
        <th>Damage</th>
        <th>Fuel</th>
      </tr>
    </thead>
    <tbody>
      {#each craft as c}
        <tr class:damaged={c.damaged}>
          <td class="name">
            {c.name}
            {#if c.damaged}
              <span class="badge damaged-badge" title={DAMAGED_HINT}>Damaged</span>
            {/if}
          </td>
          <td>{c.type}</td>
          <td>{c.baseName}</td>
          <td>
            <span class="badge" class:ready={c.status === 'Ready'}
              class:out={c.status === 'Out'}>{c.status}</span>
          </td>
          <td class="num">{c.damage}</td>
          <td class="num">{c.fuel}</td>
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
  .name { font-weight: bold; }
  .num { text-align: right; font-variant-numeric: tabular-nums; }
  .badge {
    font-size: 11px;
    padding: 2px 8px;
    border-radius: 3px;
    background: #0f3460;
  }
  .badge.ready { background: #0a3a0a; color: #44ff44; }
  .badge.out { background: #3a3a0a; color: #ffff44; }
  tr.damaged td { color: #666; }
  .badge.damaged-badge {
    margin-left: 6px;
    background: #5c0a0a;
    color: #ff8888;
    font-weight: normal;
    cursor: help;
  }
  .empty {
    color: #666;
    text-align: center;
    padding: 40px;
  }
  .error {
    color: #ff4444;
    padding: 8px;
    margin-bottom: 8px;
    background: rgba(255,68,68,0.1);
    border-radius: 4px;
  }
</style>
