#!/usr/bin/env python3

from __future__ import annotations

import html
import json
import re
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[1]
DOCS_DIR = REPO_ROOT / "docs"
RADAR_DIR = DOCS_DIR / "radar"


def read_text(path: Path) -> str:
    if not path.exists():
        return ""
    return path.read_text(encoding="utf-8")


def inline_markdown_to_html(text: str) -> str:
    escaped = html.escape(text, quote=False)
    escaped = re.sub(
        r"\[([^\]]+)\]\(([^)]+)\)",
        lambda m: f'<a href="{html.escape(m.group(2), quote=True)}">{m.group(1)}</a>',
        escaped,
    )
    escaped = re.sub(r"`([^`]+)`", r"<code>\1</code>", escaped)
    escaped = re.sub(r"\*\*([^*]+)\*\*", r"<strong>\1</strong>", escaped)
    return escaped


def get_title(text: str, fallback: str) -> str:
    for line in text.splitlines():
        if line.startswith("# "):
            return line[2:].strip()
    return fallback


def get_metadata_line(text: str) -> str:
    for prefix in ("Coverage window:", "Window:", "Source:"):
        for line in text.splitlines():
            if line.startswith(prefix):
                return line.strip()
    return ""


def get_section_bullets(text: str, heading: str, limit: int = 3) -> list[str]:
    lines = text.splitlines()
    bullets: list[str] = []
    in_section = False
    for line in lines:
        if line.strip() == f"## {heading}":
            in_section = True
            continue
        if in_section and line.startswith("## "):
            break
        if in_section and line.startswith("- "):
            bullets.append(line[2:].strip())
            if len(bullets) >= limit:
                break
    return bullets


def get_section_paragraph(text: str, heading: str) -> str:
    lines = text.splitlines()
    in_section = False
    collected: list[str] = []
    for line in lines:
        if line.strip() == f"## {heading}":
            in_section = True
            continue
        if in_section and line.startswith("## "):
            break
        if in_section:
            if line.strip():
                collected.append(line.strip())
            elif collected:
                break
    return " ".join(collected)


def get_proposal_titles(text: str, limit: int = 3) -> list[str]:
    lines = text.splitlines()
    titles: list[str] = []
    for idx, line in enumerate(lines):
        if line.strip() == "### Title":
            for next_line in lines[idx + 1 :]:
                if next_line.strip():
                    titles.append(next_line.strip())
                    break
        if len(titles) >= limit:
            break
    return titles


def recent_files(directory: Path, limit: int = 6) -> list[Path]:
    if not directory.exists():
        return []
    files = [
        path
        for path in directory.glob("*.md")
        if path.is_file() and path.name != "latest.md"
    ]
    return sorted(files, reverse=True)[:limit]


def count_files(directory: Path) -> int:
    if not directory.exists():
        return 0
    return sum(
        1
        for path in directory.glob("*.md")
        if path.is_file() and path.name != "latest.md"
    )


def rel_path(path: Path) -> str:
    return path.relative_to(DOCS_DIR).as_posix()


def card_html(
    label: str,
    title: str,
    latest_href: str,
    meta: str,
    bullets: list[str],
    archive_href: str,
    accent: str,
) -> str:
    bullet_items = (
        "".join(f"<li>{inline_markdown_to_html(item)}</li>" for item in bullets)
        or "<li class=\"empty-item\">No extracted highlights yet. Open the latest note for the full raw view.</li>"
    )
    return f"""
    <article class="brief-card {accent}">
      <div class="eyebrow">{html.escape(label)}</div>
      <h3>{html.escape(title)}</h3>
      <p class="meta">{html.escape(meta)}</p>
      <ul>{bullet_items}</ul>
      <div class="card-actions">
        <a class="primary-link" href="{html.escape(latest_href)}">Open latest</a>
        <a class="secondary-link" href="{html.escape(archive_href)}">Browse archive</a>
      </div>
    </article>
    """


def archive_column_html(title: str, entries: list[dict[str, str]]) -> str:
    items = "".join(
        f'<li><a href="{html.escape(entry["path"])}">{html.escape(entry["name"])}</a></li>'
        for entry in entries
    ) or '<li class="empty-item">No archived notes yet.</li>'
    return f"""
    <section class="archive-column">
      <h4>{html.escape(title)}</h4>
      <ul>{items}</ul>
    </section>
    """


