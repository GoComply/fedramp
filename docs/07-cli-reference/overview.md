# CLI Command Reference Overview

## Introduction

The `gocomply_fedramp` CLI provides comprehensive tools for FedRAMP compliance automation, including support for R5 Balance initiatives and the 20x Phase One pilot program.

## Installation

```bash
# From source
go install github.com/gocomply/fedramp/cli/gocomply_fedramp@latest

# Verify installation
gocomply_fedramp --version
```

## Global Options

All commands support these global options:

```bash
--help, -h      Show help
--version, -v   Print version information
--verbose       Enable verbose output
--debug         Enable debug logging
--config FILE   Specify configuration file
--output FILE   Output file (default: stdout)
--format FORMAT Output format (json|yaml|text)
```

## Command Structure

```
gocomply_fedramp [global options] <command> [command options] [arguments...]
```

## Available Commands

### Core Commands

| Command | Description | R5/20x Feature |
|---------|-------------|----------------|
| `convert` | Convert OSCAL SSP to FedRAMP Document | Legacy |
| `opencontrol` | Convert OpenControl to OSCAL | Legacy |
| `scn` | Significant Change Notification | R5.SCN |
| `ksi` | Key Security Indicators | 20x Phase One |
| `mas` | Minimum Assessment Standard | R5.MAS |
| `ssad` | Storing and Sharing Authorization Data | R5.SSAD |
| `frmr` | FedRAMP Machine Readable documents | FRMR Tools |

### Command Categories

#### 1. Document Conversion (Legacy)
- `convert` - Transform OSCAL to FedRAMP templates
- `opencontrol` - Migrate from OpenControl format

#### 2. R5 Balance Commands
- `scn` - Change management and notifications
- `mas` - Assessment standards and evidence
- `ssad` - Document storage and sharing
- `crs` - Continuous reporting (via `ksi` command)

#### 3. FedRAMP 20x Commands
- `ksi` - Key Security Indicator validation
- `ksi proposal` - Continuous reporting proposals
- `ksi report` - Generate monitoring reports

#### 4. FRMR Tools
- `frmr fetch` - Download official documents
- `frmr validate` - Validate compliance
- `frmr export` - Transform to various formats
- `frmr filter` - Filter by criteria

## Quick Command Examples

### Basic Operations

```bash
# Get help for any command
gocomply_fedramp help
gocomply_fedramp scn help
gocomply_fedramp frmr help validate

# Show version
gocomply_fedramp --version
```

### Common Workflows

```bash
# 1. Fetch latest KSI requirements
gocomply_fedramp frmr fetch ksi

# 2. Create evidence template
gocomply_fedramp frmr evidence-template FRMR.KSI.*.json --output evidence.json

# 3. Validate compliance
gocomply_fedramp frmr validate FRMR.KSI.*.json evidence.json

# 4. Generate report
gocomply_fedramp ksi report --service-id CSO-001 --output report.json
```

## Command Chaining

Commands can be chained using standard Unix pipes:

```bash
# Fetch, filter, and export in one pipeline
gocomply_fedramp frmr fetch ksi | \
gocomply_fedramp frmr filter --impact Low | \
gocomply_fedramp frmr export --format markdown > low-requirements.md
```

## Configuration Files

Create a `.fedramp.yaml` configuration file:

```yaml
# .fedramp.yaml
defaults:
  output-format: json
  verbose: true

ksi:
  service-id: CSO-12345
  evidence-path: ./evidence/

scn:
  approver-name: "Security Team"
  approver-email: "security@company.com"

frmr:
  cache-dir: ~/.fedramp/cache
  github-token: ${GITHUB_TOKEN}
```

## Environment Variables

```bash
# Set default CSO ID
export FEDRAMP_CSO_ID="CSO-12345"

# Set GitHub token for FRMR operations
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"

# Set output directory
export FEDRAMP_OUTPUT_DIR="/var/fedramp/reports"
```

