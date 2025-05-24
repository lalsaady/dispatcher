<p align="center"><img alt="xplr" src="./img/dsp.png" height="200" /></p>

<p align="center">
  <a href="https://github.com/lalsaady/dispatcher/blob/main/LICENSE"><img src="https://img.shields.io/github/license/lalsaady/dispatcher?color=blue" alt="License"></a>
  <a href="https://github.com/lalsaady/dispatcher/releases"><img src="https://img.shields.io/github/v/release/lalsaady/dispatcher?include_prereleases" alt="Release"></a>
  <a href="https://github.com/lalsaady/dispatcher/actions?workflow=test"><img src="https://img.shields.io/github/actions/workflow/status/lalsaady/dispatcher/test.yaml" alt="CI"></a>
</p>

## Overview

Dispatcher is a powerful delivery route optimization tool that uses K-means clustering to efficiently group nearby deliveries and assign them to drivers. It automatically organizes delivery routes to minimize travel time and maximize delivery efficiency.

## Requirements

- Go 1.24 or later
- Google Maps API key (set as `GOOGLE_MAPS_API_KEY` environment variable)

## Installation

```bash
go install github.com/lalsaady/dispatcher@latest
```

## Usage

```bash
# Build binary
% make build
Building dsp...
go build -o dsp main.go

# Add GOOGLE API KEY to env
export GOOGLE_MAPS_API_KEY=<enter your key>

# Example using long flags
% ./dsp run \
  --address "123 Main St Cleveland OH,456 Elm St Cleveland OH" \
  --driver "Alice,Bob" \
  --hub "789 Random Rd Cleveland OH"

# Example using short flags
% ./dsp run \
  -a "123 Main St Cleveland OH,456 Elm St Cleveland OH" \
  -d "Alice,Bob" \
  -u "789 Random Rd Cleveland OH"

# Example using CSV
% ./dsp run \
  -a addresses.csv \
  -d drivers.csv \
  -u "789 Random Rd Cleveland OH"
```

## How It Works

1. Addresses are converted to coordinates using Google Maps Geocoding API
2. K-means clustering groups nearby deliveries
3. Each cluster is assigned to a driver
4. Within each cluster, deliveries are ordered by distance from the hub