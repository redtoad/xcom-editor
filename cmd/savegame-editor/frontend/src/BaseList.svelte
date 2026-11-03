<script lang="ts">
  import { onMount } from 'svelte';
  import { listBases, getBase } from './lib/api';
  import type { BaseSummary, BaseDetail } from './lib/types';
  import BaseDetailView from './BaseDetail.svelte';

  export let gameSlot: string;

  let bases: BaseSummary[] = [];
  let selectedIdx: number | null = null;
  let error = '';

  async function load() {
    try {
      bases = await listBases(gameSlot);
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
  <BaseDetailView {gameSlot} baseIdx={selectedIdx} on:saved={handleBack} />
{:else}
  <div class="base-grid">
    {#each bases as base}
      <button class="base-card" on:click={() => selectedIdx = base.index}>
        <div class="base-name">{base.name}</div>
        <div class="base-coord">{base.coord}</div>
        <div class="base-staff">
          <span>Engineers: {base.engineers}</span>
          <span>Scientists: {base.scientists}</span>
        </div>
      </button>
    {/each}
  </div>
{/if}

<style>
  .back { margin-bottom: 12px; }
  /* bases are listed like the menu buttons of the base screen */
  .base-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 8px;
  }
  .base-card {
    text-align: left;
    text-transform: none;
    padding: 6px 10px;
  }
  .base-name {
    font-family: var(--font-title);
    font-size: var(--fs-title);
    color: var(--button-text);
  }
  .base-coord {
    color: var(--button-dark);
  }
  .base-staff {
    display: flex;
    gap: 16px;
    margin-top: 4px;
  }
</style>
