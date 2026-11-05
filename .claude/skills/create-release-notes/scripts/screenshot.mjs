// Takes a screenshot of the running savegame editor.
//
// Usage:
//   node screenshot.mjs --out FILE.png [--url http://localhost:8091]
//                       [--game GAME_4] [--tab Mission] [--click "Text"]...
//                       [--element CSS-SELECTOR] [--wait MS]
//                       [--width 1440] [--height 900]
//
// --game   clicks the savegame in the sidebar (matched by its text)
// --tab    opens a tab of the savegame (Mission, Soldiers, Bases, Crafts, ...)
// --click  clicks a button or other element by its visible text (repeatable,
//          in the given order), e.g. a toggle like "Show radar range"
// --element  screenshots only this element instead of the whole page
// --wait   extra time to wait before the screenshot (maps, sprites, globe)
import { chromium } from 'playwright';

const args = { url: 'http://localhost:8091', clicks: [], wait: 1500, width: 1440, height: 900 };
const argv = process.argv.slice(2);
for (let i = 0; i < argv.length; i++) {
  const key = argv[i];
  const value = argv[++i];
  switch (key) {
    case '--url': args.url = value; break;
    case '--out': args.out = value; break;
    case '--game': args.game = value; break;
    case '--tab': args.tab = value; break;
    case '--click': args.clicks.push(value); break;
    case '--element': args.element = value; break;
    case '--wait': args.wait = Number(value); break;
    case '--width': args.width = Number(value); break;
    case '--height': args.height = Number(value); break;
    default:
      console.error(`unknown argument ${key}`);
      process.exit(2);
  }
}
if (!args.out) {
  console.error('--out is required');
  process.exit(2);
}

// Prefer an installed Google Chrome, so no browser download is needed; fall
// back to the Chromium of Playwright (`npx playwright install chromium`).
let browser;
try {
  browser = await chromium.launch({ channel: 'chrome' });
} catch {
  browser = await chromium.launch();
}
try {
  const page = await browser.newPage({
    viewport: { width: args.width, height: args.height },
    deviceScaleFactor: 2,
  });
  const errors = [];
  page.on('pageerror', (err) => errors.push(err.message));

  await page.goto(args.url, { waitUntil: 'networkidle' });
  if (args.game) {
    // the sidebar shows the title of a savegame, not its slot
    const games = await (await fetch(new URL('/api/games', args.url))).json();
    const game = games.find((g) => g.slot === args.game);
    if (!game) throw new Error(`savegame ${args.game} not found`);
    await page.locator('button.game-card', { hasText: game.title }).first().click();
    await page.waitForLoadState('networkidle');
  }
  if (args.tab) {
    await page.locator('nav.tabs button.tab', { hasText: args.tab }).first().click();
    await page.waitForLoadState('networkidle');
  }
  for (const text of args.clicks) {
    await page.getByText(text, { exact: false }).first().click();
    await page.waitForLoadState('networkidle');
  }
  await page.waitForTimeout(args.wait);

  const target = args.element ? page.locator(args.element).first() : page;
  await target.screenshot({ path: args.out });
  console.log(`saved ${args.out}`);
  if (errors.length) {
    console.error(`page errors:\n${errors.join('\n')}`);
  }
} finally {
  await browser.close();
}
