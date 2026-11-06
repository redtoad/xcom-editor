<script lang="ts">
  import { onMount } from 'svelte';
  import { listProduction, speedupProductionProject } from './lib/api';
  import type { ProductionProject } from './lib/types';

  export let gameSlot: string;

  const SPEEDUP_HINT = 'Finish the next unit within the next hour. ' +
    'The remaining units are built at the normal pace.';
  const READY_HINT = 'The next unit is finished within the next hour, ' +
    'nothing to speed up (new projects start like this).';
  const IDLE_HINT = 'No engineers allocated: the production does not progress. ' +
    'Allocate engineers in the game.';

  // the next unit is finished within the next hour (see SpeedUp in BPROD.DAT)
  const ready = (p: ProductionProject) => p.hoursDone >= p.hoursPerUnit - 1;

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
        <th>Next Unit</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      {#each projects as p (`${p.base}/${p.product}`)}
        <tr>
          <td class="name">{p.name}</td>
          <td>{p.baseName}</td>
          <td class="num" class:idle={p.engineers === 0} title={p.engineers === 0 ? IDLE_HINT : undefined}>{p.engineers}</td>
          <td class="num" title="units produced / total">{p.produced}/{p.total}</td>
          <td class="num" title="engineer hours done / needed per unit">{p.hoursDone}/{p.hoursPerUnit}h</td>
          <td class="actions">
            {#if ready(p)}
              <span class="ready" title={READY_HINT}>Ready</span>
            {:else}
              <button on:click={() => speedup(p)} title={SPEEDUP_HINT}>Speed Up</button>
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style>
  .actions { text-align: right; }
  /* long names may wrap, so the table fits next to the sidebar */
  td.name { white-space: normal; }
  .ready {
    color: var(--ok);
    text-transform: uppercase;
    cursor: help;
  }
  /* nothing assigned, the project does not progress */
  td.num.idle {
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
