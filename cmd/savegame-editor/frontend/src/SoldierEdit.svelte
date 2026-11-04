<script lang="ts">
  import { onMount, onDestroy, createEventDispatcher } from 'svelte';
  import { getSoldier, updateSoldier, soldierImageURL } from './lib/api';
  import { registerForm, ChangeTracker } from './lib/pending';
  import type { SoldierDetail } from './lib/types';

  export let gameSlot: string;
  export let soldierIdx: number;

  const dispatch = createEventDispatcher();

  let soldier: SoldierDetail | null = null;
  let error = '';
  let saving = false;

  // Editable fields
  let name = '';
  let initialTimeUnits = 0;
  let initialHealth = 0;
  let initialEnergy = 0;
  let initialReactions = 0;
  let initialStrength = 0;
  let initialFiringAccuracy = 0;
  let initialThrowingAccuracy = 0;
  let initialMeleeAccuracy = 0;
  let initialPsionicStrength = 0;
  let initialPsionicSkill = 0;
  let armor = '';

  async function load() {
    try {
      soldier = await getSoldier(gameSlot, soldierIdx);
      name = soldier.name;
      initialTimeUnits = soldier.initialTimeUnits;
      initialHealth = soldier.initialHealth;
      initialEnergy = soldier.initialEnergy;
      initialReactions = soldier.initialReactions;
      initialStrength = soldier.initialStrength;
      initialFiringAccuracy = soldier.initialFiringAccuracy;
      initialThrowingAccuracy = soldier.initialThrowingAccuracy;
      initialMeleeAccuracy = soldier.initialMeleeAccuracy;
      initialPsionicStrength = soldier.initialPsionicStrength;
      initialPsionicSkill = soldier.initialPsionicSkill;
      armor = soldier.armor;
      tracker.reset(values());
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);

  // Tracks which values the server has, so only changes are sent.
  const tracker = new ChangeTracker((v: ReturnType<typeof values>) =>
    updateSoldier(gameSlot, soldierIdx, v));

  function values() {
    return {
      name,
      initialTimeUnits,
      initialHealth,
      initialEnergy,
      initialReactions,
      initialStrength,
      initialFiringAccuracy,
      initialThrowingAccuracy,
      initialMeleeAccuracy,
      initialPsionicStrength,
      initialPsionicSkill,
      armor,
    };
  }

  // apply sends the form values to the server if they changed since they
  // were last loaded or applied.
  async function apply() {
    if (!soldier) return;
    await tracker.apply(values);
  }

  // Let the global "Save" button apply changes not yet applied here, and
  // "Reload" refresh the form.
  onDestroy(registerForm({ apply, refresh: load }));

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

  // images of the soldier without armour and in the selected armour
  $: images = soldier?.images ?? null;
  $: armorImage = images && armor !== 'None' ? images[armor] : undefined;

  const armorOptions = ['None', 'Personal Armour', 'Power Suit', 'Flying Suit'];
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#if soldier}
  <div class="edit-form">
    <h3>Edit: {soldier.name}</h3>
    <div class="info">
      {soldier.rank} &middot; {soldier.gender} &middot; {soldier.appearance}
      {#if soldier.baseName} &middot; Base: {soldier.baseName}{/if}
      {#if soldier.craftName} &middot; Craft: {soldier.craftName}{/if}
    </div>
    <div class="info">
      Missions: {soldier.missions} &middot; Kills: {soldier.kills}
      {#if soldier.recoveryDays > 0} &middot; Recovery: {soldier.recoveryDays} days{/if}
    </div>

    {#if images?.['None']}
      <div class="portraits">
        <figure>
          <img src={soldierImageURL(images['None'])} alt="{soldier.name} without armour" />
          <figcaption>No Armour</figcaption>
        </figure>
        {#if armorImage}
          <figure>
            <img src={soldierImageURL(armorImage)} alt="{soldier.name} in {armor}" />
            <figcaption>{armor}</figcaption>
          </figure>
        {/if}
      </div>
    {/if}

    <div class="fields">
      <label>
        <span>Name</span>
        <input type="text" bind:value={name} maxlength="25" />
      </label>

      <label>
        <span>Armor</span>
        <select bind:value={armor}>
          {#each armorOptions as opt}
            <option value={opt}>{opt}</option>
          {/each}
        </select>
      </label>
    </div>

    <h4>Initial Stats</h4>
    <div class="stats-grid">
      <label><span>Time Units</span><input type="number" bind:value={initialTimeUnits} min="0" max="255" /></label>
      <label><span>Health</span><input type="number" bind:value={initialHealth} min="0" max="255" /></label>
      <label><span>Energy</span><input type="number" bind:value={initialEnergy} min="0" max="255" /></label>
      <label><span>Reactions</span><input type="number" bind:value={initialReactions} min="0" max="255" /></label>
      <label><span>Strength</span><input type="number" bind:value={initialStrength} min="0" max="255" /></label>
      <label><span>Firing Acc.</span><input type="number" bind:value={initialFiringAccuracy} min="0" max="255" /></label>
      <label><span>Throwing Acc.</span><input type="number" bind:value={initialThrowingAccuracy} min="0" max="255" /></label>
      <label><span>Melee Acc.</span><input type="number" bind:value={initialMeleeAccuracy} min="0" max="255" /></label>
      <label><span>Psi Strength</span><input type="number" bind:value={initialPsionicStrength} min="0" max="255" /></label>
      <label><span>Psi Skill</span><input type="number" bind:value={initialPsionicSkill} min="0" max="255" /></label>
    </div>

    <h4>Current Totals (Initial + Improvements)</h4>
    <div class="stats-display">
      <div><span>TU:</span> {soldier.timeUnits}</div>
      <div><span>HP:</span> {soldier.health}</div>
      <div><span>Energy:</span> {soldier.energy}</div>
      <div><span>React:</span> {soldier.reactions}</div>
      <div><span>Str:</span> {soldier.strength}</div>
      <div><span>Fire:</span> {soldier.firingAccuracy}</div>
      <div><span>Throw:</span> {soldier.throwingAccuracy}</div>
      <div><span>Melee:</span> {soldier.meleeAccuracy}</div>
      <div><span>Psi Str:</span> {soldier.psionicStrength}</div>
      <div><span>Psi Skl:</span> {soldier.psionicSkill}</div>
      <div><span>Bravery:</span> {soldier.bravery}</div>
    </div>

    <div class="actions">
      <button class="save-btn" on:click={save} disabled={saving}>
        {saving ? 'Applying...' : 'Apply Changes'}
      </button>
    </div>
  </div>
{/if}

<style>
  /* like the soldier screen of the game */
  h3 {
    background: var(--row-2);
    border: 2px solid;
    border-color: var(--row-1) var(--row-4) var(--row-4) var(--row-1);
    padding: 2px 10px;
  }
  .info {
    color: var(--muted);
    margin-bottom: 2px;
  }
  .portraits {
    display: flex;
    gap: 24px;
    margin: 12px 0 4px 0;
  }
  .portraits figure {
    margin: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-end;
    gap: 4px;
  }
  .portraits img {
    /* pixel art, scaled up without blurring */
    zoom: 2;
    image-rendering: pixelated;
  }
  .portraits figcaption {
    color: var(--muted);
    text-transform: uppercase;
  }
  .fields {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin: 12px 0 8px 0;
  }
  label {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  label span {
    color: var(--muted);
    text-transform: uppercase;
  }
  .fields input[type='text'] {
    width: 300px;
  }
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 2px 8px;
  }
  .stats-grid label {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    background: var(--row-2);
    padding: 1px 2px 1px 8px;
  }
  .stats-grid label:nth-child(4n + 1) { background: var(--row-1); }
  .stats-grid label:nth-child(4n + 3) { background: var(--row-3); }
  .stats-grid label:nth-child(4n + 4) { background: var(--row-4); }
  .stats-grid label span {
    color: var(--text);
  }
  .stats-grid input {
    width: 70px;
    text-align: right;
  }
  .stats-display {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 2px 16px;
    color: var(--value);
  }
  .stats-display span {
    color: var(--muted);
    text-transform: uppercase;
  }
  .actions {
    margin-top: 16px;
  }
</style>
