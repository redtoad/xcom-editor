<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { getBase, updateBase } from './lib/api';
  import type { BaseDetail } from './lib/types';

  export let gameSlot: string;
  export let baseIdx: number;

  const dispatch = createEventDispatcher();

  let base: BaseDetail | null = null;
  let error = '';
  let saving = false;
  let engineers = 0;
  let scientists = 0;

  async function load() {
    try {
      base = await getBase(gameSlot, baseIdx);
      engineers = base.engineers;
      scientists = base.scientists;
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);

  async function save() {
    saving = true;
    try {
      await updateBase(gameSlot, baseIdx, { engineers, scientists });
      dispatch('saved');
    } catch (e: any) {
      error = e.message;
    }
    saving = false;
  }
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#if base}
  <h3>{base.name}</h3>
  <div class="meta">{base.coord}</div>

  <div class="fields">
    <label>
      <span>Engineers</span>
      <input type="number" bind:value={engineers} min="0" max="255" />
    </label>
    <label>
      <span>Scientists</span>
      <input type="number" bind:value={scientists} min="0" max="255" />
    </label>
  </div>

  <h4>Facility Grid (6x6)</h4>
  <div class="grid">
    {#each { length: 6 } as _, y}
      <div class="grid-row">
        {#each { length: 6 } as _, x}
          {@const tile = base.tiles[x + y * 6]}
          <div
            class="tile"
            class:empty={tile.type === 'Empty'}
            class:building={tile.daysToCompletion > 0}
            title={tile.type + (tile.daysToCompletion > 0 ? ` (${tile.daysToCompletion} days)` : '')}
          >
            <span class="tile-name">{tile.type === 'Empty' ? '' : tile.type.replace(/ *\(.*\)/, '')}</span>
            {#if tile.daysToCompletion > 0}
              <span class="tile-days">{tile.daysToCompletion}d</span>
            {/if}
          </div>
        {/each}
      </div>
    {/each}
  </div>

  {#if Object.keys(base.inventory).length > 0}
    <h4>Inventory</h4>
    <div class="inventory">
      {#each Object.entries(base.inventory).sort((a, b) => a[0].localeCompare(b[0])) as [item, count]}
        <div class="inv-item">
          <span class="inv-name">{item}</span>
          <span class="inv-count">{count}</span>
        </div>
      {/each}
    </div>
  {/if}

  <div class="actions">
    <button class="save-btn" on:click={save} disabled={saving}>
      {saving ? 'Applying...' : 'Apply Changes'}
    </button>
  </div>
{/if}

<style>
  h3 {
    color: #00ff41;
    margin: 0 0 4px 0;
  }
  .meta {
    color: #666;
    font-size: 12px;
    margin-bottom: 16px;
  }
  h4 {
    color: #aaa;
    font-size: 13px;
    margin: 20px 0 8px 0;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  .fields {
    display: flex;
    gap: 16px;
    margin-bottom: 8px;
  }
  label {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  label span {
    font-size: 11px;
    color: #888;
    text-transform: uppercase;
  }
  input, select {
    background: #16213e;
    border: 1px solid #0f3460;
    color: #e0e0e0;
    padding: 6px 8px;
    border-radius: 4px;
    font-family: inherit;
    font-size: 13px;
    width: 80px;
  }
  input:focus { outline: none; border-color: #00ff41; }
  .grid {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .grid-row {
    display: flex;
    gap: 2px;
  }
  .tile {
    width: 80px;
    height: 48px;
    background: #0f3460;
    border-radius: 3px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    font-size: 9px;
    color: #ccc;
    position: relative;
    overflow: hidden;
  }
  .tile.empty {
    background: #1a1a2e;
    border: 1px dashed #333;
  }
  .tile.building {
    border: 1px solid #ffaa44;
  }
  .tile-name {
    text-align: center;
    line-height: 1.1;
    padding: 0 2px;
  }
  .tile-days {
    color: #ffaa44;
    font-size: 8px;
    font-weight: bold;
  }
  .inventory {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 4px;
    font-size: 13px;
  }
  .inv-item {
    display: flex;
    justify-content: space-between;
    padding: 3px 8px;
    background: #16213e;
    border-radius: 3px;
  }
  .inv-name { color: #ccc; }
  .inv-count { color: #00ff41; font-weight: bold; }
  .actions {
    margin-top: 20px;
  }
  .save-btn {
    background: #0a5c0a;
    border: 1px solid #0f8f0f;
    color: #e0e0e0;
    padding: 8px 20px;
    border-radius: 4px;
    cursor: pointer;
    font-family: inherit;
    font-size: 14px;
  }
  .save-btn:hover { background: #0f8f0f; }
  .save-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .error {
    color: #ff4444;
    padding: 8px;
    margin-bottom: 8px;
    background: rgba(255,68,68,0.1);
    border-radius: 4px;
  }
</style>
