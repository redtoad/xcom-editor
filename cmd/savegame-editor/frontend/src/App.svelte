<script lang="ts">
  import Sidebar from './Sidebar.svelte';
  import GameDetail from './GameDetail.svelte';

  let selectedSlot: string | null = null;

  function handleSelect(event: CustomEvent<string>) {
    selectedSlot = event.detail;
  }
</script>

<main>
  <div class="layout">
    <aside class="sidebar">
      <h1>X-COM Editor</h1>
      <Sidebar {selectedSlot} on:select={handleSelect} />
    </aside>
    <section class="detail">
      {#if selectedSlot}
        <GameDetail slot={selectedSlot} />
      {:else}
        <div class="placeholder">Select a savegame from the sidebar.</div>
      {/if}
    </section>
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, monospace;
    background: #1a1a2e;
    color: #e0e0e0;
  }
  :global(*) {
    box-sizing: border-box;
  }
  .layout {
    display: grid;
    grid-template-columns: 280px 1fr;
    height: 100vh;
  }
  .sidebar {
    background: #16213e;
    padding: 16px;
    overflow-y: auto;
    border-right: 1px solid #0f3460;
  }
  .sidebar h1 {
    font-size: 18px;
    color: #00ff41;
    margin: 0 0 16px 0;
    padding-bottom: 8px;
    border-bottom: 1px solid #0f3460;
  }
  .detail {
    padding: 16px 24px;
    overflow-y: auto;
  }
  .placeholder {
    color: #666;
    text-align: center;
    padding: 60px 20px;
    font-size: 16px;
  }
</style>
