<script lang="ts">
  import { onMount } from 'svelte';
  import { getFinancials, updateFinancials } from './lib/api';
  import type { Financials } from './lib/types';

  export let gameSlot: string;

  let financials: Financials | null = null;
  let error = '';
  let balance = 0;
  let saving = false;
  let saved = false;

  async function load() {
    try {
      financials = await getFinancials(gameSlot);
      balance = financials.currentBalance;
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);
  $: if (gameSlot) load();

  async function save() {
    saving = true;
    try {
      await updateFinancials(gameSlot, { currentBalance: balance });
      saved = true;
      setTimeout(() => saved = false, 2000);
    } catch (e: any) {
      error = e.message;
    }
    saving = false;
  }

  function formatCurrency(val: number): string {
    return '$' + val.toLocaleString();
  }

  const months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#if financials}
  <div class="balance-section">
    <label>
      <span>Current Balance</span>
      <div class="balance-edit">
        <input type="number" bind:value={balance} />
        <button class="save-btn" on:click={save} disabled={saving}>
          {saving ? '...' : saved ? 'Saved!' : 'Update'}
        </button>
      </div>
    </label>
    <div class="balance-display">{formatCurrency(balance)}</div>
  </div>

  <h4>Monthly History (Last 12 Months)</h4>
  <table>
    <thead>
      <tr>
        <th>Month</th>
        <th>Expenditure</th>
        <th>Maintenance</th>
        <th>Balance</th>
      </tr>
    </thead>
    <tbody>
      {#each financials.expenditure as _, i}
        <tr>
          <td>{months[i] ?? `M${i+1}`}</td>
          <td class="num">{formatCurrency(financials.expenditure[i])}</td>
          <td class="num">{formatCurrency(financials.maintenance[i])}</td>
          <td class="num" class:negative={financials.balance[i] < 0}>
            {formatCurrency(financials.balance[i])}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style>
  .balance-section {
    display: flex;
    align-items: flex-end;
    gap: 24px;
    margin-bottom: 8px;
  }
  label {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  label span {
    color: var(--danger);
    text-transform: uppercase;
  }
  .balance-edit {
    display: flex;
    gap: 6px;
  }
  .balance-edit input {
    width: 220px;
  }
  .balance-display {
    font-family: var(--font-title);
    font-size: var(--fs-title);
    color: var(--value);
  }
  td.negative {
    color: var(--danger);
  }
</style>
