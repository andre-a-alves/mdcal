# mdcal

A CLI tool for creating markdown calendars.

## Description

mdcal is a command-line utility that generates customizable markdown calendars for a specific month or an entire year. It can be run either with command-line flags or in interactive mode, allowing you to create calendars with various options such as custom week start days, calendar week numbers, and more.

## Installation

### Prerequisites

- Go 1.24 or higher

### From Source

```bash
# Clone the repository
git clone https://github.com/andre-a-alves/mdcal.git

# Navigate to the project directory
cd mdcal

# Build the project
go build

# Optionally, install the binary to your GOPATH
go install
```

### Using Go Install

```bash
go install github.com/andre-a-alves/mdcal@latest
```

## Usage

### Interactive Mode

Run mdcal without any arguments to enter interactive mode:

```bash
mdcal
```

In interactive mode, you'll be prompted to enter:
- Year (defaults to current year)
- Month (1-12, empty for whole year)
- Week start day (monday, sunday, etc.)
- Whether to show calendar week numbers
- Whether to show only workdays (Monday-Friday)
- Whether to add a comments column

### Command-Line Mode

Run mdcal with arguments and flags:

```bash
# Generate calendar for current year and month with default options
mdcal

# Generate calendar for a specific year
mdcal 2025

# Generate calendar for a specific year and month
mdcal 2025 12

# Generate calendar with custom options
mdcal 2025 12 --weekstart=sunday --week=false --workweek=true --comments=false --justify=center
```

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `--weekstart` | First day of the week (monday, sunday, etc.) | monday |
| `--week` | Show calendar week numbers | true |
| `--workweek` | Show only workdays (Monday-Friday) | false |
| `--comments` | Add a comments column | true |
| `--justify` | Cell justification: left, center, or right | left |
| `--version` | Print version information | - |
| `--help` | Show help information | - |

## Example Output

```markdown
| CW | Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday | Comments |
| :-: | :----: | :-----: | :-------: | :------: | :----: | :------: | :----: | :------: |
| 49 | 1      | 2       | 3         | 4        | 5      | 6        | 7      |          |
| 50 | 8      | 9       | 10        | 11       | 12     | 13       | 14     |          |
| 51 | 15     | 16      | 17        | 18       | 19     | 20       | 21     |          |
| 52 | 22     | 23      | 24        | 25       | 26     | 27       | 28     |          |
| 1  | 29     | 30      | 31        |          |        |          |        |          |
```

When rendered in Markdown, this produces a calendar table like:

| CW | Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday | Comments |
| :-: | :----: | :-----: | :-------: | :------: | :----: | :------: | :----: | :------: |
| 49 | 1      | 2       | 3         | 4        | 5      | 6        | 7      |          |
| 50 | 8      | 9       | 10        | 11       | 12     | 13       | 14     |          |
| 51 | 15     | 16      | 17        | 18       | 19     | 20       | 21     |          |
| 52 | 22     | 23      | 24        | 25       | 26     | 27       | 28     |          |
| 1  | 29     | 30      | 31        |          |        |          |        |          |

## Version

Current version: 0.1.1

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Inspiration

The original inspiration for this project came from [mdcal (Python)](https://github.com/pn11/mdcal) after a coworker spent a long time making changes to the output.
