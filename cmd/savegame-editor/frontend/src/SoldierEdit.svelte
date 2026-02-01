<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { getSoldier, updateSoldier } from './lib/api';
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
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);

  async function save() {
    saving = true;
    try {
      await updateSoldier(gameSlot, soldierIdx, {
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
      });
      dispatch('saved');
    } catch (e: any) {
      error = e.message;
    }
    saving = false;
  }

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
  .edit-form h3 {
    color: #00ff41;
    margin: 0 0 8px 0;
  }
  .info {
    color: #888;
    font-size: 13px;
    margin-bottom: 4px;
  }
  h4 {
    color: #aaa;
    font-size: 13px;
    margin: 16px 0 8px 0;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  .fields {
    display: flex;
    gap: 16px;
    margin-top: 16px;
    flex-wrap: wrap;
  }
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 8px;
  }
  .stats-display {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 4px;
    font-size: 13px;
  }
  .stats-display span {
    color: #888;
  }
  label {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  label span {
    font-size: 11px;
    color: #888;
    text-transform: uppercase;
    letter-spacing: 0.3px;
  }
  input, select {
    background: #16213e;
    border: 1px solid #0f3460;
    color: #e0e0e0;
    padding: 6px 8px;
    border-radius: 4px;
    font-family: inherit;
    font-size: 13px;
  }
  input:focus, select:focus {
    outline: none;
    border-color: #00ff41;
  }
  input[type="number"] {
    width: 80px;
  }
  .actions {
    margin-top: 20px;
  }
  .save-btn {
    background: #0a5c0a;
    border: 1px solid #0f8f0f;
    color: #e0e0e0;
    padding: 8px 20px;
    border-radius: 4px;
    cursor: pointer;
    font-family: inherit;
    font-size: 14px;
  }
  .save-btn:hover { background: #0f8f0f; }
  .save-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .error {
    color: #ff4444;
    padding: 8px;
    margin-bottom: 8px;
    background: rgba(255,68,68,0.1);
    border-radius: 4px;
  }
</style>
