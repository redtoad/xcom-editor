# CLAUDE.md

## Branching

- When starting a new task, always create a new branch with prefix `feature/` for new features and `fix/` for bugfixes.

## Commits

- Commit often.
- Never force-push.

## Testing

- Always make sure that tests run successfully before considering a task complete.
- Loading of ALL savegame files must be tested with a round-trip (load and save identical to original).
- Always update tests and docs when making changes.
- Make sure that all examples (in `examples/`) and CLI tools (in `cmd/`, each in its own subdirectory) still run after making changes.
