<script lang="ts">
  // Isometric view of the Battlescape map, similar to the view of the game.
  // The terrain is drawn with the sprites of the game if they are available
  // (see terrain), otherwise with simple shapes: floors are diamonds, walls
  // and objects are blocks. Units are always drawn as pillars.
  import { createEventDispatcher, onMount } from 'svelte';
  import type { BattleMap, BattleTerrain, BattleUnit } from './lib/types';
  import { terrainSheetURL } from './lib/api';
  import { TILE_GROUND, TILE_WEST_WALL, TILE_NORTH_WALL, TILE_OBJECT, TILE_SEEN, tileAt, darken } from './lib/battle';

  export let map: BattleMap;
  // the highest visible level (in the order of MAP.DAT, 0 = top level)
  export let level: number;
  export let fogOfWar = false;
  export let showDead = true;
  export let selected: BattleUnit | null = null;
  // scale of the map (1 = one canvas pixel per CSS pixel)
  export let zoom = 1;
  // sprites of the map, null to draw the terrain with simple shapes
  export let terrain: BattleTerrain | null = null;

  const dispatch = createEventDispatcher<{
    hover: { event: MouseEvent; unit: BattleUnit | null };
    select: BattleUnit;
  }>();

  // size of a tile (width and height of the diamond) and of a level
  const TW = 32;
  const TH = 16;
  const LH = 24;
  const PAD = 8;
  // space above the top level for walls and objects reaching higher
  const TOP = PAD + 2 * LH;
  // size of a terrain sprite, the north corner of the floor is at (16, 24)
  const SPRITE_W = 32;
  const SPRITE_H = 40;
  const SPRITE_FLOOR_Y = 24;

  let canvas: HTMLCanvasElement;
  let scroller: HTMLDivElement;
  let colors: Record<string, string> = {};
  // screen areas of the drawn units, used to find the unit under the mouse
  let hitAreas: { unit: BattleUnit; x: number; y: number; w: number; h: number }[] = [];
  // sprite sheets of the terrain sets, in normal and darkened (fog of war)
  // colors, null until they are loaded
  let sheets: { normal: CanvasImageSource; dark: CanvasImageSource }[] | null = null;

  $: loadSheets(terrain);

  async function loadSheets(t: BattleTerrain | null) {
    sheets = null;
    if (!t) return;
    const loaded = await Promise.all(t.sets.map(name => new Promise<HTMLImageElement>((resolve, reject) => {
      const img = new Image();
      img.onload = () => resolve(img);
      img.onerror = () => reject(new Error(`could not load terrain ${name}`));
      img.src = terrainSheetURL(name);
    })));
    if (t !== terrain) return; // changed while loading
    sheets = loaded.map(img => ({ normal: img, dark: darkened(img) }));
  }

  // returns a copy of the image with darkened colors
  function darkened(img: HTMLImageElement): HTMLCanvasElement {
    const c = document.createElement('canvas');
    c.width = img.width;
    c.height = img.height;
    const ctx = c.getContext('2d')!;
    ctx.drawImage(img, 0, 0);
    ctx.globalCompositeOperation = 'source-atop';
    ctx.fillStyle = 'rgba(0, 0, 0, 0.7)';
    ctx.fillRect(0, 0, c.width, c.height);
    return c;
  }

  onMount(() => {
    const style = getComputedStyle(canvas);
    const v = (name: string) => style.getPropertyValue(name).trim();
    colors = {
      floor: v('--row-3'), floorDark: v('--row-4'),
      wallLeft: v('--button'), wallRight: v('--button-dark'), wallTop: v('--button-light'),
      object: v('--dim'), objectTop: v('--muted'),
      xcom: v('--xcom'), alien: v('--alien'), civilian: v('--base'), 'left-behind': v('--muted'),
      selected: v('--value'), bg: v('--bg'),
    };
    draw();
    centerOn(selected ?? map.units.filter(u => u.owner === 'xcom' && u.alive));
  });

  // scrolls the view so that the given units are in the middle
  export function centerOn(target: BattleUnit | BattleUnit[]) {
    const units = Array.isArray(target) ? target : [target];
    if (!scroller || units.length === 0) return;
    let sx = 0, sy = 0;
    for (const u of units) {
      const [px, py] = project(u.x + 0.5, u.y + 0.5, height(u.z));
      sx += px;
      sy += py;
    }
    scroller.scrollTo({
      left: (sx / units.length) * zoom - scroller.clientWidth / 2,
      top: (sy / units.length) * zoom - scroller.clientHeight / 2,
    });
  }

  // keep the selected unit in view
  $: if (selected && scroller) centerOn(selected);

  $: if (canvas && colors.floor) draw(map, level, fogOfWar, showDead, selected, zoom, sheets);

  // height of level z above the ground
  function height(z: number): number {
    return map.sizeZ - 1 - z;
  }

  $: width = (map.sizeX + map.sizeY) * (TW / 2) + 2 * PAD;
  $: canvasHeight = TOP + (map.sizeZ - 1) * LH + (map.sizeX + map.sizeY) * (TH / 2) + PAD;

  // screen position of the north corner of tile (x, y) at height h
  function project(x: number, y: number, h: number): [number, number] {
    return [
      PAD + map.sizeY * (TW / 2) + (x - y) * (TW / 2),
      TOP + (map.sizeZ - 1) * LH + (x + y) * (TH / 2) - h * LH,
    ];
  }

  function polygon(ctx: CanvasRenderingContext2D, points: [number, number][], fill: string, stroke?: string) {
    ctx.beginPath();
    ctx.moveTo(points[0][0], points[0][1]);
    for (const p of points.slice(1)) ctx.lineTo(p[0], p[1]);
    ctx.closePath();
    ctx.fillStyle = fill;
    ctx.fill();
    if (stroke) {
      ctx.strokeStyle = stroke;
      ctx.lineWidth = 1;
      ctx.stroke();
    }
  }

  // a block on the tile, inset from the tile edges, with the given height
  function block(ctx: CanvasRenderingContext2D, x0: number, y0: number, x1: number, y1: number, h: number, bh: number,
                 top: string, left: string, right: string) {
    const n = project(x0, y0, h), e = project(x1, y0, h), s = project(x1, y1, h), w = project(x0, y1, h);
    const up = (p: [number, number]): [number, number] => [p[0], p[1] - bh];
    polygon(ctx, [w, s, up(s), up(w)], left);
    polygon(ctx, [s, e, up(e), up(s)], right);
    polygon(ctx, [up(n), up(e), up(s), up(w)], top);
  }

  function shade(color: string, unseen: boolean): string {
    return unseen ? darken(color, 0.3) : color;
  }

  function draw(..._deps: unknown[]) {
    if (!canvas || !colors.floor) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;
    const dpr = (window.devicePixelRatio || 1) * zoom;
    canvas.width = width * dpr;
    canvas.height = canvasHeight * dpr;
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.clearRect(0, 0, width, canvasHeight);
    ctx.imageSmoothingEnabled = false;
    hitAreas = [];
    const tileCount = map.sizeX * map.sizeY;

    // units by tile, large units are drawn at their lower right tile so that
    // they are not hidden by the floors of the other 3 tiles
    const unitsAt = new Map<string, BattleUnit[]>();
    for (const u of map.units) {
      if (!showDead && !u.alive) continue;
      const off = u.large ? 1 : 0;
      const key = `${u.x + off},${u.y + off},${u.z}`;
      unitsAt.set(key, [...(unitsAt.get(key) ?? []), u]);
    }

    // from the bottom level up to the selected level, back to front
    for (let z = map.sizeZ - 1; z >= level; z--) {
      const h = height(z);
      for (let d = 0; d <= map.sizeX + map.sizeY - 2; d++) {
        for (let x = Math.max(0, d - map.sizeY + 1); x <= Math.min(d, map.sizeX - 1); x++) {
          const y = d - x;
          const t = tileAt(map, x, y, z);
          const unseen = fogOfWar && !(t & TILE_SEEN);
          if (terrain && sheets) {
            // floor, west wall, north wall and object with the game sprites
            const base = (z * tileCount + y * map.sizeX + x) * 4;
            const [nx, ny] = project(x, y, h);
            for (let p = 0; p < 4; p++) {
              const part = terrain.parts[terrain.tiles[base + p]];
              if (!terrain.tiles[base + p] || !part) continue;
              const sheet = sheets[part.set];
              ctx.drawImage(unseen ? sheet.dark : sheet.normal,
                part.frame * SPRITE_W, 0, SPRITE_W, SPRITE_H,
                nx - SPRITE_W / 2, ny - SPRITE_FLOOR_Y - part.y, SPRITE_W, SPRITE_H);
            }
            for (const u of unitsAt.get(`${x},${y},${z}`) ?? []) drawUnit(ctx, u, h);
            continue;
          }
          if (t & TILE_GROUND) {
            const floor = (x + y) % 2 ? colors.floor : colors.floorDark;
            polygon(ctx, [project(x, y, h), project(x + 1, y, h), project(x + 1, y + 1, h), project(x, y + 1, h)],
              shade(floor, unseen));
          }
          const wallTop = shade(colors.wallTop, unseen);
          if (t & TILE_NORTH_WALL) {
            block(ctx, x, y, x + 1, y + 0.12, h, LH, wallTop, shade(colors.wallLeft, unseen), shade(colors.wallRight, unseen));
          }
          if (t & TILE_WEST_WALL) {
            block(ctx, x, y, x + 0.12, y + 1, h, LH, wallTop, shade(colors.wallLeft, unseen), shade(colors.wallRight, unseen));
          }
          if (t & TILE_OBJECT) {
            block(ctx, x + 0.2, y + 0.2, x + 0.8, y + 0.8, h, LH * 0.5,
              shade(colors.objectTop, unseen), shade(colors.object, unseen), shade(colors.floorDark, unseen));
          }
          for (const u of unitsAt.get(`${x},${y},${z}`) ?? []) drawUnit(ctx, u, h);
        }
      }
    }

    // outlines of the crafts on the ground level
    ctx.setLineDash([6, 4]);
    ctx.lineWidth = 2;
    for (const c of terrain && sheets ? [] : map.crafts) {
      const pts = [project(c.x, c.y, 0), project(c.x + c.width, c.y, 0),
                   project(c.x + c.width, c.y + c.height, 0), project(c.x, c.y + c.height, 0)];
      ctx.beginPath();
      pts.forEach(([px, py], i) => (i ? ctx.lineTo(px, py) : ctx.moveTo(px, py)));
      ctx.closePath();
      ctx.strokeStyle = c.alien ? colors.alien : colors.xcom;
      ctx.stroke();
      ctx.fillStyle = ctx.strokeStyle;
      ctx.font = '12px ' + getComputedStyle(canvas).getPropertyValue('--font-text');
      ctx.fillText(`${c.alien ? 'UFO' : 'X-COM'}: ${c.name}`, pts[0][0] - 30, pts[0][1] - 6);
    }
    ctx.setLineDash([]);
  }

  function drawUnit(ctx: CanvasRenderingContext2D, u: BattleUnit, h: number) {
    const size = u.large ? 2 : 1;
    // center of the unit on the floor
    const [cx, cy] = project(u.x + size / 2, u.y + size / 2, h);
    const color = colors[u.owner] ?? colors.xcom;
    const isSelected = selected?.index === u.index;
    const r = size * 6;

    // shadow
    ctx.beginPath();
    ctx.ellipse(cx, cy, r + 2, (r + 2) / 2, 0, 0, Math.PI * 2);
    ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
    ctx.fill();

    if (!u.alive || u.unconscious) {
      // lying on the floor
      ctx.beginPath();
      ctx.ellipse(cx, cy - 2, r + 2, r / 2, 0, 0, Math.PI * 2);
      ctx.fillStyle = u.alive ? colors.bg : color;
      ctx.strokeStyle = isSelected ? colors.selected : color;
      ctx.lineWidth = 2;
      ctx.fill();
      ctx.stroke();
      if (!u.alive) {
        ctx.beginPath();
        ctx.moveTo(cx - 4, cy - 6); ctx.lineTo(cx + 4, cy + 2);
        ctx.moveTo(cx + 4, cy - 6); ctx.lineTo(cx - 4, cy + 2);
        ctx.strokeStyle = colors.bg;
        ctx.stroke();
      }
      hitAreas.push({ unit: u, x: cx - r - 2, y: cy - r, w: 2 * r + 4, h: r + 4 });
      return;
    }

    // standing pillar with a head
    const ph = size * 18;
    ctx.fillStyle = color;
    ctx.strokeStyle = isSelected ? colors.selected : colors.bg;
    ctx.lineWidth = isSelected ? 3 : 1.5;
    ctx.beginPath();
    ctx.moveTo(cx - r, cy);
    ctx.lineTo(cx - r, cy - ph + r);
    ctx.arc(cx, cy - ph + r, r, Math.PI, 0);
    ctx.lineTo(cx + r, cy);
    ctx.ellipse(cx, cy, r, r / 2, 0, 0, Math.PI);
    ctx.closePath();
    ctx.fill();
    ctx.stroke();
    // visor, to tell the units apart from walls
    ctx.fillStyle = colors.bg;
    ctx.fillRect(cx - r * 0.6, cy - ph + r * 0.8, r * 1.2, Math.max(2, r * 0.35));
    hitAreas.push({ unit: u, x: cx - r, y: cy - ph, w: 2 * r, h: ph + r / 2 });
  }

  function unitAt(event: MouseEvent): BattleUnit | null {
    const rect = canvas.getBoundingClientRect();
    const sx = ((event.clientX - rect.left) * width) / rect.width;
    const sy = ((event.clientY - rect.top) * canvasHeight) / rect.height;
    // the last drawn unit is in front
    for (let i = hitAreas.length - 1; i >= 0; i--) {
      const a = hitAreas[i];
      if (sx >= a.x && sx <= a.x + a.w && sy >= a.y && sy <= a.y + a.h) return a.unit;
    }
    return null;
  }

  function onMove(event: MouseEvent) {
    const unit = unitAt(event);
    canvas.style.cursor = unit ? 'pointer' : 'default';
    dispatch('hover', { event, unit });
  }

  function onClick(event: MouseEvent) {
    const unit = unitAt(event);
    if (unit) dispatch('select', unit);
  }
</script>

<!-- svelte-ignore a11y-no-static-element-interactions a11y-click-events-have-key-events -->
<div class="scroller" bind:this={scroller}>
<canvas
  bind:this={canvas}
  style="width: {width * zoom}px; height: {canvasHeight * zoom}px"
  on:mousemove={onMove}
  on:mouseleave={e => dispatch('hover', { event: e, unit: null })}
  on:click={onClick}
  aria-label="Isometric mission map"
></canvas>
</div>

<style>
  .scroller {
    max-height: 70vh;
    overflow: auto;
    background: var(--bg);
    border: 1px solid var(--border);
  }
  canvas {
    display: block;
    max-width: none;
  }
</style>
