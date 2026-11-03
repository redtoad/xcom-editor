<script lang="ts">
  import { onMount } from 'svelte';
  import { geoOrthographic, geoPath, geoGraticule10, geoDistance } from 'd3-geo';
  import type { GeoPermissibleObjects } from 'd3-geo';
  import { feature } from 'topojson-client';
  import type { Topology, GeometryCollection } from 'topojson-specification';
  import landTopology from 'world-atlas/land-110m.json';
  import { listLocations, getWorld, textureURL } from './lib/api';
  import type { GlobeLocation, WorldData } from './lib/types';
  import { loadImages, buildTerrainMap, renderGlobe } from './lib/terrainMap';

  export let gameSlot: string;

  const SIZE = 360;
  const RADIUS = SIZE / 2 - 4;
  // The textured globe is rendered with doubled pixels, like on the original
  // 320x200 screen.
  const PIXEL_SCALE = 2;

  // Location type codes (see geoscape.LocationType).
  const ALIEN_SHIP = 1;
  const XCOM_SHIP = 2;
  const XCOM_BASE = 3;

  const colorMap: Record<number, string> = {
    3: '#00ff41', // X-COM Base
    2: '#00ccff', // X-COM Craft
    1: '#ff4444', // Alien Ship
    4: '#ff6600', // Alien Base
    5: '#ff8800', // Crash Site
    6: '#ff4444', // Landed UFO
    8: '#ffff00', // Terror Site
  };

  const legendEntries = [
    { label: 'X-COM Craft', color: colorMap[2], shape: 'triangle' },
    { label: 'Alien Ship', color: colorMap[1], shape: 'diamond' },
    { label: 'X-COM Base', color: colorMap[3], shape: 'square' },
    { label: 'Other sites', color: colorMap[5], shape: 'circle' },
  ];

  const land = feature(
    landTopology as unknown as Topology,
    (landTopology as unknown as Topology).objects.land as GeometryCollection,
  );
  const graticule = geoGraticule10();

  let locations: GlobeLocation[] = [];
  // terrain from the game files, null if not available
  let world: WorldData | null = null;
  // the terrain drawn into an equirectangular map with the selected textures
  let terrainMap: ImageData | null = null;
  let canvases: HTMLCanvasElement[] = [];
  let textureSet = 0;
  let textureSets = 0;
  let error = '';
  let centerLon = 0;
  let showSites = true;
  let hovered: GlobeLocation[] | null = null;
  let tooltipX = 0;
  let tooltipY = 0;
  let wrapper: HTMLDivElement;

  async function load(slot: string) {
    try {
      locations = await listLocations(slot);
      error = '';
      // center the front view on the first X-COM base
      const base = locations.find(l => l.typeCode === XCOM_BASE);
      centerLon = base ? Math.round(base.coord.lon) : 0;
    } catch (e: any) {
      error = e.message;
    }
  }

  $: load(gameSlot);

  onMount(async () => {
    try {
      const data = await getWorld();
      if (data.available) {
        world = data;
        textureSets = data.textureSets ?? 0;
      }
    } catch {
      // fall back to the generic map
    }
  });

  $: if (world) updateTerrainMap(world, textureSet);

  async function updateTerrainMap(world: WorldData, set: number) {
    try {
      const urls = Array.from({ length: world.terrains ?? 0 }, (_, i) => textureURL(set, i));
      const textures = await loadImages(urls);
      terrainMap = buildTerrainMap(world.polygons ?? [], textures, world.ocean ?? '#0b1d3a');
    } catch (e: any) {
      error = e.message;
    }
  }

  function normalize(lon: number): number {
    return ((lon + 540) % 360) - 180;
  }

  $: hemispheres = [
    { label: 'Front', lon: normalize(centerLon) },
    { label: 'Back', lon: normalize(centerLon + 180) },
  ].map(h => {
    const projection = geoOrthographic()
      .scale(RADIUS)
      .translate([SIZE / 2, SIZE / 2])
      .rotate([-h.lon, 0])
      .clipAngle(90);
    // projection for the textured globe (rendered at a lower resolution)
    const pixelProjection = geoOrthographic()
      .scale(RADIUS / PIXEL_SCALE)
      .translate([SIZE / PIXEL_SCALE / 2, SIZE / PIXEL_SCALE / 2])
      .rotate([-h.lon, 0]);
    return { ...h, projection, pixelProjection, path: geoPath(projection) };
  });

  // render the textured globe whenever the view or the textures change
  $: if (terrainMap && canvases.length === hemispheres.length) {
    hemispheres.forEach((h, i) => {
      const ctx = canvases[i]?.getContext('2d');
      if (ctx) renderGlobe(ctx, h.pixelProjection, terrainMap!);
    });
  }

  $: ships = locations.filter(l => l.typeCode === XCOM_SHIP || l.typeCode === ALIEN_SHIP);
  $: sites = locations.filter(l => l.typeCode !== XCOM_SHIP && l.typeCode !== ALIEN_SHIP);
  $: markers = showSites ? [...sites, ...ships] : ships;
  $: groups = groupByPosition(markers);
  $: flights = ships
    .filter(s => s.destCoord)
    .map(s => ({ ship: s, dest: s.destCoord! }));

  interface MarkerGroup {
    coord: { lat: number; lon: number };
    // the location that is drawn (ships take precedence over sites)
    primary: GlobeLocation;
    items: GlobeLocation[];
  }

  // groupByPosition combines locations at the same position (e.g. crafts
  // parked at a base) into a single marker.
  function groupByPosition(locs: GlobeLocation[]): MarkerGroup[] {
    const groups = new Map<string, GlobeLocation[]>();
    for (const loc of locs) {
      const key = `${loc.coord.lat},${loc.coord.lon}`;
      groups.set(key, [...(groups.get(key) || []), loc]);
    }
    return [...groups.values()].map(items => {
      const sorted = [...items].sort((a, b) => priority(b) - priority(a));
      return { coord: items[0].coord, primary: sorted[0], items: sorted };
    });
  }

  function priority(loc: GlobeLocation): number {
    const inFlight = loc.destCoord ? 10 : 0;
    switch (loc.typeCode) {
      case ALIEN_SHIP: return inFlight + 3;
      case XCOM_SHIP: return inFlight + 2;
      case XCOM_BASE: return 1;
      default: return 0;
    }
  }

  function isVisible(loc: { lat: number; lon: number }, centerLon: number): boolean {
    return geoDistance([loc.lon, loc.lat], [centerLon, 0]) < Math.PI / 2;
  }

  function flightPath(ship: GlobeLocation, dest: { lat: number; lon: number }): GeoPermissibleObjects {
    return {
      type: 'LineString',
      coordinates: [
        [ship.coord.lon, ship.coord.lat],
        [dest.lon, dest.lat],
      ],
    };
  }

  function color(loc: GlobeLocation): string {
    return colorMap[loc.typeCode] || '#ffffff';
  }

  function shape(loc: GlobeLocation): string {
    switch (loc.typeCode) {
      case XCOM_SHIP: return 'M0,-6 L5,4 L-5,4 Z';
      case ALIEN_SHIP: return 'M0,-6 L6,0 L0,6 L-6,0 Z';
      case XCOM_BASE: return 'M-4,-4 L4,-4 L4,4 L-4,4 Z';
      default: return 'M-3.5,0 A3.5,3.5 0 1,0 3.5,0 A3.5,3.5 0 1,0 -3.5,0';
    }
  }

  function formatCoord(c: { lat: number; lon: number }): string {
    const lat = `${Math.abs(c.lat).toFixed(1)}°${c.lat < 0 ? 'S' : 'N'}`;
    const lon = `${Math.abs(c.lon).toFixed(1)}°${c.lon < 0 ? 'W' : 'E'}`;
    return `${lat} ${lon}`;
  }

  function showTooltip(event: MouseEvent, items: GlobeLocation[]) {
    const rect = wrapper.getBoundingClientRect();
    tooltipX = event.clientX - rect.left + 14;
    tooltipY = event.clientY - rect.top + 14;
    hovered = items;
  }

  function hideTooltip() {
    hovered = null;
  }
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

