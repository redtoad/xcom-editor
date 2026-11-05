<script lang="ts">
  import { onMount } from 'svelte';
  import { listBases } from './lib/api';
  import { applyPendingChanges } from './lib/pending';
  import type { BaseSummary } from './lib/types';
  import BaseDetailView from './BaseDetail.svelte';
  import MiniBases from './MiniBases.svelte';

  export let gameSlot: string;

  let bases: BaseSummary[] = [];
  let selectedIdx: number | null = null;
  let error = '';
  let loaded = false;

  async function load() {
    try {
      bases = await listBases(gameSlot);
      // like the game, start with the first base
      if (!bases.some(b => b.index === selectedIdx)) {
        selectedIdx = bases.length > 0 ? bases[0].index : null;
      }
      error = '';
      loaded = true;
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);
  $: if (gameSlot) { selectedIdx = null; loaded = false; load(); }

  // select switches to another base, like clicking a mini base on the base
  // screen. Changes not applied yet are applied first so they are not lost.
  async function select(idx: number) {
    if (idx === selectedIdx) return;
    try {
      await applyPendingChanges();
    } catch (e: any) {
      error = e.message;
      return;
    }
    selectedIdx = idx;
    load();
  }
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#if selectedIdx !== null}
  <BaseDetailView {gameSlot} baseIdx={selectedIdx} on:saved={load}>
    <MiniBases slot="selector" {bases} selected={selectedIdx} on:select={e => select(e.detail)} />
  </BaseDetailView>
{:else if loaded && bases.length === 0}
  <div class="empty">No bases.</div>
{/if}

<style>
  .empty {
    color: var(--muted);
  }
</style>
