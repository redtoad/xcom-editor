<script lang="ts">
  import { onMount } from 'svelte';
  import { listResearch, speedupResearchProject } from './lib/api';
  import type { ResearchProject } from './lib/types';

  export let gameSlot: string;

  const SPEEDUP_HINT = 'Reduce the remaining time to one scientist day, ' +
    'so the research is completed on the next day.';
  const READY_HINT = 'The research is completed on the next day, nothing to speed up.';
  const IDLE_HINT = 'No scientists assigned: the research does not progress. ' +
    'Assign scientists in the game.';

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
          <td class="name">{p.name}</td>
          <td>{p.baseName}</td>
          <td class="num" class:idle={p.scientists === 0} title={p.scientists === 0 ? IDLE_HINT : undefined}>{p.scientists}</td>
          <td class="num">{p.daysRemaining}</td>
          <td class="actions">
            {#if p.daysRemaining <= 1}
              <span class="ready" title={READY_HINT}>Ready</span>
            {:else}
              <button on:click={() => speedup(p)} title={SPEEDUP_HINT}>Speed Up</button>
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
  <p class="hint">Days left are scientist days: the research progresses by the number of assigned scientists each day.</p>
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
  .hint { color: var(--muted); }
  .empty {
    color: var(--muted);
    text-align: center;
    padding: 40px;
    text-transform: uppercase;
  }
</style>