def recent_entries(directory: Path, limit: int = 6) -> list[dict[str, str]]:
    return [{"name": path.stem, "path": rel_path(path)} for path in recent_files(directory, limit)]


def stream_data(
    label: str,
    latest_path: Path,
    text: str,
    archive_directory: Path,
    archive_href: str,
    bullet_heading: str,
    bullet_limit: int,
    fallback_title: str,
    fallback_meta: str,
) -> dict[str, object]:
    return {
        "label": label,
        "title": get_title(text, fallback_title),
        "meta": get_metadata_line(text) or fallback_meta,
        "highlights": get_section_bullets(text, bullet_heading, bullet_limit),
        "latest_path": rel_path(latest_path),
        "archive_path": archive_href,
        "recent": recent_entries(archive_directory),
        "count": count_files(archive_directory),
    }


def build_surface_data() -> dict[str, object]:
    weekly_latest = RADAR_DIR / "weekly" / "latest.md"
    daily_latest = RADAR_DIR / "daily" / "latest.md"
    trending_latest = RADAR_DIR / "trending" / "latest.md"

    weekly_text = read_text(weekly_latest)
    daily_text = read_text(daily_latest)
    trending_text = read_text(trending_latest)
    weekly = stream_data(
        label="Weekly",
        latest_path=weekly_latest,
        text=weekly_text,
        archive_directory=RADAR_DIR / "weekly",
        archive_href="radar/weekly/",
        bullet_heading="Field Summary",
        bullet_limit=3,
        fallback_title="Weekly synthesis unavailable",
        fallback_meta="Latest durable weekly synthesis",
    )
    daily = stream_data(
        label="Daily",
        latest_path=daily_latest,
        text=daily_text,
        archive_directory=RADAR_DIR / "daily",
        archive_href="radar/daily/",
        bullet_heading="New Signals",
        bullet_limit=3,
        fallback_title="Daily radar unavailable",
        fallback_meta="Latest daily field signal intake",
    )
    trending = stream_data(
        label="Trending",
        latest_path=trending_latest,
        text=trending_text,
        archive_directory=RADAR_DIR / "trending",
        archive_href="radar/trending/",
        bullet_heading="Notable Projects",
        bullet_limit=3,
        fallback_title="Trending snapshot unavailable",
        fallback_meta="Latest GitHub Trending board summary",
    )
    deep_dive = get_section_paragraph(weekly_text, "Theme Deep Dive") or (
        "No weekly deep dive is available yet. Once the first synthesis lands, "
        "this block should summarize the durable interpretation rather than just the raw signal list."
    )
    proposal_titles = get_proposal_titles(weekly_text, 3)
    flow = [
        {"name": "Radar", "description": "Collect authoritative external signal and keep noise bounded."},
        {"name": "Feedback", "description": "Normalize real Harness usage pressure into repeatable cases."},
        {"name": "Knowledge", "description": "Connect source material, usage pressure, and internal gaps into durable understanding."},
        {"name": "Proposals", "description": "Turn mature evidence into selective, reviewable change candidates."},
        {"name": "Implementation", "description": "Only reviewed proposals should feed Harness, Skills, or eval work."},
    ]
    layers = [
        {"name": "Radar", "href": "radar/README.md", "description": "External first-party signal, recurring notes, source policy, and the freshest field view."},
        {"name": "Feedback", "href": "feedback/README.md", "description": "Structured Harness usage cases instead of raw session dumping."},
        {"name": "Knowledge", "href": "knowledge/README.md", "description": "Durable synthesis, timelines, concept maps, and topic continuity."},
        {"name": "Proposals", "href": "proposals/README.md", "description": "Selective reviewed change candidates after evidence is mature enough."},
        {"name": "Research", "href": "research/README.md", "description": "One-off analysis and incubation, kept useful but out of the active lane."},
    ]
    total_notes = int(weekly["count"]) + int(daily["count"]) + int(trending["count"])
    return {
        "hero_title": "Readable field intelligence instead of disappearing thread residue.",
        "intro": (
            "This page is the human-facing front door for the repository's knowledge system. "
            "It turns the latest radar outputs into a stable briefing surface, keeps the five-layer "
            "operating model visible, and makes historical notes discoverable without digging through archived Codex threads."
        ),
        "counts": {
            "total_notes": total_notes,
            "stream_count": 3,
            "layer_count": 5,
        },
        "streams": {
            "weekly": weekly,
            "daily": daily,
            "trending": trending,
        },
        "deep_dive": deep_dive,
        "proposal_candidates": proposal_titles,
        "flow": flow,
        "layers": layers,
        "canonical_paths": {
            "visual_surface": "index.html",
            "briefing": "knowledge/briefing.md",
            "machine_index": "knowledge/index.json",
            "operating_model": "knowledge/operating-model.md",
            "learning_loop": "harness-learning-loop.md",
        },
    }