<div class="controls">
  <label>
    Rotate
    <input type="range" min="-180" max="180" step="5" bind:value={centerLon} />
    <span class="value">{formatCoord({ lat: 0, lon: normalize(centerLon) }).split(' ')[1]}</span>
  </label>
  <label>
    <input type="checkbox" bind:checked={showSites} />
    Show bases and sites
  </label>
  {#if world && textureSets > 1}
    <label>
      Textures
      <select bind:value={textureSet}>
        {#each Array(textureSets) as _, i}
          <option value={i}>Set {i + 1}</option>
        {/each}
      </select>
    </label>
  {/if}
  <span class="count">{ships.length} ships &middot; {flights.length} in flight</span>
</div>

<div class="map-wrapper" bind:this={wrapper}>
  {#each hemispheres as h, i (h.label)}
    <figure class="hemisphere" style="width:{SIZE}px">
      {#if world}
        <canvas
          class="globe-texture"
          bind:this={canvases[i]}
          width={SIZE / PIXEL_SCALE}
          height={SIZE / PIXEL_SCALE}
          style="width:{SIZE}px; height:{SIZE}px"
        ></canvas>
      {/if}
      <svg viewBox="0 0 {SIZE} {SIZE}" width={SIZE} height={SIZE} role="img" aria-label="{h.label} hemisphere">
        {#if world}
          <path class="rim" d={h.path({ type: 'Sphere' })} />
          <path class="graticule faint" d={h.path(graticule)} />
        {:else}
          <path class="ocean" d={h.path({ type: 'Sphere' })} />
          <path class="graticule" d={h.path(graticule)} />
          <path class="land" d={h.path(land)} />
        {/if}

        {#each flights as { ship, dest }}
          {@const d = h.path(flightPath(ship, dest))}
          {#if d}
            <!-- svelte-ignore a11y-no-static-element-interactions -->
            <path
              class="flight"
              {d}
              stroke={color(ship)}
              on:mousemove={(e) => showTooltip(e, [ship])}
              on:mouseleave={hideTooltip}
            />
            {#if isVisible(dest, h.lon)}
              {@const p = h.projection([dest.lon, dest.lat])}
              {#if p}
                <circle class="destination" cx={p[0]} cy={p[1]} r="4" stroke={color(ship)} />
              {/if}
            {/if}
          {/if}
        {/each}

        {#each groups as group}
          {#if isVisible(group.coord, h.lon)}
            {@const p = h.projection([group.coord.lon, group.coord.lat])}
            {#if p}
              <!-- svelte-ignore a11y-no-static-element-interactions -->
              <g
                class="marker"
                class:ship={group.primary.typeCode === XCOM_SHIP || group.primary.typeCode === ALIEN_SHIP}
                class:active={hovered === group.items}
                transform="translate({p[0]},{p[1]})"
                on:mousemove={(e) => showTooltip(e, group.items)}
                on:mouseleave={hideTooltip}
              >
                <path d={shape(group.primary)} fill={color(group.primary)} />
                {#if group.items.length > 1}
                  <text class="count-badge" x="7" y="-5">{group.items.length}</text>
                {/if}
              </g>
            {/if}
          {/if}
        {/each}
      </svg>
      <figcaption>{h.label} &middot; centered on {formatCoord({ lat: 0, lon: h.lon }).split(' ')[1]}</figcaption>
    </figure>
  {/each}

  {#if hovered}
    <div class="tooltip" style="left:{tooltipX}px; top:{tooltipY}px">
      {#each hovered as loc}
        <div class="tooltip-entry">
          <div class="tooltip-title" style="color:{color(loc)}">{loc.name}</div>
          <div class="tooltip-type">{loc.type}</div>
          {#if loc.craft?.damaged}
            <div class="tooltip-damaged">Damaged entry in CRAFT.DAT, values are meaningless</div>
          {/if}
          {#if loc.craft}
            <dl>
              <dt>Craft</dt><dd>{loc.craft.type}</dd>
              <dt>Status</dt><dd>{loc.craft.status}</dd>
              <dt>Altitude</dt><dd>{loc.craft.altitude}</dd>
              <dt>Speed</dt><dd>{loc.craft.speed}</dd>
              {#if loc.typeCode === XCOM_SHIP}
                <dt>Fuel</dt><dd>{loc.craft.fuel}</dd>
                <dt>Damage</dt><dd>{loc.craft.damage}</dd>
                {#if loc.craft.baseName}<dt>Base</dt><dd>{loc.craft.baseName}</dd>{/if}
              {:else}
                {#if loc.craft.race}<dt>Race</dt><dd>{loc.craft.race}</dd>{/if}
                {#if loc.craft.mission}<dt>Mission</dt><dd>{loc.craft.mission}</dd>{/if}
                {#if loc.craft.zone}<dt>Zone</dt><dd>{loc.craft.zone}</dd>{/if}
              {/if}
              {#if loc.destCoord}
                <dt>Heading to</dt>
                <dd>{loc.craft.destination || formatCoord(loc.destCoord)}</dd>
              {/if}
            </dl>
          {/if}
        </div>
      {/each}
      <div class="tooltip-position">{formatCoord(hovered[0].coord)}</div>
    </div>
  {/if}
</div>

<div class="legend">
  {#each legendEntries as entry}
    <span class="legend-item">
      <span class="swatch {entry.shape}" style="background:{entry.color}"></span>
      {entry.label}
    </span>
  {/each}
  <span class="legend-item"><span class="swatch line"></span> Flight path</span>
</div>

<style>
  .controls {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 20px;
    margin-bottom: 12px;
    color: var(--muted);
  }
  .controls label {
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .controls input[type='range'] {
    width: 200px;
    accent-color: var(--button);
  }
  .controls .value {
    width: 48px;
    color: var(--text);
  }
  .count {
    margin-left: auto;
    color: var(--muted);
  }
  .map-wrapper {
    position: relative;
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 24px;
  }
  .hemisphere {
    position: relative;
    margin: 0;
    text-align: center;
  }
  .globe-texture {
    position: absolute;
    top: 0;
    left: 0;
    image-rendering: pixelated;
  }
  .hemisphere svg {
    position: relative;
  }
  .rim {
    fill: none;
    stroke: var(--border);
    stroke-width: 1.5;
  }
  .hemisphere svg {
    max-width: 100%;
    height: auto;
  }
  figcaption {
    margin-top: 6px;
    color: var(--muted);
  }
  .ocean {
    fill: #0b1d3a;
    stroke: #1a4a8a;
    stroke-width: 1.5;
  }
  .graticule {
    fill: none;
    stroke: #16325c;
    stroke-width: 0.5;
  }
  .graticule.faint {
    stroke: rgba(255, 255, 255, 0.08);
  }
  .controls select {
    background: var(--panel);
    color: var(--text);
    border: 1px solid var(--border);
    font-family: inherit;
  }
  .land {
    fill: #1f3d2b;
    stroke: #2f6b45;
    stroke-width: 0.5;
  }
  .flight {
    fill: none;
    stroke-width: 2;
    stroke-dasharray: 6 4;
    animation: dash 1.2s linear infinite;
    cursor: pointer;
  }
  .flight:hover {
    stroke-width: 3.5;
  }
  @keyframes dash {
    to { stroke-dashoffset: -20; }
  }
  .destination {
    fill: none;
    stroke-width: 1.5;
    pointer-events: none;
  }
  .marker {
    cursor: pointer;
  }
  .marker path {
    stroke: #000;
    stroke-width: 1;
  }
  .marker.ship path {
    stroke: #fff;
  }
  .marker.active path,
  .marker:hover path {
    stroke: #fff;
    stroke-width: 2;
  }
  .count-badge {
    font-family: 'XcomSmall', monospace;
    fill: #fff;
    font-size: 9px;
    paint-order: stroke;
    stroke: #000;
    stroke-width: 2px;
  }
  .tooltip {
    position: absolute;
    z-index: 10;
    pointer-events: none;
    min-width: 240px;
    /* like the interception window of the game */
    background: var(--dialog);
    color: var(--dialog-text);
    border: 2px solid var(--dialog-border);
    padding: 6px 10px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.6);
  }
  .tooltip-type {
    color: var(--muted);
    margin-bottom: 4px;
  }
  .tooltip-entry + .tooltip-entry {
    margin-top: 8px;
    padding-top: 8px;
    border-top: 1px solid var(--border);
  }
  .tooltip-damaged {
    color: var(--danger);
    margin-bottom: 4px;
  }
  .tooltip-position {
    margin-top: 8px;
    color: var(--muted);
  }
  dl {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 2px 10px;
    margin: 0;
  }
  dt {
    color: var(--dialog-text);
    text-transform: uppercase;
  }
  dd {
    margin: 0;
    color: var(--value);
  }
  .legend {
    display: flex;
    flex-wrap: wrap;
    gap: 14px;
    justify-content: center;
    margin-top: 12px;
    color: var(--muted);
  }
  .legend-item {
    display: flex;
    align-items: center;
    gap: 5px;
  }
  .swatch {
    display: inline-block;
    width: 9px;
    height: 9px;
  }
  .swatch.circle { border-radius: 50%; }
  .swatch.diamond { transform: rotate(45deg); width: 8px; height: 8px; }
  .swatch.triangle { clip-path: polygon(50% 0, 100% 100%, 0 100%); }
  .swatch.line {
    width: 16px;
    height: 0;
    border-top: 2px dashed #00ccff;
  }
</style>
