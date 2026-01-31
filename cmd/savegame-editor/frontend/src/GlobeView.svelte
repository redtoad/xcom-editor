<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { listLocations } from './lib/api';
  import type { GlobeLocation } from './lib/types';
  import Globe from 'globe.gl';
  import { MeshBasicMaterial, BackSide } from 'three';

  export let gameSlot: string;

  let locations: GlobeLocation[] = [];
  let error = '';
  let container: HTMLDivElement;
  let globe: ReturnType<typeof Globe> | null = null;

  const EARTH_NIGHT_IMG = '/earth-day.jpg';

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
    { typeCode: 3, label: 'X-COM Base', color: '#00ff41' },
    { typeCode: 2, label: 'X-COM Craft', color: '#00ccff' },
    { typeCode: 1, label: 'Alien Ship', color: '#ff4444' },
    { typeCode: 4, label: 'Alien Base', color: '#ff6600' },
    { typeCode: 5, label: 'Crash Site', color: '#ff8800' },
    { typeCode: 6, label: 'Landed UFO', color: '#ff4444' },
    { typeCode: 8, label: 'Terror Site', color: '#ffff00' },
  ];

  function pointColor(d: GlobeLocation): string {
    return colorMap[d.typeCode] || '#ffffff';
  }

  function pointRadius(d: GlobeLocation): number {
    return d.typeCode === 3 ? 0.6 : 0.4;
  }

  function pointLabel(d: GlobeLocation): string {
    const c = pointColor(d);
    return `<div style="color:${c};font-family:monospace;font-size:13px;text-shadow:0 0 4px #000">
      <b>${d.name}</b><br/><span style="opacity:0.7">${d.type}</span>
    </div>`;
  }

  async function load() {
    try {
      locations = await listLocations(gameSlot);
      error = '';
      if (globe) {
        globe.pointsData(locations);
      }
    } catch (e: any) {
      error = e.message;
    }
  }

  function initGlobe() {
    if (!container) return;

    const material = new MeshBasicMaterial({
      transparent: true,
      opacity: 0.7,
    });
    // Show the back side as well for a see-through effect
    const backMaterial = new MeshBasicMaterial({
      transparent: true,
      opacity: 0.3,
      side: BackSide,
    });

    const width = container.clientWidth;
    const height = Math.min(width, 600);

    globe = Globe()
      (container)
      .globeImageUrl(EARTH_NIGHT_IMG)
      .globeMaterial(material)
      .backgroundColor('rgba(0,0,0,0)')
      .showAtmosphere(false)
      .width(width)
      .height(height)
      .pointsData(locations)
      .pointLat((d: any) => d.coord.lat)
      .pointLng((d: any) => d.coord.lon)
      .pointColor(pointColor as any)
      .pointRadius(pointRadius as any)
      .pointAltitude(0.01)
      .pointLabel(pointLabel as any);

    // Add a translucent inner sphere for the back-side visibility
    const scene = globe.scene();
    if (scene && scene.children) {
      for (const child of scene.children) {
        if ((child as any).type === 'Mesh' && (child as any).__globeObjType === 'globe') {
          const backGlobe = child.clone();
          (backGlobe as any).material = backMaterial;
          scene.add(backGlobe);
          break;
        }
      }
    }
  }

  onMount(async () => {
    await load();
    initGlobe();
  });

  $: if (gameSlot) load();

  onDestroy(() => {
    if (globe) {
      try {
        globe.renderer()?.dispose();
        globe.controls()?.dispose();
      } catch {}
      globe = null;
    }
  });
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

<div class="globe-wrapper">
  <div class="globe-container" bind:this={container}></div>
  <div class="legend">
    {#each legendEntries as entry}
      <span class="legend-item">
        <span class="dot" style="background:{entry.color}"></span>
        {entry.label}
      </span>
    {/each}
  </div>
</div>

<style>
  .globe-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .globe-container {
    width: 100%;
    max-width: 600px;
  }
  .globe-container :global(canvas) {
    cursor: grab;
  }
  .globe-container :global(canvas:active) {
    cursor: grabbing;
  }
  .legend {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: center;
    margin-top: 12px;
    font-size: 12px;
    color: #aaa;
  }
  .legend-item {
    display: flex;
    align-items: center;
    gap: 4px;
  }
  .dot {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }
  .error {
    color: #ff4444;
    padding: 8px;
    margin-bottom: 8px;
    background: rgba(255,68,68,0.1);
    border-radius: 4px;
  }
</style>
