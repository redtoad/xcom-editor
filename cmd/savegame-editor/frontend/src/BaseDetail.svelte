<script lang="ts">
  import { onMount, onDestroy, createEventDispatcher } from 'svelte';
  import { getBase, updateBase, adjustInventory, baseBitURL } from './lib/api';
  import { registerForm, ChangeTracker } from './lib/pending';
  import type { BaseDetail, InventoryItem } from './lib/types';

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
  // show items not in stores too, so they can be added
  let showAllItems = false;

  // the quantity is stored as int16
  const MAX_ITEMS = 32767;

  $: items = (base?.items ?? [])
    .filter((it) => showAllItems || it.count > 0)
    .sort((a, b) => a.name.localeCompare(b.name));

  async function load() {
    try {
      base = await getBase(gameSlot, baseIdx);
      engineers = base.engineers;
      scientists = base.scientists;
      tracker.reset({ engineers, scientists });
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);

  // Tracks which values the server has, so only changes are sent.
  const tracker = new ChangeTracker((v: { engineers: number; scientists: number }) =>
    updateBase(gameSlot, baseIdx, v));

  // apply sends the form values to the server if they changed since they
  // were last loaded or applied.
  async function apply() {
    if (!base) return;
    await tracker.apply(() => ({ engineers, scientists }));
  }

  // Let the global "Save" button apply changes not yet applied here, and
  // "Reload" refresh the form.
  onDestroy(registerForm({ apply, refresh: load }));

  // adjust changes the quantity of an item, by 10 with the shift key.
  async function adjust(item: InventoryItem, sign: number, e: MouseEvent) {
    if (!base) return;
    try {
      item.count = await adjustInventory(gameSlot, baseIdx, item.index, sign * (e.shiftKey ? 10 : 1));
      base = base;
      error = '';
    } catch (err: any) {
      error = err.message;
    }
  }

  async function save() {
    saving = true;
    try {
      await apply();
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

  <div class="inventory-header">
    <h4>Inventory</h4>
    <label class="show-all">
      <input type="checkbox" bind:checked={showAllItems} />
      Show all items
    </label>
  </div>
  {#if items.length === 0}
    <div class="meta">No items in stores.</div>
  {:else}
    <div class="inventory">
      {#each items as item (item.index)}
        <div class="inv-item" class:none={item.count === 0}>
          <span class="inv-name">{item.name}</span>
          <span class="inv-edit">
            <button class="inv-btn" title="Remove one (shift: ten)" disabled={item.count <= 0}
              on:click={(e) => adjust(item, -1, e)}>-</button>
            <span class="inv-count">{item.count}</span>
            <button class="inv-btn" title="Add one (shift: ten)" disabled={item.count >= MAX_ITEMS}
              on:click={(e) => adjust(item, 1, e)}>+</button>
          </span>
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
  .fields label {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--row-2);
    padding: 1px 2px 1px 8px;
  }
  .fields label:nth-child(2) {
    background: var(--row-3);
  }
  .fields label span {
    color: var(--text);
    text-transform: uppercase;
  }
  .fields input {
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
  .inv-item.none .inv-name {
    color: var(--muted);
  }
  .inv-edit {
    display: flex;
    align-items: center;
    gap: 4px;
  }
  .inv-count {
    color: var(--value);
    min-width: 5ch;
    text-align: right;
  }
  .inv-btn {
    min-width: 0;
    padding: 0 6px;
    line-height: 1.2;
  }
  .inventory-header {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  /* like the controls of the map */
  .show-all {
    display: flex;
    align-items: center;
    gap: 6px;
    color: var(--muted);
  }
</style>
