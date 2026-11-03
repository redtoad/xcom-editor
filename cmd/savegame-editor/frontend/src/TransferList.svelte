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
  .empty {
    color: var(--muted);
    text-align: center;
    padding: 40px;
    text-transform: uppercase;
  }
</style>
