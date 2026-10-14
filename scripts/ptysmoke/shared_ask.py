#!/usr/bin/env python3
"""PTY smoke: interactive collect of a *shared* (inherited) option, deferred to
the executing command (`Command.prepareSharedOpts`).

Covers the real-terminal behaviour unit tests cannot: a shared option bound with
`SharedOpts()` + `gflag.WithQuestion` asks **once**, on the command that actually
runs, and never asks when the value was given in any command segment (parent or
leaf) — asking in an intermediate segment would ignore a value given further right
and prompt even though the user passed it.

  go build -o tmp/smoke/cliapp ./_examples/cliapp   # add GOOS/GOARCH for docker
  python3 -I scripts/ptysmoke/shared_ask.py --bin tmp/smoke/cliapp

Exit status 0 when every check passes.
"""
import argparse
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from ptydrive import Checker, PtySession  # noqa: E402

ASK = r"Please input your api token:"


def scenario(binary, env, cwd, ck):
    # 1. empty value: exactly one question, at the executed (leaf) command
    s = PtySession([binary, "shared-ask", "run"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("empty value: question asked", s.until(ASK, m))
        s.line("shared-token")
        ck.check("empty value: exit code 0", s.wait_exit() == 0)
        ck.check("empty value: leaf uses collected value", s.contains(r"leaf: token=shared-token", m))
        ck.check("empty value: asked only once", s.count(ASK, m) == 1)
    finally:
        s.close(grace=0)

    # 2. value in the leaf segment: no question at all
    s = PtySession([binary, "shared-ask", "run", "--token=leaf-token"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("leaf value: exit code 0", s.wait_exit() == 0)
        ck.check("leaf value: no question", not s.contains(ASK, m))
        ck.check("leaf value: leaf uses given value", s.contains(r"leaf: token=leaf-token", m))
    finally:
        s.close(grace=0)

    # 3. value in the parent segment: no question either
    s = PtySession([binary, "shared-ask", "--token=parent-token", "run"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("parent value: exit code 0", s.wait_exit() == 0)
        ck.check("parent value: no question", not s.contains(ASK, m))
        ck.check("parent value: leaf uses given value", s.contains(r"leaf: token=parent-token", m))
    finally:
        s.close(grace=0)

    # 4. -h on the leaf: help shows the inherited option and never collects
    s = PtySession([binary, "shared-ask", "run", "-h"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("help: exit code 0", s.wait_exit() == 0)
        ck.check("help: inherited option listed", s.contains(r"--token", m))
        ck.check("help: no question", not s.contains(ASK, m))
    finally:
        s.close(grace=0)


def main():
    p = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    p.add_argument("--bin", default="tmp/smoke/cliapp", help="cliapp binary (built from ./_examples/cliapp)")
    p.add_argument("--cwd", default=".", help="working dir of the tested process")
    args = p.parse_args()

    binary = os.path.abspath(args.bin)
    if not os.path.exists(binary):
        raise SystemExit(f"binary not found: {binary}\nbuild it first: go build -o {args.bin} ./_examples/cliapp")

    cwd = os.path.abspath(args.cwd)
    env = dict(os.environ, TERM="xterm", NO_COLOR="1")
    ck = Checker()
    scenario(binary, env, cwd, ck)
    sys.exit(ck.report())


if __name__ == "__main__":
    main()
