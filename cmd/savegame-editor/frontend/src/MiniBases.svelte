<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { baseBitURL } from './lib/api';
  import type { BaseSummary } from './lib/types';

  // MiniBases is the base selector of the base screen of the game: all 8
  // base slots side by side, each base drawn as a tiny map of its
  // facilities (built ones green, ones under construction red). Clicking a
  // base selects it.

  export let bases: BaseSummary[];
  export let selected: number | null;

  const dispatch = createEventDispatcher<{ select: number }>();

  const SLOTS = 8;
  const MINI = 14; // size of a mini base in pixels (without the border)
  const SCALE = 2;
  const WIDTH = SLOTS * (MINI + 2);
  const HEIGHT = MINI + 2;
  // background of a mini base (BASEBITS.PCK)
  const BACKGROUND_SPRITE = 41;

  // colours of the base screen palette, by palette index (shades of green for
  // built facilities, red for facilities under construction, white frame)
  const GREEN = 48;
  const RED = 32;
  const PALETTE: Record<number, string> = {
    1: '#fcfcfc',
    33: '#f06868', 34: '#e45c5c', 35: '#d85050', 37: '#c03c3c',
    49: '#b8d454', 50: '#a4c448', 51: '#90b83c', 53: '#689c28',
  };

  let canvas: HTMLCanvasElement;
  // the background is drawn in plain colour if the game files are missing
  let backgroundLoaded = false;
  const background = new Image();
  background.onload = () => { backgroundLoaded = true; };
  background.src = baseBitURL(BACKGROUND_SPRITE);

  $: bySlot = new Map(bases.map(b => [b.index, b]));
  $: if (canvas) draw(bases, selected, backgroundLoaded);

  // draw draws the mini bases like MiniBaseView of OpenXcom.
  function draw(..._deps: unknown[]) {
    const ctx = canvas.getContext('2d');
    if (!ctx) return;
    ctx.clearRect(0, 0, WIDTH, HEIGHT);
    const rect = (x: number, y: number, w: number, h: number, c: number) => {
      ctx.fillStyle = PALETTE[c];
      ctx.fillRect(x, y, w, h);
    };
    for (let i = 0; i < SLOTS; i++) {
      const ox = i * (MINI + 2);
      if (i === selected) rect(ox, 0, MINI + 2, MINI + 2, 1);
      if (backgroundLoaded) {
        ctx.drawImage(background, 0, 0, MINI + 2, MINI + 2, ox, 0, MINI + 2, MINI + 2);
      } else {
        ctx.fillStyle = '#202020';
        ctx.fillRect(ox + 1, 1, MINI, MINI);
      }
      const base = bySlot.get(i);
      if (!base) continue;
      for (const f of base.facilities ?? []) {
        const color = f.built ? GREEN : RED;
        let x = ox + 2 + f.x * 2, y = 2 + f.y * 2, w = f.size * 2, h = f.size * 2;
        rect(x, y, w, h, color + 3);
        x++; y++; w--; h--;
        rect(x, y, w, h, color + 5);
        x--; y--;
        rect(x, y, w, h, color + 2);
        x++; y++; w--; h--;
        rect(x, y, w, h, color + 3);
        x--; y--;
        rect(x, y, 1, 1, color + 1);
      }
    }
  }

  function slotAt(e: MouseEvent): number {
    const r = canvas.getBoundingClientRect();
    return Math.floor((e.clientX - r.left) / (r.width / SLOTS));
  }

  function handleClick(e: MouseEvent) {
    const slot = slotAt(e);
    if (bySlot.has(slot)) dispatch('select', slot);
  }

  let hovered: BaseSummary | undefined;
  function handleMove(e: MouseEvent) {
    hovered = bySlot.get(slotAt(e));
  }
</script>

<canvas
  bind:this={canvas}
  width={WIDTH}
  height={HEIGHT}
  style="width:{WIDTH * SCALE}px; height:{HEIGHT * SCALE}px"
  class:clickable={hovered !== undefined}
  title={hovered?.name ?? ''}
  on:click={handleClick}
  on:mousemove={handleMove}
  on:mouseleave={() => hovered = undefined}
></canvas>

<style>
  canvas {
    display: block;
    image-rendering: pixelated;
  }
  .clickable {
    cursor: pointer;
  }
</style>
