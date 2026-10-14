#!/usr/bin/env python3
"""PTY smoke: gcli collects a *local* option value with `CliOpt.Question`.

Runs the example app (./_examples/cliapp) in a pseudo terminal, covering what unit
tests cannot: the question is really printed on the terminal, what the user types
becomes the option value, a value given on the command line skips the question, and
`-h` never asks.

  go build -o tmp/smoke/cliapp ./_examples/cliapp   # add GOOS/GOARCH for docker
  python3 -I scripts/ptysmoke/ask_demo.py --bin tmp/smoke/cliapp

Exit status 0 when every check passes.
"""
import argparse
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from ptydrive import Checker, PtySession  # noqa: E402

ASK = r"Please input your access token:"


def scenario(binary, env, cwd, ck):
    # 1. empty value: ask once and use what the user typed
    s = PtySession([binary, "ask-demo"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("empty value: question asked", s.until(ASK, m))
        s.line("smoke-token")
        ck.check("empty value: exit code 0", s.wait_exit() == 0)
        ck.check("empty value: typed text is used", s.contains(r"smoke-token", m))
        ck.check("empty value: asked only once", s.count(ASK, m) == 1)
    finally:
        s.close(grace=0)

    # 2. value on the command line: no question
    s = PtySession([binary, "ask-demo", "--token=cli-token"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("cli value: exit code 0", s.wait_exit() == 0)
        ck.check("cli value: no question", not s.contains(ASK, m))
        ck.check("cli value: value is used", s.contains(r"cli-token", m))
    finally:
        s.close(grace=0)

    # 3. -h: help only, never collects
    s = PtySession([binary, "ask-demo", "-h"], env=env, cwd=cwd, checker=ck)
    try:
        m = s.mark()
        ck.check("help: exit code 0", s.wait_exit() == 0)
        ck.check("help: option listed", s.contains(r"--token", m))
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
