#!/usr/bin/env python3
"""Verify the Kaitai format descriptions against the savegames in testdata/.

Every .ksy file in this directory is compiled to Python with
kaitai-struct-compiler and used to parse every file of the same name found
in testdata/ (e.g. CRAFT.DAT.ksy parses testdata/*/CRAFT.DAT). A file passes
if it parses without error and all of its bytes are consumed.

Requirements:
    brew install kaitai-struct-compiler   (or see https://kaitai.io)
    pip install kaitaistruct

Usage:
    python3 misc/kaitai-structs/verify.py [NAME.DAT.ksy ...]
"""

import importlib.util
import pathlib
import subprocess
import sys
import tempfile

HERE = pathlib.Path(__file__).resolve().parent
TESTDATA = HERE.parent.parent / "testdata"


def compile_ksy(ksy: pathlib.Path, outdir: pathlib.Path) -> pathlib.Path:
    before = set(outdir.glob("*.py"))
    subprocess.run(
        ["kaitai-struct-compiler", "-t", "python", "--outdir", str(outdir), str(ksy)],
        check=True,
        capture_output=True,
        text=True,
    )
    created = set(outdir.glob("*.py")) - before
    if len(created) != 1:
        raise RuntimeError(f"expected one generated module for {ksy.name}, got {created}")
    return created.pop()


def load_class(module_path: pathlib.Path):
    spec = importlib.util.spec_from_file_location(module_path.stem, module_path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    name = "".join(part.capitalize() for part in module_path.stem.split("_"))
    return getattr(module, name)


def touch(obj, seen=None):
    """Evaluate all instances recursively so that errors in them surface."""
    from kaitaistruct import KaitaiStruct

    if seen is None:
        seen = set()
    if isinstance(obj, list):
        for item in obj:
            touch(item, seen)
        return
    if not isinstance(obj, KaitaiStruct) or id(obj) in seen:
        return
    seen.add(id(obj))
    for attr in dir(obj):
        if attr.startswith("_") or attr[0].isupper():
            continue
        value = getattr(obj, attr)
        if callable(value):
            continue
        touch(value, seen)


def main(argv):
    ksys = [HERE / a for a in argv] if argv else sorted(HERE.glob("*.ksy"))
    failures = 0
    with tempfile.TemporaryDirectory() as tmp:
        outdir = pathlib.Path(tmp)
        sys.path.insert(0, tmp)
        for ksy in ksys:
            filename = ksy.name.removesuffix(".ksy")
            try:
                cls = load_class(compile_ksy(ksy, outdir))
            except subprocess.CalledProcessError as e:
                print(f"FAIL {ksy.name}: compile error\n{e.stdout}{e.stderr}")
                failures += 1
                continue
            samples = sorted(TESTDATA.glob(f"*/{filename}"))
            if not samples:
                print(f"SKIP {ksy.name}: compiled, no sample files in testdata/")
                continue
            for sample in samples:
                rel = sample.relative_to(TESTDATA)
                try:
                    obj = cls.from_file(str(sample))
                    touch(obj)
                    size = sample.stat().st_size
                    pos = obj._io.pos()
                    if pos != size:
                        raise ValueError(f"consumed {pos} of {size} bytes")
                    print(f"OK   {ksy.name}: {rel}")
                except Exception as e:  # noqa: BLE001
                    print(f"FAIL {ksy.name}: {rel}: {type(e).__name__}: {e}")
                    failures += 1
    print(f"\n{failures} failure(s)")
    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
