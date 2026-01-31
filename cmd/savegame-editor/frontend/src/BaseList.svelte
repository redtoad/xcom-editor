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
  .base-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 12px;
  }
  .base-card {
    display: block;
    width: 100%;
    text-align: left;
    background: #16213e;
    border: 1px solid #0f3460;
    border-radius: 6px;
    padding: 14px;
    cursor: pointer;
    color: #e0e0e0;
    font-family: inherit;
    font-size: 13px;
    transition: border-color 0.15s;
  }
  .base-card:hover {
    border-color: #00ff41;
  }
  .base-name {
    font-weight: bold;
    color: #00ff41;
    font-size: 15px;
    margin-bottom: 4px;
  }
  .base-coord {
    color: #666;
    font-size: 11px;
    margin-bottom: 8px;
  }
  .base-staff {
    display: flex;
    gap: 16px;
    font-size: 12px;
    color: #888;
  }
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
