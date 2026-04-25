# SmartAlac
A tool for managing my music library.

## Project Structure

This repository contains multiple tools for managing ALAC music files, unified into a single module.

- **`curate`** (`cmd/curate`): Curates and tags existing ALAC files using MusicBrainz.
- **`bme`** (`cmd/bme`): Batch Music Encoder - Rips CDs, encodes them to ALAC, and tags them.

### Shared Packages
- **`pkg/sa`**: Core logic for the `curate` tool.
- **`pkg/bme`**: Core logic for the `bme` tool.
- **`pkg/mb5`**: Shared `purego` bindings for `libmusicbrainz5`.
- **`pkg/cdio`**: Shared `purego` bindings for `libcdio`, `libcdio_cdda`, and `libcdio_paranoia`.

## History
Originally multiple projects (shell/C, Python, and then separate Go modules), these have been unified for better maintainability and shared MusicBrainz integration.

## Goals

1.  **Tagging/Curation**: Update ALAC tags using MusicBrainz (and potentially Discogs) data.
2.  **Vinyl Rips**: Encode and tag high-quality vinyl rips.
3.  **CD Ripping**: A fully automated pipeline to rip, encode, and tag CDs with a TUI dashboard.

## Usage

### Curate
```bash
go run ./cmd/curate --dir /path/to/music
```

### BME
```bash
go run ./cmd/bme --dir /path/to/work/dir
```

## Vision

- Integrated TUI for matching files missing MusicBrainz tags (e.g. vinyl rips) without Picard.
- Unified tagging logic across all tools.
