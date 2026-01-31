<script lang="ts">
  import { onMount } from 'svelte';
  import { getGame, saveGame, reloadGame, healAll, completeConstructions, speedupDeliveries } from './lib/api';
  import type { GameDetail } from './lib/types';
  import SoldierList from './SoldierList.svelte';
  import BaseList from './BaseList.svelte';
  import CraftList from './CraftList.svelte';
  import TransferList from './TransferList.svelte';
  import FinancialView from './FinancialView.svelte';

  // Svelte reserves "slot" as a keyword, use a different prop name
  export let slot: string;

  let game: GameDetail | null = null;
  let activeTab = 'soldiers';
  let status = '';
  let error = '';

  const tabs = [
    { id: 'soldiers', label: 'Soldiers' },
    { id: 'bases', label: 'Bases' },
    { id: 'craft', label: 'Craft' },
    { id: 'transfers', label: 'Transfers' },
    { id: 'financials', label: 'Financials' },
  ];

  async function load() {
    try {
      game = await getGame(slot);
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);
  $: if (slot) load();

  async function doAction(name: string, fn: () => Promise<void>) {
    try {
      status = name + '...';
      await fn();
      status = name + ' done!';
      await load();
      setTimeout(() => status = '', 2000);
    } catch (e: any) {
      status = '';
      error = e.message;
    }
  }
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#if game}
  <div class="header">
    <h2>{game.title}</h2>
    <div class="meta">{game.time} &middot; Balance: ${game.balance?.toLocaleString()}</div>
  </div>

  <nav class="tabs">
    {#each tabs as tab}
      <button
        class="tab"
        class:active={activeTab === tab.id}
        on:click={() => activeTab = tab.id}
      >{tab.label}</button>
    {/each}
  </nav>

  <div class="tab-content">
    {#if activeTab === 'soldiers'}
      <SoldierList gameSlot={slot} />
    {:else if activeTab === 'bases'}
      <BaseList gameSlot={slot} />
    {:else if activeTab === 'craft'}
      <CraftList gameSlot={slot} />
    {:else if activeTab === 'transfers'}
      <TransferList gameSlot={slot} />
    {:else if activeTab === 'financials'}
      <FinancialView gameSlot={slot} />
    {/if}
  </div>

  <div class="action-bar">
    <button class="action" on:click={() => doAction('Healing', () => healAll(slot))}>Heal All</button>
    <button class="action" on:click={() => doAction('Completing', () => completeConstructions(slot))}>Complete Constructions</button>
    <button class="action" on:click={() => doAction('Speeding up', () => speedupDeliveries(slot))}>Speed Up Deliveries</button>
    <div class="spacer"></div>
    <button class="action reload" on:click={() => doAction('Reloading', () => reloadGame(slot))}>Reload</button>
    <button class="action save" on:click={() => doAction('Saving', () => saveGame(slot))}>Save</button>
    {#if status}
      <span class="status">{status}</span>
    {/if}
  </div>
{/if}

<style>
  .header h2 {
    margin: 0;
    color: #00ff41;
  }
  .meta {
    color: #888;
    font-size: 13px;
    margin-top: 4px;
  }
  .tabs {
    display: flex;
    gap: 0;
    margin: 16px 0 0 0;
    border-bottom: 1px solid #0f3460;
  }
  .tab {
    background: none;
    border: none;
    border-bottom: 2px solid transparent;
    color: #888;
    padding: 8px 16px;
    cursor: pointer;
    font-family: inherit;
    font-size: 14px;
    transition: color 0.15s, border-color 0.15s;
  }
  .tab:hover { color: #e0e0e0; }
  .tab.active {
    color: #00ff41;
    border-bottom-color: #00ff41;
  }
  .tab-content {
    margin-top: 16px;
    min-height: 300px;
  }
  .action-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 24px;
    padding-top: 16px;
    border-top: 1px solid #0f3460;
    flex-wrap: wrap;
  }
  .action {
    background: #0f3460;
    border: 1px solid #1a4a8a;
    color: #e0e0e0;
    padding: 6px 14px;
    border-radius: 4px;
    cursor: pointer;
    font-family: inherit;
    font-size: 13px;
    transition: background 0.15s;
  }
  .action:hover { background: #1a4a8a; }
  .action.save {
    background: #0a5c0a;
    border-color: #0f8f0f;
  }
  .action.save:hover { background: #0f8f0f; }
  .action.reload {
    background: #5c3a0a;
    border-color: #8f5c0f;
  }
  .action.reload:hover { background: #8f5c0f; }
  .spacer { flex: 1; }
  .status {
    font-size: 12px;
    color: #00ff41;
  }
  .error {
    color: #ff4444;
    padding: 8px;
    margin-bottom: 8px;
    background: rgba(255,68,68,0.1);
    border-radius: 4px;
  }
</style>
