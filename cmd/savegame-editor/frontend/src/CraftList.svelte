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
  .badge { color: var(--muted); }
  .badge.ready { color: var(--ok); }
  .badge.out { color: var(--value); }
  tr.damaged td { color: var(--dim); }
  tr.damaged td.num { color: var(--dim); }
  .badge.damaged-badge {
    margin-left: 6px;
    color: var(--danger);
    cursor: help;
  }
  .empty {
    color: var(--muted);
    text-align: center;
    padding: 40px;
    text-transform: uppercase;
  }
</style>
