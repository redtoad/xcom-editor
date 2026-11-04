// Edit forms (soldier, base, financials) only send their values to the server
// when "Apply Changes" is clicked. To make sure the global "Save" button does
// not silently drop values that were edited but not applied yet, every open
// form registers itself here. "Save" calls applyPendingChanges() before
// writing the savegame to disk, and "Reload" calls refreshForms() so open
// forms show (and compare against) the reloaded values.

export interface EditForm {
  // apply sends unapplied changes to the server.
  apply: () => Promise<void>;
  // refresh reloads the form values from the server.
  refresh: () => Promise<void>;
}

const forms = new Set<EditForm>();

// registerForm registers an open form and returns a function that
// unregisters it again (use it in onDestroy).
export function registerForm(form: EditForm): () => void {
  forms.add(form);
  return () => {
    forms.delete(form);
  };
}

// applyPendingChanges applies the unapplied changes of all open forms.
export async function applyPendingChanges(): Promise<void> {
  for (const form of [...forms]) {
    await form.apply();
  }
}

// refreshForms reloads all open forms from the server.
export async function refreshForms(): Promise<void> {
  for (const form of [...forms]) {
    await form.refresh();
  }
}

// ChangeTracker sends form values to the server only if they differ from the
// values last loaded from or sent to the server. Sends are serialized: each
// apply waits for the previous one to finish before it reads the current
// values and compares them, so an in-flight update can never make a later
// apply skip values that differ from what the server ends up with.
export class ChangeTracker<T> {
  private applied = '';
  private queue: Promise<void> = Promise.resolve();

  constructor(private send: (values: T) => Promise<void>) {}

  // reset records values as being the server state (e.g. after loading).
  reset(values: T): void {
    this.applied = JSON.stringify(values);
  }

  // apply sends the values returned by current() if they changed. current is
  // called only once all previous applies have finished.
  apply(current: () => T): Promise<void> {
    const run = this.queue.then(async () => {
      const values = current();
      const snap = JSON.stringify(values);
      if (snap === this.applied) return;
      await this.send(values);
      this.applied = snap;
    });
    // Keep the queue going even if this apply fails; the caller sees the error.
    this.queue = run.catch(() => {});
    return run;
  }
}
