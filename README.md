# SmartAlac
A tool for managing my music library

## History
Years ago I wrote a python tool to manage my 50k FLAC files using MusicBrainz as the primary data source, but also pulling from Amazon and Discogs. It worked well for me. It was written in Python 1. Things have changed since then and it has become unmaintainable. A rewrite is necessary.

I'd written a cd-ripping shell script (with a helper written in C) to digitize my CD collection, encode them, and tag them. I digitized all my CDs. I lost the script (and the helper). I've since bought more CDs that I need to rip.

I've converted my entire collection from FLAC to ALAC because Apple-Reasons.

I've learned Go. Go is good.

## Goals

Replace my curate.py script with something written in Go that will updated the tags on my ALAC files every month or so based on the data from MusicBrainz and Discogs

A tool to help with the final stages of my vinyl rips (encode, tag, move into place)

A tool to help rip CDs (rip, encode, tag, move into place)

## Current state of "curate" tool

Basic proof-of-concept functionality in place. Reading files created with ffmpeg (from flac) works. Writing files works, but munches up the custom tags (UPPERCASE). 

Need to fix mp4tag to allow lowercase/CamelCase tag names in Customs. Because MusicBrainz picard tags CamelCase.

Need to figure out why files tagged with (newer?) Picard lose their Custom Tag key/value matching.

Need to add features to gomusicbrainz to get additional fields (e.g. ISRC and others).

ASIN/Amazon stuff is non-existant and will probably remain so. It never did much for me

Discogs stuff is non-existant and is low priority, except making sure my discogs library is up-to-date.

Getting Apple iTunes artist/album ID (did not exist in batchme) might be a nice feature.

## New TUI tool visioning

I want the ability to match/tag files missing musicbrainz tags (e.g. my vinyl rips) without having to go through Picard first.

I want the ability to encode wav files (vinyl rips) from within the tool.
