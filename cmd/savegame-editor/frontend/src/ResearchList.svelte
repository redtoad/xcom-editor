<script lang="ts">
  import { onMount } from 'svelte';
  import { listResearch, speedupResearchProject } from './lib/api';
  import type { ResearchProject } from './lib/types';

  export let gameSlot: string;

  const SPEEDUP_HINT = 'Reduce the remaining time to one scientist day, ' +
    'so the research is completed on the next day.';

  let projects: ResearchProject[] = [];
  let error = '';

  async function load() {
    try {
      projects = await listResearch(gameSlot);
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  async function speedup(p: ResearchProject) {
    try {
      await speedupResearchProject(gameSlot, p.base, p.topic);
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
  <div class="empty">No ongoing research.</div>
{:else}
  <table>
    <thead>
      <tr>
        <th>Research Project</th>
        <th>Base</th>
        <th>Scientists</th>
        <th>Days Left</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      {#each projects as p (`${p.base}/${p.topic}`)}
        <tr>
          <td>{p.name}</td>
          <td>{p.baseName}</td>
          <td class="num">{p.scientists}</td>
          <td class="num">{p.daysRemaining}</td>
          <td class="actions">
            <button on:click={() => speedup(p)} disabled={p.daysRemaining <= 1} title={SPEEDUP_HINT}>Speed Up</button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
  <p class="hint">Days left are scientist days: the research progresses by the number of assigned scientists each day.</p>
{/if}

<style>
  .actions { text-align: right; }
  .hint { color: var(--muted); }
  .empty {
    color: var(--muted);
    text-align: center;
    padding: 40px;
    text-transform: uppercase;
  }
</style>
