// Edit forms (soldier, base, financials) only send their values to the server
// when "Apply Changes" is clicked. To make sure the global "Save" button does
// not silently drop values that were edited but not applied yet, every open
// form registers a function here that applies its pending changes. "Save"
// calls applyPendingChanges() before writing the savegame to disk.

type Applier = () => Promise<void>;

const appliers = new Set<Applier>();

// registerPendingChanges registers an applier and returns a function that
// unregisters it again (use it in onDestroy).
export function registerPendingChanges(fn: Applier): () => void {
  appliers.add(fn);
  return () => {
    appliers.delete(fn);
  };
}

// applyPendingChanges runs all registered appliers.
export async function applyPendingChanges(): Promise<void> {
  for (const fn of [...appliers]) {
    await fn();
  }
}

// snapshot returns a comparable representation of form values, used to
// detect whether a form has unapplied changes.
export function snapshot(values: Record<string, unknown>): string {
  return JSON.stringify(values);
}
