# SmartAlac
A tool for managing my music library

## History
Years ago I wrote a python tool to manage my 50k FLAC files using MusicBrainz as the primary data source, but also pulling from Amazon and Discogs. It worked well for me. It was written in Python 1. Things have changed since then and it has become unmaintainable. A rewrite is necessary.

I'd written a cd-ripping shell script (with a helper written in C) to digitize my CD collection, encode them, and tag them. I digitized all my CDs. I lost the script (and the helper). I've since bought more CDs that I need to rip.

I've learned Go. Go is good.

## Goals

Replace my curate.py script with something written in Go that will updated the tags on my ALAC files every month or so based on the data from MusicBrainz and Discogs

A tool to help with the final stages of my vinyl rips (encode, tag, move into place)

A tool to help rip CDs (rip, encode, tag, move into place)
