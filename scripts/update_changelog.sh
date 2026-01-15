#!/usr/bin/env bash
set -euo pipefail

TAG="${1:-}"
if [[ -z "$TAG" ]]; then
  echo "usage: $0 <tag>"
  exit 1
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

python3 - "$TAG" <<'PY'
import re
import subprocess
import sys
from datetime import datetime, timezone

TAG = sys.argv[1]
CHLOG = "CHANGELOG.md"

def run(cmd: list[str]) -> str:
    return subprocess.check_output(cmd, text=True).strip()

if not subprocess.run(["test", "-f", CHLOG]).returncode:
    pass

# Ensure CHANGELOG exists with base template.
try:
    with open(CHLOG, "r", encoding="utf-8") as f:
        content = f.read()
except FileNotFoundError:
    content = """# Changelog

All notable changes to this project will be documented in this file.

## Unreleased

### Breaking Changes

None.

### Added

Nothing yet.

### Changed

Nothing yet.

### Fixed

Nothing yet.
"""

def parse_semver(v: str):
    m = re.match(r"^v(\d+)\.(\d+)\.(\d+)$", v)
    if not m:
        return None
    return tuple(int(x) for x in m.groups())

tags_raw = run(["git", "tag", "-l", "v*"]).splitlines()
tags = [t for t in tags_raw if parse_semver(t) is not None]
tags.sort(key=parse_semver)

if TAG not in tags:
    raise SystemExit(f"Tag not found in local semver tags: {TAG}")

idx = tags.index(TAG)
prev = tags[idx - 1] if idx > 0 else None
range_expr = f"{prev}..{TAG}" if prev else TAG

subjects = run(["git", "log", range_expr, "--format=%s", "--no-merges"]) if range_expr else ""

breaking = []
added = []
changed = []
fixed = []
other = []

markers = re.compile(r"(^BREAKING:|^BREAKING CHANGE:|BREAKING CHANGE|BREAKING)", re.I)
def cat(msg: str):
    if not msg:
        return
    if markers.search(msg):
        breaking.append(msg)
    elif re.match(r"^feat(\(|:| )", msg, re.I):
        added.append(msg)
    elif re.match(r"^fix(\(|:| )", msg, re.I):
        fixed.append(msg)
    elif re.match(r"^(refactor|docs|ci)(\(|:| )", msg, re.I):
        changed.append(msg)
    else:
        other.append(msg)

for s in subjects.splitlines():
    cat(s.strip())

def bullets(items):
    if not items:
        return ["Nothing yet."]
    return [f"- {x}" for x in items]

section = []
ct = int(run(["git", "show", "-s", "--format=%ct", f"{TAG}^{{commit}}"]))
DATE = datetime.fromtimestamp(ct, tz=timezone.utc).date().isoformat()
section.append(f"## {TAG} - {DATE}")
section.append("")
section.append("### Breaking Changes")
section.extend(bullets(breaking))
section.append("")
section.append("### Added")
section.extend(bullets(added))
section.append("")
section.append("### Changed")
section.extend(bullets(changed))
section.append("")
section.append("### Fixed")
section.extend(bullets(fixed))
if other:
    section.append("")
    section.append("### Misc")
    section.extend([f"- {x}" for x in other])
section.append("")
section_text = "\n".join(section).rstrip() + "\n\n"

# Remove any existing section for TAG (from '## TAG' to the next '## ' header).
pattern = re.compile(rf"^##\s+{re.escape(TAG)}(?:\s+-[0-9]{{4}}-[0-9]{{2}}-[0-9]{{2}})?\s*$", re.M)
lines = content.splitlines(keepends=True)

out = []
i = 0
skip = False
while i < len(lines):
    line = lines[i]
    if not skip and re.match(rf"^##\s+{re.escape(TAG)}", line.strip()):
        skip = True
        i += 1
        while i < len(lines) and not re.match(r"^##\s+", lines[i].strip()):
            i += 1
        continue
    if skip:
        skip = False
    out.append(line)
    i += 1

content = "".join(out)

# Insert section just before the first '## Unreleased' header.
m = re.search(r"^##\s+Unreleased\s*$", content, flags=re.M)
if not m:
    # If template missing, append at end.
    content = content.rstrip() + "\n\n" + section_text
else:
    insert_at = m.start()
    content = content[:insert_at].rstrip() + "\n\n" + section_text + content[insert_at:]

with open(CHLOG, "w", encoding="utf-8") as f:
    f.write(content)

print(f"Updated {CHLOG} for {TAG}")
PY