def build_page(data: dict[str, object]) -> str:
    streams = data["streams"]
    weekly = streams["weekly"]
    daily = streams["daily"]
    trending = streams["trending"]
    counts = data["counts"]
    proposal_titles = data["proposal_candidates"]
    proposal_chips = "".join(f"<li>{inline_markdown_to_html(title)}</li>" for title in proposal_titles)
    flow_steps = "".join(
        f'<div class="flow-step"><strong>{html.escape(step["name"])}</strong><p>{html.escape(step["description"])}</p></div>'
        for step in data["flow"]
    )
    layer_cards = "".join(
        f'<article class="layer-card"><h3><a href="{html.escape(layer["href"])}">{html.escape(layer["name"])}</a></h3><p>{html.escape(layer["description"])}</p></article>'
        for layer in data["layers"]
    )
    return f"""<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Erwin Knowledge Surface</title>
  <style>
    :root {{
      --paper: #f7f1e7;
      --paper-deep: #efe4d4;
      --ink: #1f2430;
      --muted: #586172;
      --line: rgba(31, 36, 48, 0.12);
      --accent: #b85d2b;
      --shadow: 0 24px 70px rgba(31, 36, 48, 0.08);
    }}
    * {{ box-sizing: border-box; }}
    html {{ scroll-behavior: smooth; }}
    body {{
      margin: 0;
      color: var(--ink);
      background:
        radial-gradient(circle at top left, rgba(184, 93, 43, 0.12), transparent 30%),
        radial-gradient(circle at 85% 10%, rgba(33, 91, 107, 0.14), transparent 26%),
        linear-gradient(180deg, #fbf6ee 0%, var(--paper) 48%, #f3ebdd 100%);
      font-family: "Avenir Next", "Segoe UI", "Helvetica Neue", sans-serif;
    }}
    a {{ color: inherit; text-decoration: none; }}
    code {{
      font-family: "SF Mono", "JetBrains Mono", "Menlo", monospace;
      background: rgba(31, 36, 48, 0.06);
      padding: 0.12rem 0.35rem;
      border-radius: 0.45rem;
      font-size: 0.95em;
    }}
    .shell {{
      width: min(1280px, calc(100vw - 32px));
      margin: 0 auto;
      padding: 32px 0 56px;
    }}
    .masthead {{
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      margin-bottom: 18px;
      padding: 14px 18px;
      border: 1px solid rgba(31, 36, 48, 0.08);
      border-radius: 999px;
      background: rgba(255,255,255,0.68);
      backdrop-filter: blur(12px);
      box-shadow: 0 14px 38px rgba(31, 36, 48, 0.06);
    }}
    .masthead strong {{
      font-size: 0.95rem;
      letter-spacing: 0.02em;
    }}
    .masthead nav {{
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
    }}
    .nav-pill {{
      padding: 10px 14px;
      border-radius: 999px;
      font-size: 0.88rem;
      border: 1px solid rgba(31, 36, 48, 0.08);
      color: var(--muted);
      background: rgba(255,255,255,0.72);
      transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
    }}
    .nav-pill:hover {{
      transform: translateY(-1px);
      background: rgba(255,255,255,0.92);
      box-shadow: 0 10px 24px rgba(31, 36, 48, 0.08);
    }}
    .hero {{
      position: relative;
      overflow: hidden;
      border: 1px solid var(--line);
      border-radius: 28px;
      padding: 32px;
      background:
        linear-gradient(135deg, rgba(255,255,255,0.76), rgba(255,255,255,0.54)),
        linear-gradient(135deg, rgba(184, 93, 43, 0.10), rgba(33, 91, 107, 0.10));
      box-shadow: var(--shadow);
    }}
    .hero::after {{
      content: "";
      position: absolute;
      inset: auto -80px -80px auto;
      width: 240px;
      height: 240px;
      border-radius: 999px;
      background: radial-gradient(circle, rgba(184, 93, 43, 0.22), transparent 70%);
      pointer-events: none;
    }}
    .eyebrow {{
      display: inline-flex;
      align-items: center;
      gap: 8px;
      padding: 8px 12px;
      border-radius: 999px;
      font-size: 12px;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      background: rgba(255,255,255,0.65);
      border: 1px solid rgba(31, 36, 48, 0.08);
      color: var(--muted);
    }}
    h1, h2, h3, h4 {{
      font-family: "Iowan Old Style", "Palatino Linotype", "Book Antiqua", Georgia, serif;
      margin: 0;
      letter-spacing: -0.03em;
    }}
    h1 {{
      font-size: clamp(2.2rem, 5vw, 4.3rem);
      line-height: 0.95;
      margin-top: 16px;
      max-width: 12ch;
    }}
    .hero p {{
      max-width: 66ch;
      margin: 18px 0 0;
      color: var(--muted);
      font-size: 1.02rem;
      line-height: 1.7;
    }}
    .hero-stats {{
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      margin-top: 22px;
    }}
    .stat-chip {{
      min-width: 150px;
      padding: 14px 16px;
      border-radius: 18px;
      border: 1px solid rgba(31, 36, 48, 0.08);
      background: rgba(255,255,255,0.66);
      box-shadow: 0 14px 30px rgba(31, 36, 48, 0.05);
    }}
    .stat-chip strong {{
      display: block;
      font-size: 1.2rem;
      line-height: 1;
    }}
    .stat-chip span {{
      display: block;
      margin-top: 8px;
      color: var(--muted);
      font-size: 0.84rem;
      text-transform: uppercase;
      letter-spacing: 0.08em;
    }}
    .hero-grid {{
      display: grid;
      grid-template-columns: minmax(0, 1.2fr) minmax(280px, 0.8fr);
      gap: 24px;
      margin-top: 28px;
    }}
    .hero-panel {{
      border: 1px solid var(--line);
      border-radius: 22px;
      padding: 22px;
      background: rgba(255,255,255,0.55);
    }}
    .hero-panel p {{
      margin-top: 12px;
      font-size: 0.96rem;
      line-height: 1.65;
    }}
    .flow {{
      display: grid;
      grid-template-columns: repeat(5, minmax(0, 1fr));
      gap: 12px;
      margin-top: 20px;
    }}
    .flow-step {{
      border: 1px solid var(--line);
      border-radius: 18px;
      padding: 14px;
      min-height: 110px;
      background: rgba(255,255,255,0.62);
    }}
    .flow-step strong {{
      display: block;
      margin-bottom: 8px;
      font-size: 1rem;
    }}
    .flow-step p {{
      margin: 0;
      font-size: 0.92rem;
      color: var(--muted);
      line-height: 1.5;
    }}
    .section {{
      margin-top: 28px;
      border: 1px solid var(--line);
      border-radius: 28px;
      background: rgba(255,255,255,0.55);
      box-shadow: var(--shadow);
      padding: 28px;
      scroll-margin-top: 18px;
    }}
    .section-header {{
      display: flex;
      align-items: end;
      justify-content: space-between;
      gap: 16px;
      margin-bottom: 20px;
    }}
    .section-header p {{
      margin: 0;
      color: var(--muted);
      max-width: 60ch;
      line-height: 1.65;
    }}
    .brief-grid {{
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 18px;
    }}
    .brief-card {{
      border-radius: 24px;
      padding: 22px;
      min-height: 100%;
      border: 1px solid rgba(31, 36, 48, 0.08);
      background: linear-gradient(180deg, rgba(255,255,255,0.85), rgba(255,255,255,0.66));
      transition: transform 180ms ease, box-shadow 180ms ease;
    }}
    .brief-card.warm {{ box-shadow: 0 18px 45px rgba(184, 93, 43, 0.12); }}
    .brief-card.cool {{ box-shadow: 0 18px 45px rgba(33, 91, 107, 0.12); }}
    .brief-card.violet {{ box-shadow: 0 18px 45px rgba(123, 79, 154, 0.10); }}
    .brief-card:hover {{
      transform: translateY(-3px);
    }}
    .brief-card h3 {{
      margin-top: 10px;
      font-size: 1.65rem;
      line-height: 1.08;
    }}
    .brief-card .meta {{
      color: var(--muted);
      font-size: 0.92rem;
      margin-top: 10px;
      line-height: 1.5;
    }}
    .brief-card ul {{
      margin: 18px 0 0;
      padding-left: 18px;
      color: var(--ink);
      line-height: 1.65;
    }}
    .brief-card li + li {{ margin-top: 10px; }}
    .card-actions {{
      margin-top: 22px;
      display: flex;
      gap: 14px;
      flex-wrap: wrap;
    }}
    .primary-link, .secondary-link {{
      display: inline-flex;
      align-items: center;
      gap: 8px;
      padding: 11px 14px;
      border-radius: 999px;
      font-size: 0.92rem;
      border: 1px solid var(--line);
    }}
    .primary-link {{
      background: rgba(31, 36, 48, 0.92);
      color: #fffaf2;
      border-color: transparent;
    }}
    .secondary-link {{
      background: rgba(255,255,255,0.7);
      color: var(--muted);
    }}
    .split {{
      display: grid;
      grid-template-columns: minmax(0, 1.1fr) minmax(280px, 0.9fr);
      gap: 18px;
    }}
    .insight-card, .proposal-card {{
      border: 1px solid var(--line);
      border-radius: 22px;
      padding: 22px;
      background: rgba(255,255,255,0.72);
    }}
    .insight-card p, .proposal-card p {{
      margin: 12px 0 0;
      color: var(--muted);
      line-height: 1.72;
    }}
    .proposal-card ul {{
      margin: 16px 0 0;
      padding-left: 18px;
      line-height: 1.6;
    }}
    .archive-grid {{
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 18px;
    }}
    .archive-column {{
      border: 1px solid var(--line);
      border-radius: 22px;
      padding: 20px;
      background: rgba(255,255,255,0.7);
    }}
    .archive-column h4 {{ font-size: 1.35rem; }}
    .archive-column ul {{
      list-style: none;
      padding: 0;
      margin: 16px 0 0;
      display: grid;
      gap: 10px;
    }}
    .archive-column li {{
      padding-bottom: 10px;
      border-bottom: 1px dashed rgba(31, 36, 48, 0.08);
      color: var(--muted);
    }}
    .archive-column li:last-child {{
      border-bottom: none;
      padding-bottom: 0;
    }}
    .empty-item {{
      color: var(--muted);
      font-style: italic;
    }}
    .layer-grid {{
      display: grid;
      grid-template-columns: repeat(5, minmax(0, 1fr));
      gap: 14px;
    }}
    .layer-card {{
      border: 1px solid var(--line);
      border-radius: 20px;
      padding: 18px;
      background: rgba(255,255,255,0.74);
    }}
    .layer-card p {{
      margin-top: 10px;
      font-size: 0.93rem;
      line-height: 1.55;
      color: var(--muted);
    }}
    .footer-note {{
      margin-top: 18px;
      color: var(--muted);
      font-size: 0.92rem;
      line-height: 1.6;
    }}
    @media (max-width: 980px) {{
      .masthead {{
        border-radius: 22px;
        align-items: start;
        flex-direction: column;
      }}
      .hero-grid, .brief-grid, .split, .archive-grid, .layer-grid, .flow {{
        grid-template-columns: 1fr;
      }}
      .shell {{ width: min(100vw - 20px, 100%); }}
      .hero, .section {{ padding: 22px; }}
    }}
  </style>
</head>
<body>
  <main class="shell">
    <div class="masthead">
      <strong>Erwin Briefing Surface</strong>
      <nav>
        <a class="nav-pill" href="#briefing">Latest briefing</a>
        <a class="nav-pill" href="#interpretation">Interpretation</a>
        <a class="nav-pill" href="#architecture">Architecture</a>
        <a class="nav-pill" href="#history">History</a>
      </nav>
    </div>
    <section class="hero">
      <div class="eyebrow">Erwin Knowledge Surface</div>
      <h1>{html.escape(data["hero_title"])}</h1>
      <p>
        {html.escape(data["intro"])}
      </p>
      <div class="hero-stats">
        <div class="stat-chip"><strong>{counts["total_notes"]}</strong><span>Durable notes indexed</span></div>
        <div class="stat-chip"><strong>{counts["stream_count"]} streams</strong><span>Daily, weekly, trending</span></div>
        <div class="stat-chip"><strong>{counts["layer_count"]} layers</strong><span>Radar to implementation</span></div>
      </div>
      <div class="hero-grid">
        <div class="hero-panel">
          <div class="eyebrow">Current Lead</div>
          <h2>{html.escape(weekly["title"])}</h2>
          <p>{inline_markdown_to_html(data["deep_dive"])}</p>
        </div>
        <div class="hero-panel">
          <div class="eyebrow">At A Glance</div>
          <h2>What to open first</h2>
          <p>
            Start with <a href="{html.escape(weekly["latest_path"])}"><code>weekly/latest.md</code></a> when you need the durable weekly read, then use the archive lists below to move backward in time. This page always reflects the current <code>latest.md</code> set for each stream.
          </p>
          <div class="card-actions">
            <a class="primary-link" href="{html.escape(data["canonical_paths"]["briefing"])}">Open text briefing</a>
            <a class="secondary-link" href="{html.escape(data["canonical_paths"]["machine_index"])}">Open machine index</a>
          </div>
        </div>
      </div>
      <div class="flow">
        {flow_steps}
      </div>
    </section>

    <section class="section" id="briefing">
      <div class="section-header">
        <div>
          <div class="eyebrow">Latest Briefing</div>
          <h2>The current field read</h2>
        </div>
        <p>These cards are generated from the latest daily, weekly, and trending notes so you can scan the signal surface without opening raw markdown first.</p>
      </div>
      <div class="brief-grid">
        {card_html(weekly["label"], weekly["title"], weekly["latest_path"], weekly["meta"], weekly["highlights"], weekly["archive_path"], "warm")}
        {card_html(daily["label"], daily["title"], daily["latest_path"], daily["meta"], daily["highlights"], daily["archive_path"], "cool")}
        {card_html(trending["label"], trending["title"], trending["latest_path"], trending["meta"], trending["highlights"], trending["archive_path"], "violet")}
      </div>
    </section>

    <section class="section" id="interpretation">
      <div class="section-header">
        <div>
          <div class="eyebrow">Durable Interpretation</div>
          <h2>What the latest week is trying to teach us</h2>
        </div>
        <p>The point of the knowledge layer is not to keep more text. It is to keep a cleaner understanding of what changes, why it matters, and what deserves a proposal.</p>
      </div>
      <div class="split">
        <article class="insight-card">
          <h3>Theme Deep Dive</h3>
          <p>{inline_markdown_to_html(data["deep_dive"])}</p>
        </article>
        <article class="proposal-card">
          <h3>Current Proposal Candidates</h3>
          <p>The weekly synthesis can nominate proposal candidates, but proposal creation should remain selective and reviewable.</p>
          <ul>{proposal_chips or '<li>No proposal candidates extracted from the latest weekly synthesis yet.</li>'}</ul>
        </article>
      </div>
    </section>

    <section class="section" id="architecture">
      <div class="section-header">
        <div>
          <div class="eyebrow">Knowledge Architecture</div>
          <h2>The full system, not just Radar</h2>
        </div>
        <p>Radar is only the intake layer. These directories define how external signal, internal usage pressure, durable synthesis, and reviewed proposals fit together.</p>
      </div>
      <div class="layer-grid">
        {layer_cards}
      </div>
      <p class="footer-note">
        Canonical behavior is defined in <a href="{html.escape(data["canonical_paths"]["operating_model"])}"><code>docs/knowledge/operating-model.md</code></a> and linked back into the Harness learning loop at <a href="{html.escape(data["canonical_paths"]["learning_loop"])}"><code>docs/harness-learning-loop.md</code></a>.
      </p>
    </section>

    <section class="section" id="history">
      <div class="section-header">
        <div>
          <div class="eyebrow">Recent History</div>
          <h2>Older signal should still be visible</h2>
        </div>
        <p>This archive strip keeps non-latest notes reachable without asking you to browse hidden automation threads or raw SQLite metadata.</p>
      </div>
      <div class="archive-grid">
        {archive_column_html("Daily Radar", daily["recent"])}
        {archive_column_html("Weekly Synthesis", weekly["recent"])}
        {archive_column_html("GitHub Trending", trending["recent"])}
      </div>
    </section>
  </main>
</body>
</html>
"""


