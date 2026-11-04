<script lang="ts">
  import { onMount } from 'svelte';
  import { listProduction, speedupProductionProject } from './lib/api';
  import type { ProductionProject } from './lib/types';

  export let gameSlot: string;

  const SPEEDUP_HINT = 'Finish the next unit within the next hour. ' +
    'The remaining units are built at the normal pace.';

  let projects: ProductionProject[] = [];
  let error = '';

  async function load() {
    try {
      projects = await listProduction(gameSlot);
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  async function speedup(p: ProductionProject) {
    try {
      await speedupProductionProject(gameSlot, p.base, p.product);
      await load();
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

{#if projects.length === 0 && !error}
  <div class="empty">No ongoing production.</div>
{:else}
  <table>
    <thead>
      <tr>
        <th>Item</th>
        <th>Base</th>
        <th>Engineers</th>
        <th>Produced</th>
        <th>Total</th>
        <th>Next Unit</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      {#each projects as p (`${p.base}/${p.product}`)}
        <tr>
          <td>{p.name}</td>
          <td>{p.baseName}</td>
          <td class="num">{p.engineers}</td>
          <td class="num">{p.produced}</td>
          <td class="num">{p.total}</td>
          <td class="num" title="engineer hours done / needed per unit">{p.hoursDone}/{p.hoursPerUnit}h</td>
          <td class="actions">
            <button on:click={() => speedup(p)} disabled={p.hoursDone >= p.hoursPerUnit - 1} title={SPEEDUP_HINT}>Speed Up</button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style>
  .actions { text-align: right; }
  .empty {
    color: var(--muted);
    text-align: center;
    padding: 40px;
    text-transform: uppercase;
  }
</style>
