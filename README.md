# mdcal

A CLI tool for creating markdown calendars.

![GitHub License](https://img.shields.io/github/license/andre-a-alves/mdcal?style=for-the-badge)
![GitHub Release](https://img.shields.io/github/v/release/andre-a-alves/mdcal?style=for-the-badge)

## Description

mdcal is a command-line utility that generates customizable markdown calendars for a specific month, a range of months, or an entire year. It can be run either with command-line flags or in interactive mode, allowing you to create calendars with various options such as custom week start days, calendar week numbers, and more.

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

The basic syntax for mdcal is:

```bash
mdcal [year] [month] [endMonth|endYear endMonth]
```

Where:
- `year`: The year for the calendar (defaults to current year)
- `month`: The month to start the calendar (1-12, omit for whole year)
- `endMonth`: When specifying a range within the same year
- `endYear endMonth`: When specifying a range across different years

### Interactive Mode

Run mdcal without any arguments to enter interactive mode:

```bash
mdcal
```

In interactive mode, you'll be prompted to enter:
- Year (defaults to current year)
- Month (1-12, empty for whole year)
- First day of the week (monday/mon, sunday/sun, etc.)
- Whether to show calendar week numbers
- Whether to leave weekends off the calendar
- Whether to leave the comments column off
- Cell justification (left, center, or right)

Note: The month range feature is currently only available through command-line arguments, not in interactive mode.

### Command-Line Mode

Run mdcal with arguments and flags:

```bash
# Generate calendar for current year and month with default options
mdcal

# Generate calendar for a specific year
mdcal 2025

# Generate calendar for a specific year and month
mdcal 2025 12

# Generate calendar for a range of months within the same year (March through May 2025)
mdcal 2025 3 5

# Generate calendar for a range of months across different years (December 2025 through January 2026)
mdcal 2025 12 2026 1

# Generate calendar with custom options
mdcal 2025 12 --start=sunday --no-week-no --workweek=true --no-comment=true --justify=center

# Using short flags
mdcal 2025 12 -s sunday -w -W -c -j center
```

Note: When specifying a range of months, the end date must be after the start date. If an invalid range is provided (e.g., `mdcal 2026 1 2025 12`), an error message will be displayed.

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `-s, --start` | First day of the week (monday/mon) | monday |
| `-w, --no-week-no` | Leave week numbers off the calendar | false |
| `-W, --workweek` | Leave weekends off the calendar | false |
| `-c, --no-comment` | Leave the comments column off | false |
| `-j, --justify` | Cell justification: left, center, or right | left |
| `-v, --version` | Print version information | - |
| `-h, --help` | Show help information | - |

## Example Output

```markdown
| CW   | Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday | Comments |
| :--: | :----: | :-----: | :-------: | :------: | :----: | :------: | :----: | :------: |
| _49_ | 1      | 2       | 3         | 4        | 5      | 6        | 7      |          |
| _50_ | 8      | 9       | 10        | 11       | 12     | 13       | 14     |          |
| _51_ | 15     | 16      | 17        | 18       | 19     | 20       | 21     |          |
| _52_ | 22     | 23      | 24        | 25       | 26     | 27       | 28     |          |
| _1_  | 29     | 30      | 31        |          |        |          |        |          |
```

When rendered in Markdown, this produces a calendar table like:

| CW   | Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday | Comments |
|:----:| :----: | :-----: | :-------: | :------: | :----: | :------: | :----: | :------: |
| _49_ | 1      | 2       | 3         | 4        | 5      | 6        | 7      |          |
| _50_ | 8      | 9       | 10        | 11       | 12     | 13       | 14     |          |
| _51_ | 15     | 16      | 17        | 18       | 19     | 20       | 21     |          |
| _52_ | 22     | 23      | 24        | 25       | 26     | 27       | 28     |          |
| _1_  | 29     | 30      | 31        |          |        |          |        |          |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Inspiration

The original inspiration for this project came from [mdcal (Python)](https://github.com/pn11/mdcal) after a coworker spent a long time making changes to the output.
