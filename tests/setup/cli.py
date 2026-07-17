"""AxisML Lite environment lifecycle entrypoints."""

from __future__ import annotations

import argparse

from setup import lite


def setup_main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(prog="test-setup", description="Bring up a test environment")
    parser.parse_args(argv)
    return lite.setup()


def teardown_main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(prog="test-teardown", description="Tear down a test environment")
    parser.add_argument("--clean", action="store_true", help="lite: also remove the data volumes")
    args = parser.parse_args(argv)
    return lite.teardown(clean=args.clean)
