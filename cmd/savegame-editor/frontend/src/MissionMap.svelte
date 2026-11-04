<script lang="ts">
  import { getBattleMap } from './lib/api';
  import type { BattleMap, BattleUnit } from './lib/types';

  export let gameSlot: string;

  // flags of the tiles (see battleMapJS in api.go)
  const TILE_GROUND = 1;
  const TILE_WEST_WALL = 2;
  const TILE_NORTH_WALL = 4;
  const TILE_OBJECT = 8;
  const TILE_SEEN = 16;

  // size of a tile in SVG units
  const T = 12;
  // map modules are 10x10 tiles
  const MODULE = 10;

  let map: BattleMap | null = null;
  let error = '';
  let level = 0;
  let fogOfWar = false;
  let showDead = true;
  let showOtherLevels = true;
  let hovered: BattleUnit | null = null;
  let selected: BattleUnit | null = null;
  let tooltipX = 0;
  let tooltipY = 0;
  let wrapper: HTMLDivElement;

  async function load(slot: string) {
    try {
      map = await getBattleMap(slot);
      error = '';
      // start on the level with the most X-COM units
      const counts = new Array(map.sizeZ).fill(0);
      map.units.filter(u => u.owner === 'xcom' && u.alive).forEach(u => counts[u.z]++);
      level = counts.indexOf(Math.max(...counts));
    } catch (e: any) {
      error = e.message;
      map = null;
    }
  }

  $: load(gameSlot);

  function tile(m: BattleMap, x: number, y: number, z: number): number {
    return m.tiles[z * m.sizeX * m.sizeY + y * m.sizeX + x];
  }

  // The terrain of a level is drawn as one path per kind of tile part.
  function terrain(m: BattleMap, z: number, fog: boolean) {
    let ground = '', objects = '', walls = '', unseen = '';
    for (let y = 0; y < m.sizeY; y++) {
      for (let x = 0; x < m.sizeX; x++) {
        const t = tile(m, x, y, z);
        const px = x * T, py = y * T;
        if (t & TILE_GROUND) ground += `M${px},${py}h${T}v${T}h${-T}z`;
        if (t & TILE_OBJECT) objects += `M${px + 2},${py + 2}h${T - 4}v${T - 4}h${-(T - 4)}z`;
        if (t & TILE_WEST_WALL) walls += `M${px},${py}v${T}`;
        if (t & TILE_NORTH_WALL) walls += `M${px},${py}h${T}`;
        if (fog && !(t & TILE_SEEN)) unseen += `M${px},${py}h${T}v${T}h${-T}z`;
      }
    }
    return { ground, objects, walls, unseen };
  }

  $: paths = map ? terrain(map, level, fogOfWar) : null;

  $: units = map
    ? map.units
        .filter(u => showDead || u.alive)
        .filter(u => showOtherLevels || u.z === level)
        // units on the current level are drawn last (on top)
        .sort((a, b) => Number(a.z === level) - Number(b.z === level))
    : [];

  $: counts = map ? countUnits(map.units) : null;

  function countUnits(all: BattleUnit[]) {
    const c = { xcom: [0, 0], alien: [0, 0], civilian: [0, 0] } as Record<string, number[]>;
    for (const u of all) {
      if (!c[u.owner]) continue;
      c[u.owner][1]++;
      if (u.alive) c[u.owner][0]++;
    }
    return c;
  }

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

  function levelLabel(z: number): string {
    return `Level ${z + 1}`;
  }

  function onHover(event: MouseEvent, u: BattleUnit | null) {
    hovered = u;
    if (u && wrapper) {
      const rect = wrapper.getBoundingClientRect();
      tooltipX = event.clientX - rect.left + 14;
      tooltipY = event.clientY - rect.top + 14;
    }
  }

  function select(u: BattleUnit) {
    selected = selected?.index === u.index ? null : u;
    if (selected) level = u.z;
  }

  $: listed = map
    ? [...map.units].sort((a, b) =>
        ['xcom', 'alien', 'civilian', 'left-behind'].indexOf(a.owner) -
        ['xcom', 'alien', 'civilian', 'left-behind'].indexOf(b.owner))
    : [];
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#if map && paths && counts}
  <div class="controls">
    <span class="levels">
      {#each Array(map.sizeZ) as _, z}
        <button class="tab" class:active={level === z} on:click={() => level = z}>{levelLabel(z)}</button>
      {/each}
    </span>
    <label><input type="checkbox" bind:checked={fogOfWar} /> Fog of war</label>
    <label><input type="checkbox" bind:checked={showDead} /> Dead units</label>
    <label><input type="checkbox" bind:checked={showOtherLevels} /> Other levels</label>
    <span class="summary">
      Turn <b>{map.turn}</b>
      &middot; <span class="xcom">X-COM {counts.xcom[0]}/{counts.xcom[1]}</span>
      &middot; <span class="alien">Aliens {counts.alien[0]}/{counts.alien[1]}</span>
      {#if counts.civilian[1]}&middot; <span class="civilian">Civilians {counts.civilian[0]}/{counts.civilian[1]}</span>{/if}
    </span>
  </div>

  <div class="layout">
    <div class="map-wrapper" bind:this={wrapper}>
      <svg viewBox="-1 -1 {map.sizeX * T + 2} {map.sizeY * T + 2}" role="img" aria-label="Mission map">
        <rect class="void" x="0" y="0" width={map.sizeX * T} height={map.sizeY * T} />
        <path class="ground" d={paths.ground} />
        <path class="object" d={paths.objects} />
        <path class="wall" d={paths.walls} />

        <!-- map modules -->
        {#each Array(Math.ceil(map.sizeX / MODULE) + 1) as _, i}
          <line class="module" x1={i * MODULE * T} y1="0" x2={i * MODULE * T} y2={map.sizeY * T} />
        {/each}
        {#each Array(Math.ceil(map.sizeY / MODULE) + 1) as _, i}
          <line class="module" x1="0" y1={i * MODULE * T} x2={map.sizeX * T} y2={i * MODULE * T} />
        {/each}

        <path class="unseen" d={paths.unseen} />

        {#each map.crafts as c}
          <g class="craft" class:alien={c.alien}>
            <rect x={c.x * T} y={c.y * T} width={c.width * T} height={c.height * T} />
            <text x={c.x * T + 4} y={c.y * T + 14}>{c.alien ? 'UFO' : 'X-COM'}: {c.name}</text>
          </g>
        {/each}

        {#each units as u (u.index)}
          {@const size = u.large ? 2 : 1}
          {@const cx = (u.x + size / 2) * T}
          {@const cy = (u.y + size / 2) * T}
          {@const r = (size * T) / 2 - 1}
          <!-- svelte-ignore a11y-no-static-element-interactions a11y-click-events-have-key-events -->
          <g
            class="unit {u.owner}"
            class:dead={!u.alive}
            class:unconscious={u.unconscious}
            class:other-level={u.z !== level}
            class:selected={selected?.index === u.index}
            on:mousemove={e => onHover(e, u)}
            on:mouseleave={e => onHover(e, null)}
            on:click={() => select(u)}
          >
            {#if !u.alive}
              <path d="M{cx - r},{cy - r}L{cx + r},{cy + r}M{cx + r},{cy - r}L{cx - r},{cy + r}" />
            {:else if u.large}
              <rect x={cx - r} y={cy - r} width={2 * r} height={2 * r} />
            {:else}
              <circle {cx} {cy} {r} />
            {/if}
            {#if u.z !== level}
              <text x={cx} y={cy + 3}>{u.z > level ? '▼' : '▲'}</text>
            {/if}
          </g>
        {/each}
      </svg>

      {#if hovered}
        <div class="tooltip" style="left:{tooltipX}px; top:{tooltipY}px">
          <div class="tooltip-title {hovered.owner}">{hovered.name || '(unnamed)'}</div>
          <div class="tooltip-type">{ownerLabel(hovered)}</div>
          <dl>
            <dt>State</dt><dd>{state(hovered)}</dd>
            <dt>Position</dt><dd>{hovered.x}, {hovered.y}, {levelLabel(hovered.z)}</dd>
            <dt>Health</dt><dd>{hovered.health}/{hovered.maxHealth}</dd>
            <dt>Time units</dt><dd>{hovered.timeUnits}/{hovered.maxTimeUnits}</dd>
            <dt>Morale</dt><dd>{hovered.morale}</dd>
            {#if hovered.owner !== 'xcom'}
              <dt>Spotted</dt><dd>{hovered.visible ? 'Yes' : 'No'}</dd>
            {/if}
          </dl>
        </div>
      {/if}
    </div>

    <ul class="unit-list">
      {#each listed as u (u.index)}
        <li>
          <button
            class="unit-entry {u.owner}"
            class:dead={!u.alive}
            class:selected={selected?.index === u.index}
            on:click={() => select(u)}
          >
            <span class="name">{u.name || '(unnamed)'}</span>
            <span class="where">{u.x},{u.y} L{u.z + 1}</span>
          </button>
        </li>
      {/each}
    </ul>
  </div>

  <div class="legend">
    <span class="legend-item"><span class="swatch xcom"></span> X-COM</span>
    <span class="legend-item"><span class="swatch alien"></span> Alien</span>
    <span class="legend-item"><span class="swatch civilian"></span> Civilian</span>
    <span class="legend-item"><span class="swatch large"></span> Large unit</span>
    <span class="legend-item">&times; Dead</span>
    <span class="legend-item">▲▼ Other level</span>
    {#if map.crafts.length}<span class="legend-item"><span class="swatch craft"></span> Craft (approx. area)</span>{/if}
  </div>
{/if}

<style>
  .controls {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 16px;
    margin-bottom: 8px;
  }
  .levels {
    display: inline-flex;
    gap: 4px;
  }
  .tab.active {
    background: var(--button-dark);
    color: var(--value);
  }
  .summary {
    color: var(--muted);
    text-transform: uppercase;
  }
  .summary b {
    color: var(--value);
    font-weight: normal;
  }
  .xcom { color: var(--xcom); }
  .alien { color: var(--alien); }
  .civilian { color: var(--base); }

  .layout {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    align-items: flex-start;
  }
  .map-wrapper {
    position: relative;
    flex: 1 1 480px;
    max-width: 760px;
  }
  svg {
    display: block;
    width: 100%;
    height: auto;
    background: var(--bg);
  }
  .void { fill: #050510; }
  .ground { fill: var(--row-4); }
  .object { fill: var(--dim); }
  .wall {
    stroke: var(--button-light);
    stroke-width: 2;
    stroke-linecap: square;
    fill: none;
  }
  .module {
    stroke: var(--border);
    stroke-width: 0.5;
    stroke-dasharray: 2 3;
  }
  .unseen { fill: rgba(0, 0, 0, 0.75); }

  .craft rect {
    fill: rgba(116, 196, 252, 0.08);
    stroke: var(--xcom);
    stroke-width: 2;
    stroke-dasharray: 6 4;
  }
  .craft.alien rect {
    fill: rgba(236, 84, 84, 0.1);
    stroke: var(--alien);
  }
  .craft text {
    font-family: var(--font-text);
    font-size: 11px;
    fill: var(--xcom);
  }
  .craft.alien text { fill: var(--alien); }

  .unit { cursor: pointer; }
  .unit circle, .unit rect {
    stroke: var(--bg);
    stroke-width: 1.5;
  }
  .unit.xcom circle, .unit.xcom rect { fill: var(--xcom); }
  .unit.alien circle, .unit.alien rect { fill: var(--alien); }
  .unit.civilian circle, .unit.civilian rect { fill: var(--base); }
  .unit.left-behind circle, .unit.left-behind rect { fill: var(--muted); }
  .unit.unconscious circle, .unit.unconscious rect {
    fill: var(--bg);
    stroke-width: 2;
  }
  .unit.xcom.unconscious * { stroke: var(--xcom); }
  .unit.alien.unconscious * { stroke: var(--alien); }
  .unit.civilian.unconscious * { stroke: var(--base); }
  .unit.dead path {
    stroke-width: 2.5;
    stroke-linecap: round;
  }
  .unit.xcom.dead path { stroke: var(--xcom); }
  .unit.alien.dead path { stroke: var(--alien); }
  .unit.civilian.dead path { stroke: var(--base); }
  .unit.left-behind.dead path { stroke: var(--muted); }
  .unit.other-level { opacity: 0.45; }
  .unit text {
    font-size: 7px;
    text-anchor: middle;
    fill: var(--bg);
    pointer-events: none;
  }
  .unit.selected circle, .unit.selected rect, .unit.selected path {
    stroke: var(--value);
    stroke-width: 3;
  }
  .unit:hover { opacity: 1; }

  .tooltip {
    position: absolute;
    z-index: 10;
    min-width: 200px;
    padding: 6px 10px;
    pointer-events: none;
    background: var(--dialog);
    border: 2px solid var(--dialog-border);
    color: var(--dialog-text);
  }
  .tooltip-title { color: var(--value); }
  .tooltip-type {
    color: var(--muted);
    text-transform: uppercase;
  }
  .tooltip dl {
    display: grid;
    grid-template-columns: auto auto;
    gap: 0 12px;
    margin: 4px 0 0 0;
  }
  .tooltip dt {
    color: var(--dialog-text);
    white-space: nowrap;
  }
  .tooltip dd {
    margin: 0;
    white-space: nowrap;
    color: var(--value);
  }

  .unit-list {
    flex: 0 1 300px;
    list-style: none;
    margin: 0;
    padding: 0;
    max-height: 760px;
    overflow-y: auto;
  }
  .unit-entry {
    display: flex;
    justify-content: space-between;
    gap: 8px;
    width: 100%;
    padding: 2px 6px;
    background: none;
    border: 1px solid transparent;
    text-align: left;
    cursor: pointer;
    font: inherit;
  }
  .unit-entry.xcom { color: var(--xcom); }
  .unit-entry.alien { color: var(--alien); }
  .unit-entry.civilian { color: var(--base); }
  .unit-entry.left-behind { color: var(--muted); }
  .unit-entry.dead { opacity: 0.5; text-decoration: line-through; }
  .unit-entry.selected { border-color: var(--value); }
  .unit-entry .where { color: var(--muted); }

  .legend {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin-top: 8px;
    color: var(--muted);
  }
  .legend-item {
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }
  .swatch {
    display: inline-block;
    width: 12px;
    height: 12px;
    border-radius: 50%;
  }
  .swatch.xcom { background: var(--xcom); }
  .swatch.alien { background: var(--alien); }
  .swatch.civilian { background: var(--base); }
  .swatch.large {
    border-radius: 0;
    background: var(--muted);
  }
  .swatch.craft {
    border-radius: 0;
    border: 2px dashed var(--alien);
  }
</style>
