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
    margin-bottom: 24px;
  }
  .balance-edit {
    display: flex;
    gap: 8px;
    align-items: center;
    margin-top: 4px;
  }
  .balance-display {
    font-size: 28px;
    color: #00ff41;
    font-weight: bold;
    margin-top: 8px;
  }
  h4 {
    color: #aaa;
    font-size: 13px;
    margin: 20px 0 8px 0;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  label span {
    font-size: 11px;
    color: #888;
    text-transform: uppercase;
    letter-spacing: 0.3px;
  }
  input {
    background: #16213e;
    border: 1px solid #0f3460;
    color: #e0e0e0;
    padding: 6px 8px;
    border-radius: 4px;
    font-family: inherit;
    font-size: 13px;
    width: 160px;
  }
  input:focus { outline: none; border-color: #00ff41; }
  .save-btn {
    background: #0a5c0a;
    border: 1px solid #0f8f0f;
    color: #e0e0e0;
    padding: 6px 14px;
    border-radius: 4px;
    cursor: pointer;
    font-family: inherit;
    font-size: 13px;
  }
  .save-btn:hover { background: #0f8f0f; }
  .save-btn:disabled { opacity: 0.5; cursor: not-allowed; }
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
  .negative { color: #ff4444; }
  .error {
    color: #ff4444;
    padding: 8px;
    margin-bottom: 8px;
    background: rgba(255,68,68,0.1);
    border-radius: 4px;
  }
</style>