def build_briefing(data: dict[str, object]) -> str:
    streams = data["streams"]
    weekly = streams["weekly"]
    daily = streams["daily"]
    trending = streams["trending"]
    proposal_titles = data["proposal_candidates"]

    def render_bullets(items: list[str], empty: str) -> str:
        if not items:
            return f"- {empty}"
        return "\n".join(f"- {item}" for item in items)

    def render_archive(title: str, entries: list[dict[str, str]]) -> str:
        lines = [f"### {title}", ""]
        if not entries:
            lines.append("- No archived notes yet.")
            return "\n".join(lines)
        lines.extend(f"- [{entry['name']}](../{entry['path']})" for entry in entries)
        return "\n".join(lines)

    proposal_block = render_bullets(
        proposal_titles,
        "No proposal candidates extracted from the latest weekly synthesis yet.",
    )

    sections = [
        "# Current Knowledge Briefing",
        "",
        "> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.",
        "",
        "## Open First",
        "",
        f"- [Visual surface](../{data['canonical_paths']['visual_surface']})",
        f"- [Latest weekly source](../{weekly['latest_path']})",
        f"- [Machine-readable index]({Path(data['canonical_paths']['machine_index']).name})",
        f"- [Operating model]({Path(data['canonical_paths']['operating_model']).name})",
        "",
        "## At A Glance",
        "",
        f"- `{data['counts']['total_notes']}` durable notes are currently indexed across `{data['counts']['stream_count']}` recurring streams.",
        "- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.",
        "- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.",
        "",
        "## Weekly Lead",
        "",
        f"### {weekly['title']}",
        "",
        f"- Metadata: `{weekly['meta']}`",
        render_bullets(weekly["highlights"], "Open the weekly note for the full synthesis."),
        "",
        "## Daily Signal Snapshot",
        "",
        f"### {daily['title']}",
        "",
        f"- Metadata: `{daily['meta']}`",
        render_bullets(daily["highlights"], "Open the daily note for the full signal list."),
        "",
        "## Trending Snapshot",
        "",
        f"### {trending['title']}",
        "",
        f"- Metadata: `{trending['meta']}`",
        render_bullets(trending["highlights"], "Open the trending note for the full board summary."),
        "",
        "## Durable Interpretation",
        "",
        data["deep_dive"],
        "",
        "## Current Proposal Candidates",
        "",
        proposal_block,
        "",
        "## Recent History",
        "",
        render_archive("Daily Radar", daily["recent"]),
        "",
        render_archive("Weekly Synthesis", weekly["recent"]),
        "",
        render_archive("GitHub Trending", trending["recent"]),
        "",
    ]
    return "\n".join(sections).rstrip() + "\n"


def build_index_json(data: dict[str, object]) -> str:
    payload = {
        "surface_version": 1,
        "hero_title": data["hero_title"],
        "intro": data["intro"],
        "counts": data["counts"],
        "canonical_paths": data["canonical_paths"],
        "flow": data["flow"],
        "layers": data["layers"],
        "streams": data["streams"],
        "deep_dive": data["deep_dive"],
        "proposal_candidates": data["proposal_candidates"],
    }
    return json.dumps(payload, indent=2, ensure_ascii=False) + "\n"


def main() -> int:
    data = build_surface_data()
    html_output = "\n".join(line.rstrip() for line in build_page(data).splitlines()) + "\n"
    briefing_output = build_briefing(data)
    json_output = build_index_json(data)
    (DOCS_DIR / "index.html").write_text(html_output, encoding="utf-8")
    (DOCS_DIR / "knowledge" / "briefing.md").write_text(briefing_output, encoding="utf-8")
    (DOCS_DIR / "knowledge" / "index.json").write_text(json_output, encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
