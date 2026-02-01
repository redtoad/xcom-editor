<script lang="ts">
  import { onMount } from 'svelte';
  import { listTransfers } from './lib/api';
  import type { TransferSummary } from './lib/types';

  export let gameSlot: string;

  let transfers: TransferSummary[] = [];
  let error = '';

  async function load() {
    try {
      transfers = await listTransfers(gameSlot);
      error = '';
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

{#if transfers.length === 0 && !error}
  <div class="empty">No active transfers.</div>
{:else}
  <table>
    <thead>
      <tr>
        <th>Origin</th>
        <th>Destination</th>
        <th>Hours Left</th>
        <th>Type</th>
        <th>Quantity</th>
      </tr>
    </thead>
    <tbody>
      {#each transfers as t}
        <tr>
          <td>{t.origin === 255 ? 'Purchase' : `Base ${t.origin}`}</td>
          <td>Base {t.destination}</td>
          <td class="num">{t.hoursLeft}h</td>
          <td class="num">{t.type}</td>
          <td class="num">{t.quantity}</td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style>
  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 13px;
  }
  th {
    text-align: left;
    padding: 6px 10px;
    border-bottom: 1px solid #0f3460;
    color: #888;
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  td {
    padding: 6px 10px;
    border-bottom: 1px solid #16213e;
  }
  .num { text-align: right; font-variant-numeric: tabular-nums; }
  .empty {
    color: #666;
    text-align: center;
    padding: 40px;
  }
  .error {
    color: #ff4444;
    padding: 8px;
    margin-bottom: 8px;
    background: rgba(255,68,68,0.1);
    border-radius: 4px;
  }
</style>