## Output Formats

Most commands support multiple output formats:

### JSON (Default)
```bash
gocomply_fedramp ksi validate CSO-001 --format json
```

### YAML
```bash
gocomply_fedramp scn list --format yaml
```

### Human-Readable Text
```bash
gocomply_fedramp frmr info document.json --format text
```

### Markdown
```bash
gocomply_fedramp frmr export document.json --format markdown
```

## Error Handling

The CLI uses standard exit codes:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Command line usage error |
| 3 | Configuration error |
| 4 | Network/connectivity error |
| 5 | Validation failure |
| 10 | File not found |
| 11 | Permission denied |

Example error handling:

```bash
#!/bin/bash
gocomply_fedramp ksi validate CSO-001
case $? in
  0) echo "Validation passed!" ;;
  5) echo "Validation failed - review requirements" ;;
  *) echo "Unexpected error occurred" ;;
esac
```

## Logging and Debug

Enable detailed logging:

```bash
# Verbose output
gocomply_fedramp --verbose ksi validate CSO-001

# Debug logging
gocomply_fedramp --debug frmr fetch ksi

# Log to file
gocomply_fedramp --verbose ksi validate CSO-001 2> fedramp.log
```

## Performance Tips

1. **Use caching for FRMR operations**
   ```bash
   export FEDRAMP_CACHE_DIR=~/.fedramp/cache
   ```

2. **Batch operations**
   ```bash
   # Process multiple files at once
   gocomply_fedramp frmr validate FRMR.*.json evidence.json
   ```

3. **Parallel processing**
   ```bash
   # Run validations in parallel
   parallel gocomply_fedramp ksi validate {} ::: CSO-001 CSO-002 CSO-003
   ```

## Shell Completion

Enable tab completion:

### Bash
```bash
gocomply_fedramp completion bash > /etc/bash_completion.d/gocomply_fedramp
```

### Zsh
```bash
gocomply_fedramp completion zsh > "${fpath[1]}/_gocomply_fedramp"
```

### Fish
```bash
gocomply_fedramp completion fish > ~/.config/fish/completions/gocomply_fedramp.fish
```

## Aliases and Functions

Add to your shell configuration:

```bash
# ~/.bashrc or ~/.zshrc

# Quick aliases
alias fr='gocomply_fedramp'
alias frv='gocomply_fedramp --verbose'
alias frksi='gocomply_fedramp ksi'
alias frscn='gocomply_fedramp scn'

# Useful functions
frvalidate() {
  gocomply_fedramp frmr validate FRMR.KSI.*.json "$1" --format text
}

frmonthly() {
  gocomply_fedramp ksi report \
    --service-id ${FEDRAMP_CSO_ID} \
    --output "monthly-$(date +%Y%m).json"
}
```

## Integration with Other Tools

### jq for JSON processing
```bash
gocomply_fedramp ksi validate CSO-001 | jq '.summary'
```

### Generate reports with pandoc
```bash
gocomply_fedramp frmr export doc.json --format markdown | \
  pandoc -o report.pdf
```

### Send notifications
```bash
gocomply_fedramp ksi validate CSO-001 || \
  mail -s "KSI Validation Failed" security@company.com
```

## Troubleshooting

### Command not found
```bash
# Check installation
which gocomply_fedramp

# Add to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin
```

### Network issues
```bash
# Use proxy if needed
export HTTPS_PROXY=http://proxy:8080

# Increase timeout
export FEDRAMP_TIMEOUT=300
```

### Debug mode
```bash
# Enable all debug output
export FEDRAMP_DEBUG=1
gocomply_fedramp --debug ksi validate CSO-001
```

## Getting Help

```bash
# General help
gocomply_fedramp help

# Command-specific help
gocomply_fedramp ksi help
gocomply_fedramp frmr help validate

# Online documentation
open https://github.com/gocomply/fedramp/docs
```

---

For detailed information about specific commands, see the individual command reference pages. 