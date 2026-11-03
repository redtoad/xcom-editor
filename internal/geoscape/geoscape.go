// Package geoscape
//
// These are savegame files used by the Geoscape portion of the game.
// This can be considered a *standard* save.
//
// All 27 files of a Geoscape savegame (GAME_1 ... GAME_10) are supported:
// ACTS.DAT, AKNOW.DAT, ALIEN.DAT, ASTORE.DAT, BASE.DAT, BPROD.DAT, CRAFT.DAT,
// DIPLOM.DAT, FACIL.DAT, IGLOB.DAT, INTER.DAT, LEASE.DAT, LIGLOB.DAT, LOC.DAT,
// MISSIONS.DAT, PRODUCT.DAT, PROJECT.DAT, PURCHASE.DAT, RESEARCH.DAT,
// SAVEINFO.DAT, SOLDIER.DAT, TRANSFER.DAT, UIGLOB.DAT, UP.DAT, XBASES.DAT,
// XCOM.DAT and ZONAL.DAT.
//
// GEOSCAPE.EXE saves all files at 0x146c0 and loads them at 0x15460 (linear
// addresses in the LE image). It allocates the buffers at 0x214e0 and
// initializes a new game at 0x21880; the initial tables of FACIL.DAT,
// PRODUCT.DAT, PURCHASE.DAT, RESEARCH.DAT, UP.DAT and ZONAL.DAT are
// stored in the executable.
//
// Undocumented parts of a file (e.g. all of AKNOW.DAT) are kept as raw bytes,
// so that every file can be saved byte-identical to the original.
//
// Savegames saved during a mission contain additional files, see package
// battlescape. Temporary files in MISSDAT are not supported.
//
// Real savegames used for round-trip tests are stored in testdata/ at the
// repository root.
package geoscape
