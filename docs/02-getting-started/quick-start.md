# Quick Start Guide

## 🚀 Get Started in 5 Minutes

This guide will help you quickly set up and start using the FedRAMP R5 Balance & 20x implementation toolkit.

## Prerequisites

- Go 1.19 or higher
- Git
- Basic understanding of FedRAMP requirements
- (Optional) Docker for containerized deployment

## Installation

### Option 1: Install from Source

```bash
# Clone the repository
git clone https://github.com/gocomply/fedramp.git
cd fedramp

# Build the CLI tool
go build -o gocomply_fedramp cli/gocomply_fedramp/main.go

# Add to PATH (optional)
sudo mv gocomply_fedramp /usr/local/bin/
```

### Option 2: Go Install

```bash
go install github.com/gocomply/fedramp/cli/gocomply_fedramp@latest
```

### Option 3: Docker

```bash
docker pull gocomply/fedramp:latest
alias gocomply_fedramp='docker run -it --rm -v $(pwd):/workspace gocomply/fedramp:latest'
```

## First Steps

### 1. Verify Installation

```bash
gocomply_fedramp --help
```

Expected output:
```
NAME:
   gocomply_fedramp - OSCAL-FedRAMP Workbench

USAGE:
   gocomply_fedramp [global options] command [command options] [arguments...]

COMMANDS:
   convert      Convert OSCAL SSP to FedRAMP Document
   opencontrol  Convert OpenControl masonry repo into FedRAMP formatted OSCAL
   scn          Significant Change Notification operations for R5.SCN BIR
   ksi          FedRAMP 20x Key Security Indicators operations
   mas          Minimum Assessment Standard operations for R5.MAS
   ssad         Storing and Sharing Authorization Data operations for R5.SSAD
   frmr         Work with FedRAMP Machine Readable (FRMR) documents
   help, h      Shows a list of commands or help for one command
```

### 2. Download Official FedRAMP Documents

```bash
# Fetch the latest Key Security Indicators
gocomply_fedramp frmr fetch ksi

# Fetch other FRMR documents
gocomply_fedramp frmr fetch mas
gocomply_fedramp frmr fetch scn
```

### 3. Generate Your First KSI Evidence Template

```bash
# Create an evidence template for KSI validation
gocomply_fedramp frmr evidence-template FRMR.KSI.key-security-indicators.json \
  --output my-evidence.json

# Edit the template to reflect your implementation status
# Set each requirement to true/false based on your compliance
```

### 4. Validate Your KSI Compliance

```bash
# Run KSI validation
gocomply_fedramp frmr validate FRMR.KSI.key-security-indicators.json my-evidence.json
```

### 5. Create a Significant Change Notification

```bash
# Create an SCN for a security patch
gocomply_fedramp scn create CSO-12345 "security-patch" \
  "Apply critical security patches" \
  "Remediate CVE-2024-1234" \
  --affected-controls SI-2,RA-5 \
  --approver-name "John Doe" \
  --approver-title "CISO"
```

## Common Workflows

### For Cloud Service Providers (CSPs)

#### Initial 20x Submission
```bash
# 1. Generate KSI validation report
gocomply_fedramp ksi validate MY-CSO-001 --output ksi-report.json

# 2. Create continuous reporting proposal
gocomply_fedramp ksi proposal --service-id MY-CSO-001 --output proposal.json

# 3. Generate submission package
gocomply_fedramp ksi package MY-CSO-001 --output submission.zip
```

#### Continuous Monitoring
```bash
# Generate monthly ConMon report
gocomply_fedramp crs report \
  --scan-coverage 98.5 \
  --patch-compliance 99.2 \
  --failed-logins 12 \
  --backup-success 100 \
  --encryption-coverage 100 \
  --mfa-coverage 95.8 \
  --output monthly-conmon.json
```

### For 3PAOs

#### KSI Assessment
```bash
# Validate CSP's KSI evidence
gocomply_fedramp frmr validate ksi-doc.json csp-evidence.json \
  > assessment-results.txt

# Generate assessment report
gocomply_fedramp mas create \
  --type initial \
  --cso-id CSO-12345 \
  --assessor "3PAO Company" \
  --output assessment.json
```

### For Federal Agencies

#### Review Authorization Package
```bash
# Fetch and review CSP documents
gocomply_fedramp ssad fetch CSO-12345 --output csp-package/

# Validate package completeness
gocomply_fedramp ssad validate csp-package/

# Generate review report
gocomply_fedramp ssad review csp-package/ --output review-report.json
```

## Quick Examples

### Example 1: Check Your Compliance Score
```bash
# Create a simple evidence file
cat > evidence.json << EOF
{
  "KSI-IAM-01": true,
  "KSI-IAM-02": true,
  "KSI-IAM-03": false,
  "KSI-MLA-01": true,
  "KSI-MLA-02": false
}
EOF

# Check compliance
gocomply_fedramp frmr validate FRMR.KSI.*.json evidence.json
```

### Example 2: Filter Requirements by Impact Level
```bash
# Get only Low impact requirements
gocomply_fedramp frmr filter FRMR.KSI.*.json \
  --impact Low \
  --output low-requirements.json

# Export to markdown for review
gocomply_fedramp frmr export low-requirements.json \
  --format markdown \
  --output low-requirements.md
```

### Example 3: Track Changes Over Time
```bash
# Create initial baseline
gocomply_fedramp scn create CSO-001 "initial" "Initial deployment" "v1.0"

# Record a change
gocomply_fedramp scn create CSO-001 "update" "Security update" "v1.1" \
  --change-type adaptive

# View change history
gocomply_fedramp scn list CSO-001
```

## Next Steps

- 📖 Read the [detailed documentation](../README.md)
- 💡 Explore [use cases](../10-use-cases/README.md) for your role
- 🔧 Learn about [API integration](../08-api-docs/README.md)
- 🤝 Join the [FedRAMP Community](https://github.com/FedRAMP/community)

## Troubleshooting

### Common Issues

**Command not found**
```bash
# Ensure gocomply_fedramp is in your PATH
export PATH=$PATH:/path/to/gocomply_fedramp
```

**Permission denied**
```bash
# Make the binary executable
chmod +x gocomply_fedramp
```

**Network errors fetching documents**
```bash
# Check your internet connection and proxy settings
export HTTPS_PROXY=your-proxy:port  # if behind a proxy
```

### Getting Help

```bash
# Get help for any command
gocomply_fedramp help
gocomply_fedramp ksi help
gocomply_fedramp frmr help validate

# Check version
gocomply_fedramp --version
```

## Quick Reference Card

| Task | Command |
|------|---------|
| Fetch KSI document | `gocomply_fedramp frmr fetch ksi` |
| Validate compliance | `gocomply_fedramp frmr validate ksi.json evidence.json` |
| Create SCN | `gocomply_fedramp scn create [options]` |
| Generate ConMon report | `gocomply_fedramp crs report [metrics]` |
| Export to markdown | `gocomply_fedramp frmr export doc.json --format markdown` |
| Filter by impact | `gocomply_fedramp frmr filter doc.json --impact Low` |

---

🎉 **Congratulations!** You're now ready to use the FedRAMP R5 Balance & 20x toolkit. For more detailed information, explore the full documentation. 