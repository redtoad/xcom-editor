<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { getBase, updateBase, baseBitURL } from './lib/api';
  import type { BaseDetail } from './lib/types';

  export let gameSlot: string;
  export let baseIdx: number;

  // Base sprites are 32x32 pixels (some are taller), drawn with doubled
  // pixels like on the original 320x200 screen.
  const SCALE = 2;
  const TILE = 32 * SCALE;

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
  <div class="base-screen">
    <div class="base-layout">
      {#if base.layers}
        <div class="base-view" style="width:{6 * TILE}px; height:{6 * TILE}px">
          {#each base.layers as layer}
            {#each layer as s}
              <img
                class="sprite"
                src={baseBitURL(s.sprite)}
                alt=""
                style="left:{s.x * SCALE}px; top:{s.y * SCALE}px; width:{TILE}px"
              />
            {/each}
          {/each}
          {#each base.tiles as tile, i}
            {#if tile.type !== 'Empty'}
              <div
                class="base-tile"
                style="left:{(i % 6) * TILE}px; top:{Math.floor(i / 6) * TILE}px; width:{TILE}px; height:{TILE}px"
                title={tile.type + (tile.daysToCompletion > 0 ? ` (${tile.daysToCompletion} days)` : '')}
              >
                {#if tile.daysToCompletion > 0}
                  <span class="base-days">{tile.daysToCompletion}d</span>
                {/if}
              </div>
            {/if}
          {/each}
        </div>
      {:else}
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
      {/if}
    </div>

    <!-- side panel like the base screen of the game -->
    <div class="side">
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

      <button class="save-btn" on:click={save} disabled={saving}>
        {saving ? 'Applying...' : 'Apply Changes'}
      </button>
    </div>
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
{/if}

<style>
  .base-screen {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    align-items: flex-start;
  }
  .side {
    flex: 1;
    min-width: 260px;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .side h3 {
    margin: 0;
  }
  .meta {
    color: var(--muted);
  }
  .fields {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  label {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--row-2);
    padding: 1px 2px 1px 8px;
  }
  label:nth-child(2) {
    background: var(--row-3);
  }
  label span {
    color: var(--text);
    text-transform: uppercase;
  }
  input {
    width: 80px;
    text-align: right;
  }
  .save-btn {
    align-self: stretch;
  }
  .base-view {
    position: relative;
    /* some sprites are taller than a tile */
    margin-bottom: 16px;
  }
  .sprite {
    position: absolute;
    height: auto;
    image-rendering: pixelated;
    pointer-events: none;
  }
  .base-tile {
    position: absolute;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .base-tile:hover {
    outline: 2px solid var(--value);
    outline-offset: -2px;
  }
  .base-days {
    color: var(--value);
    background: var(--bg);
    padding: 0 2px;
  }
  /* simple grid without the game graphics */
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
    width: 64px;
    height: 64px;
    background: var(--row-3);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    font-family: 'XcomSmall', monospace;
    font-size: 9px;
    line-height: 1.1;
    text-align: center;
    color: var(--text);
    overflow: hidden;
  }
  .tile.empty {
    background: var(--panel);
  }
  .tile.building {
    border: 2px dotted var(--value);
  }
  .tile-days {
    color: var(--value);
  }
  .inventory {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 2px 8px;
  }
  .inv-item {
    display: flex;
    justify-content: space-between;
    padding: 0 8px;
    background: var(--row-3);
  }
  .inv-item:nth-child(even) {
    background: var(--row-4);
  }
  .inv-count {
    color: var(--value);
  }
</style>
