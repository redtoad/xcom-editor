<script lang="ts">
  // Details of a unit on the mission map (used in the tooltip and in the
  // panel of the selected unit).
  import type { BattleMap, BattleUnit } from './lib/types';
  import { levelLabel } from './lib/battle';

  export let unit: BattleUnit;
  export let map: BattleMap;

  function ownerLabel(u: BattleUnit): string {
    switch (u.owner) {
      case 'xcom': return u.tank ? 'X-COM tank' : 'X-COM soldier';
      case 'alien': return 'Alien';
      case 'civilian': return 'Civilian';
      default: return 'Left behind';
    }
  }

  function state(u: BattleUnit): string {
    if (!u.alive) return 'Dead';
    if (u.unconscious) return 'Unconscious';
    return 'Alive';
  }
</script>

<div class="title">{unit.name || '(unnamed)'}</div>
<div class="type">{ownerLabel(unit)}</div>
<dl>
  <dt>State</dt><dd>{state(unit)}</dd>
  <dt>Position</dt><dd>{unit.x}, {unit.y}, {levelLabel(map, unit.z)}</dd>
  <dt>Health</dt><dd>{unit.health}/{unit.maxHealth}</dd>
  <dt>Time units</dt><dd>{unit.timeUnits}/{unit.maxTimeUnits}</dd>
  <dt>Morale</dt><dd>{unit.morale}</dd>
  {#if unit.owner !== 'xcom'}
    <dt>Spotted</dt><dd>{unit.visible ? 'Yes' : 'No'}</dd>
  {/if}
</dl>

<style>
  .title { color: var(--value); }
  .type {
    color: var(--muted);
    text-transform: uppercase;
  }
  dl {
    display: grid;
    grid-template-columns: auto auto;
    gap: 0 12px;
    margin: 4px 0 0 0;
  }
  dt {
    color: var(--dialog-text);
    white-space: nowrap;
  }
  dd {
    margin: 0;
    white-space: nowrap;
    color: var(--value);
  }
</style>
