# Tick

**Tick** is a simple CLI utility for tracking time spent on different activities.

## Installation

Install the latest version using Go:

```bash
go install github.com/your-username/tick@latest
```

After installation, you can run:

```bash
tick
```

## Usage

Start tracking an activity:

```bash
tick start drawing
```

Stop the current session:

```bash
tick stop
```

Check the current activity:

```bash
tick status
```

View your statistics:

```bash
tick stats
```

## Data

Tick stores your current activity and session history locally in a JSON file:

```text
.tick/tickdata.json
```

No database or external server is required — all data is stored locally.

## Goal

Tick is designed to be a small and simple tool for tracking and managing your time directly from the terminal.

## Development

This project was created with the assistance of an AI coding agent.

The project is intentionally marked as AI-assisted to distinguish it from projects written manually.
