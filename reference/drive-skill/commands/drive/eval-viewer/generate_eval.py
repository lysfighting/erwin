#!/usr/bin/env python3
"""Generate and serve a drive-skill mission evaluation report.

Reads eval.json (v2.3) or eval-data.json (v2.2 legacy) + evidence/ from a
mission's grading/ directory, embeds everything into a self-contained HTML
page, and serves it via HTTP or writes a static file.

Usage:
    python generate_eval.py <mission-dir> [--port PORT]
    python generate_eval.py <mission-dir> --static <output.html>
"""

import argparse
import base64
import json
import mimetypes
import os
import signal
import subprocess
import sys
import time
from functools import partial
from http.server import HTTPServer, BaseHTTPRequestHandler
from pathlib import Path

TEXT_EXTENSIONS = {
    ".txt", ".md", ".json", ".csv", ".py", ".js", ".ts", ".tsx", ".jsx",
    ".yaml", ".yml", ".xml", ".html", ".css", ".sh", ".rb", ".go", ".rs",
    ".java", ".log", ".toml",
}
IMAGE_EXTENSIONS = {".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp"}


def get_mime_type(path: Path) -> str:
    mime, _ = mimetypes.guess_type(str(path))
    return mime or "application/octet-stream"


def embed_file(path: Path) -> dict:
    ext = path.suffix.lower()
    if ext in TEXT_EXTENSIONS:
        try:
            content = path.read_text(errors="replace")
        except OSError:
            content = "(Error reading file)"
        return {"name": path.name, "type": "text", "content": content}
    elif ext in IMAGE_EXTENSIONS:
        try:
            raw = path.read_bytes()
            mime = get_mime_type(path)
            b64 = base64.b64encode(raw).decode("ascii")
        except OSError:
            return {"name": path.name, "type": "error", "content": "(Error reading file)"}
        return {"name": path.name, "type": "image", "mime": mime, "data_uri": f"data:{mime};base64,{b64}"}
    else:
        return {"name": path.name, "type": "binary", "content": f"({path.suffix} file)"}


def collect_files(directory: Path) -> list[dict]:
    files = []
    if directory.is_dir():
        for f in sorted(directory.iterdir()):
            if f.is_file() and not f.name.startswith("."):
                files.append(embed_file(f))
    return files


def load_eval_data(mission_dir: Path) -> dict:
    # Try v2.3 format first, fall back to v2.2 legacy
    for name in ("eval.json", "eval-data.json"):
        p = mission_dir / "grading" / name
        if p.exists():
            try:
                return json.loads(p.read_text())
            except (json.JSONDecodeError, OSError):
                pass
    return {}


def load_eval_reports(mission_dir: Path) -> list[dict]:
    grading = mission_dir / "grading"
    reports = []
    if grading.is_dir():
        for f in sorted(grading.iterdir()):
            if f.is_file() and f.suffix == ".md" and f.name.startswith("eval"):
                try:
                    reports.append({"name": f.name, "content": f.read_text(errors="replace")})
                except OSError:
                    pass
    return reports


def load_comparison_data(mission_dir: Path) -> dict | None:
    grading_dir = mission_dir / "grading"
    try:
        from aggregate_eval import build_comparison
        comparison = build_comparison(grading_dir)
        if comparison and comparison.get("iteration_count", 0) >= 1:
            return comparison
    except ImportError:
        pass
    for name in ("eval-comparison.json",):
        p = grading_dir / name
        if p.exists():
            try:
                return json.loads(p.read_text())
            except (json.JSONDecodeError, OSError):
                pass
    return None


def build_report_data(mission_dir: Path) -> dict:
    eval_data = load_eval_data(mission_dir)
    evidence = collect_files(mission_dir / "grading" / "evidence")
    screenshots = collect_files(mission_dir / "grading" / "screenshots")
    reports = load_eval_reports(mission_dir)
    comparison = load_comparison_data(mission_dir)

    contract = ""
    contract_path = mission_dir / "contract.md"
    if contract_path.exists():
        try:
            contract = contract_path.read_text(errors="replace")
        except OSError:
            pass

    data = {
        "mission_name": eval_data.get("mission", mission_dir.name),
        "eval_data": eval_data,
        "evidence": evidence,
        "screenshots": screenshots,
        "reports": reports,
        "contract": contract,
    }
    if comparison:
        data["comparison"] = comparison
    return data


def generate_html(data: dict) -> str:
    template_path = Path(__file__).parent / "viewer.html"
    template = template_path.read_text()
    data_json = json.dumps(data)
    return template.replace("/*__EMBEDDED_DATA__*/", f"const EMBEDDED_DATA = {data_json};")


def _kill_port(port: int) -> None:
    try:
        result = subprocess.run(
            ["lsof", "-ti", f":{port}"],
            capture_output=True, text=True, timeout=5,
        )
        for pid_str in result.stdout.strip().split("\n"):
            if pid_str.strip():
                try:
                    os.kill(int(pid_str.strip()), signal.SIGTERM)
                except (ProcessLookupError, ValueError):
                    pass
        if result.stdout.strip():
            time.sleep(0.5)
    except (subprocess.TimeoutExpired, FileNotFoundError):
        pass


class EvalHandler(BaseHTTPRequestHandler):
    def __init__(self, mission_dir: Path, *args, **kwargs):
        self.mission_dir = mission_dir
        super().__init__(*args, **kwargs)

    def do_GET(self) -> None:
        if self.path == "/" or self.path == "/index.html":
            data = build_report_data(self.mission_dir)
            html = generate_html(data)
            content = html.encode("utf-8")
            self.send_response(200)
            self.send_header("Content-Type", "text/html; charset=utf-8")
            self.send_header("Content-Length", str(len(content)))
            self.end_headers()
            self.wfile.write(content)
        else:
            self.send_error(404)

    def log_message(self, format: str, *args: object) -> None:
        pass


def main() -> None:
    parser = argparse.ArgumentParser(description="Generate drive-skill eval report (HTML)")
    parser.add_argument("mission_dir", type=Path, help="Path to mission directory")
    parser.add_argument("--port", "-p", type=int, default=3118, help="Server port (default: 3118)")
    parser.add_argument("--static", "-s", type=Path, default=None,
                        help="Write standalone HTML to this path instead of starting server")
    args = parser.parse_args()

    mission_dir = args.mission_dir.resolve()
    if not mission_dir.is_dir():
        print(f"Error: {mission_dir} is not a directory", file=sys.stderr)
        sys.exit(1)

    data = build_report_data(mission_dir)

    if args.static:
        html = generate_html(data)
        args.static.parent.mkdir(parents=True, exist_ok=True)
        args.static.write_text(html)
        print(f"\n  Static eval report written to: {args.static}\n")
        sys.exit(0)

    port = args.port
    _kill_port(port)
    handler = partial(EvalHandler, mission_dir)
    try:
        server = HTTPServer(("127.0.0.1", port), handler)
    except OSError:
        server = HTTPServer(("127.0.0.1", 0), handler)
        port = server.server_address[1]

    url = f"http://localhost:{port}"
    print(f"\n  Drive Eval Viewer")
    print(f"  ─────────────────────────────────")
    print(f"  URL:     {url}")
    print(f"  Mission: {mission_dir}")
    print(f"\n  Press Ctrl+C to stop.\n")

    try:
        import webbrowser
        webbrowser.open(url)
    except Exception:
        pass

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nStopped.")
        server.server_close()


if __name__ == "__main__":
    main()
