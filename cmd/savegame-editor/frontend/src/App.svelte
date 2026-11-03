<script lang="ts">
  import Sidebar from './Sidebar.svelte';
  import GameDetail from './GameDetail.svelte';

  let selectedSlot: string | null = null;
  let refresh = 0;

  function handleSelect(event: CustomEvent<string>) {
    selectedSlot = event.detail;
  }
</script>

<main>
  <div class="layout">
    <aside class="sidebar">
      <h1>X-COM Editor</h1>
      <Sidebar {selectedSlot} {refresh} on:select={handleSelect} />
    </aside>
    <section class="detail">
      {#if selectedSlot}
        <GameDetail slot={selectedSlot} on:changed={() => refresh++} />
      {:else}
        <div class="placeholder">Select a savegame from the sidebar.</div>
      {/if}
    </section>
  </div>
</main>

<style>
  .layout {
    display: grid;
    grid-template-columns: 280px 1fr;
    height: 100vh;
  }
  .sidebar {
    background: var(--panel);
    padding: 12px;
    overflow-y: auto;
    border-right: 4px solid var(--border);
  }
  .sidebar h1 {
    margin: 0 0 12px 0;
    padding-bottom: 8px;
    border-bottom: 2px solid var(--border);
  }
  .detail {
    padding: 12px 20px;
    overflow-y: auto;
  }
  .placeholder {
    color: var(--muted);
    text-align: center;
    padding: 60px 20px;
    text-transform: uppercase;
  }
</style>
