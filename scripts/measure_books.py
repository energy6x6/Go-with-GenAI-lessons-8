#!/usr/bin/env python3
"""Capture real endpoint timings and SQL logs with no third-party dependencies."""
import argparse
import json
import math
import statistics
import time
import urllib.request
from pathlib import Path

parser = argparse.ArgumentParser()
parser.add_argument("--label", choices=["before", "after"], required=True)
parser.add_argument("--url", default="http://localhost:8080/api/v1/books?limit=100")
parser.add_argument("--log", type=Path, default=Path("app.log"))
parser.add_argument("--rounds", type=int, default=11)
args = parser.parse_args()
if args.rounds < 1:
    parser.error("--rounds must be positive")
if not args.log.is_file():
    parser.error("server log missing; start the API with 2>&1 | tee app.log")

with urllib.request.urlopen(args.url, timeout=10) as response:
    json.load(response)  # Warm the connection pool before measuring.

timings, counts, sql_lines = [], [], []
for _ in range(args.rounds):
    offset = args.log.stat().st_size
    start = time.perf_counter()
    with urllib.request.urlopen(args.url, timeout=10) as response:
        books = json.load(response)
    timings.append((time.perf_counter() - start) * 1000)
    with args.log.open("rb") as log:
        log.seek(offset)
        lines = log.read().decode("utf-8", errors="replace").splitlines()
    sql_lines = [line for line in lines if "SELECT " in line]
    counts.append(len(sql_lines))

ordered = sorted(timings)
result = {
    "label": args.label,
    "url": args.url,
    "books": len(books),
    "rounds": args.rounds,
    "selects_per_request": counts,
    "median_ms": round(statistics.median(timings), 3),
    "p95_ms": round(ordered[math.ceil(len(ordered) * 0.95) - 1], 3),
    "latencies_ms": [round(t, 3) for t in timings],
}
output = Path("evidence")
output.mkdir(exist_ok=True)
(output / f"n1-{args.label}.json").write_text(json.dumps(result, indent=2) + "\n")
(output / f"n1-{args.label}.log").write_text("\n".join(sql_lines) + "\n")
print(json.dumps(result, indent=2))
print(f"SQL excerpt: evidence/n1-{args.label}.log")
if not counts or any(n == 0 for n in counts):
    raise SystemExit("No SELECTs captured: check logger level, log path and concurrent requests.")
