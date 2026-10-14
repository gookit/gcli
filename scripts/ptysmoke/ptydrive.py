"""Drive an interactive CLI process in a pseudo terminal (Linux/macOS only).

Same approach as suga's ``scripts/ptysmoke/ptydrive.py`` (see its README for the
origin of the lessons); on Windows a console cannot be driven from Python, so run
these cases on Linux/macOS or inside a ``python:3.12-slim`` container.

Lessons built in:
  - assertions only look at output produced after a mark, never the whole buffer;
  - wait for the expected prompt before typing the next input;
  - Enter is "\\r" on the terminal, Ctrl-C is "\\x03", Ctrl-U is "\\x15";
  - ``text()`` strips ANSI escapes, so colour-only output needs ``text_raw()``.
"""
import os
import pty
import re
import select
import time

ANSI = re.compile(r"\x1b\[[0-9;?]*[A-Za-z]")
ENTER = "\r"
CTRL_C = "\x03"
CTRL_D = "\x04"
CTRL_U = "\x15"


class Checker:
    """Collects check results across several short-lived sessions."""

    def __init__(self):
        self.results = []

    def check(self, name, ok):
        self.results.append((name, bool(ok)))
        print(("PASS " if ok else "FAIL ") + name, flush=True)
        return ok

    def failed(self):
        return [name for name, ok in self.results if not ok]

    def report(self):
        """Print the summary line, return the process exit code."""
        failed = self.failed()
        print(f"\n{len(self.results) - len(failed)}/{len(self.results)} passed")
        if failed:
            print("failed: " + ", ".join(failed))
            return 1
        return 0


class PtySession:
    def __init__(self, argv, env=None, cwd=None, checker=None):
        self.buf = b""
        self.checker = checker or Checker()
        self.pid, self.fd = pty.fork()
        if self.pid == 0:  # child
            if cwd:
                os.chdir(cwd)
            os.execve(argv[0], argv, env or dict(os.environ))

    # --- output ---------------------------------------------------------
    def pump(self, seconds):
        """Read and DROP output for seconds (keeps the child from blocking on
        a full pty). To wait and still check what came, use until."""
        end = time.time() + seconds
        while time.time() < end:
            ready, _, _ = select.select([self.fd], [], [], 0.1)
            if not ready:
                continue
            try:
                data = os.read(self.fd, 4096)
            except OSError:
                return False
            if not data:
                return False
        return True

    def _read_once(self, timeout):
        ready, _, _ = select.select([self.fd], [], [], timeout)
        if ready:
            try:
                data = os.read(self.fd, 4096)
            except OSError:
                return False
            if not data:
                return False
            self.buf += data
        return True

    def mark(self):
        """Position to search from; use it right before sending input."""
        return len(self.buf)

    def text(self, start=0):
        return ANSI.sub("", self.buf[start:].decode("utf-8", "replace"))

    def until(self, pattern, start, timeout=20.0):
        """Wait until pattern (regex) appears in the output after start."""
        end = time.time() + timeout
        while time.time() < end:
            if re.search(pattern, self.text(start)):
                return True
            if not self._read_once(0.2):
                return bool(re.search(pattern, self.text(start)))
        return False

    def contains(self, pattern, start):
        return bool(re.search(pattern, self.text(start)))

    def count(self, pattern, start=0):
        return len(re.findall(pattern, self.text(start)))

    def text_raw(self, start=0):
        """The raw output after start, escape codes kept: text() strips them, so
        it hides colour-only output such as a highlighted input line."""
        return self.buf[start:].decode("utf-8", "replace")

    def until_raw(self, pattern, start, timeout=20.0):
        end = time.time() + timeout
        while time.time() < end:
            if re.search(pattern, self.text_raw(start)):
                return True
            if not self._read_once(0.2):
                return bool(re.search(pattern, self.text_raw(start)))
        return False

    def contains_raw(self, pattern, start):
        return bool(re.search(pattern, self.text_raw(start)))

    # --- input ----------------------------------------------------------
    def send(self, s):
        os.write(self.fd, s.encode())

    def line(self, s):
        self.send(s + ENTER)

    # --- lifecycle ------------------------------------------------------
    def wait_exit(self, timeout=20.0):
        """Drain output until EOF, then return the child's exit code.

        Returns None when the child did not exit in time (the caller's close()
        then kills it), so a hanging command is reported as a failed check.
        """
        end = time.time() + timeout
        while time.time() < end:
            if not self._read_once(0.2):
                break

        end = time.time() + 5
        while time.time() < end:
            try:
                pid, status = os.waitpid(self.pid, os.WNOHANG)
            except OSError:
                return None
            if pid == self.pid:
                return os.waitstatus_to_exitcode(status)
            time.sleep(0.05)
        return None

    def close(self, grace=2.0):
        self.pump(grace)
        try:
            os.kill(self.pid, 9)
        except OSError:
            pass
        try:
            os.waitpid(self.pid, 0)
        except OSError:
            pass

    # --- checks ---------------------------------------------------------
    def check(self, name, ok):
        return self.checker.check(name, ok)

    def failed(self):
        return self.checker.failed()
