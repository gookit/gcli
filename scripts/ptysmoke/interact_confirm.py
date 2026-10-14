#!/usr/bin/env python3
"""PTY smoke: cliui's line-based prompts inside a gcli app.

`builtin.GenAutoComplete` and friends ask with `interact.Unconfirmed` /
`interact.AnswerIsYes`; v0.5.1 changed those to return `(bool, error)` and broke the
build, v0.5.3 restored the legacy signatures. These cases keep the *runtime* side
covered: a real terminal, one answer per line, "n" and "y" taking the intended path.

  go build -o tmp/smoke/cliapp ./_examples/cliapp   # add GOOS/GOARCH for docker
  python3 -I scripts/ptysmoke/interact_confirm.py --bin tmp/smoke/cliapp

Exit status 0 when every check passes.
"""
import argparse
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from ptydrive import Checker, PtySession  # noqa: E402

CONFIRM = r"Ensure continue \[yes\|no\]"


def scenario(binary, env, cwd, ck):
    # 1. "n" -> not confirmed
    s = PtySession([binary, "interact", "confirm"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("confirm: question asked", s.until(CONFIRM, m))
        s.line("n")
        ck.check("confirm: exit code 0", s.wait_exit() == 0)
        ck.check("confirm: n is unconfirmed", s.contains(r"Unconfirmed", m))
    finally:
        s.close(grace=0)

    # 2. "y" -> confirmed
    s = PtySession([binary, "interact", "confirm"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("confirm: question asked again", s.until(CONFIRM, m))
        s.line("y")
        ck.check("confirm: exit code 0 (y)", s.wait_exit() == 0)
        ck.check("confirm: y is confirmed", s.contains(r"Confirmed", m))
        ck.check("confirm: no retry hint", not s.contains(r"Please try again", m))
    finally:
        s.close(grace=0)

    # 3. read a line and print it back
    s = PtySession([binary, "interact", "read"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("read: question asked", s.until(r"Your name\?", m))
        s.line("gopher")
        ck.check("read: exit code 0", s.wait_exit() == 0)
        ck.check("read: input echoed back", s.contains(r"Your input:\s+gopher", m))
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
