# FedRAMP R5 Balance & 20x Phase One Implementations

This repository now includes implementations for:
1. **FedRAMP's R5 Balance Improvement Releases (BIRs)** as outlined in the [FedRAMP roadmap](https://github.com/FedRAMP/roadmap/issues?q=type:%22R5%20Balance%22)
2. **FedRAMP 20x Phase One Pilot** for automated Key Security Indicators (KSI) validation

## 🚀 What We've Built

### FedRAMP 20x Phase One Support (NEW!)
**Status**: ✅ Complete Implementation for 20x Pilot  
**Documentation**: [FedRAMP 20x Phase One](https://www.fedramp.gov/20x/)

**Features:**
- Complete Key Security Indicators (KSI) validation framework
- Machine-readable assessment reports in JSON format
- Automated validation for cloud-native security capabilities
- Continuous reporting proposal generator
- 3PAO attestation support
- Trust Center configuration templates

**Key Files:**
- `pkg/fedramp/ksi.go` - KSI validation framework
- `pkg/fedramp/continuous_reporting.go` - Continuous reporting capabilities
- `cli/cmd/ksi.go` - CLI commands for KSI operations
- `examples/fedramp_20x_demo.go` - Complete 20x submission demo

## 🚀 R5 Balance Implementations

### 1. R5.SCN - Significant Change Notification Support
**Status**: ✅ Initial Implementation Complete  
**GitHub Issue**: [#21](https://github.com/FedRAMP/roadmap/issues/21)

**Features:**
- Complete SCN data structures following RFC-0007 specifications
- Automatic classification of change types (Adaptive, Transformative, Impact Change)
- JSON export for machine-readable format requirements
- Validation engine for submission readiness
- CLI tools for SCN creation, validation, and export

**Key Files:**
- `pkg/fedramp/scn.go` - Core SCN data structures and logic
- `cli/cmd/scn.go` - CLI commands for SCN operations

**Usage Example:**
```bash
# Create a new SCN
gocomply_fedramp scn create CSO-12345 "security-patch" "Apply critical patches" "CVE remediation" \
  --affected-controls SI-2,RA-5 \
  --approver-name "Jane Smith" \
  --approver-title "CISO"

# Validate an SCN
gocomply_fedramp scn validate scn.json

# Export summary
gocomply_fedramp scn export summary scn.json summary.txt
```

### 2. R5.CRS - Continuous Reporting Standard Support  
**Status**: ✅ Initial Implementation Complete  
**GitHub Issue**: [#23](https://github.com/FedRAMP/roadmap/issues/23)

**Features:**
- Key Security Metrics (KSM) data structures
- Standard FedRAMP metrics (vulnerability scanning, patch compliance, MFA coverage, etc.)
- Continuous monitoring report generation
- Dashboard data aggregation
- Automated trend analysis framework

**Key Files:**
- `pkg/fedramp/crs.go` - Continuous Reporting Standard implementation

**Standard Metrics Included:**
- Vulnerability Scan Coverage
- Security Patch Compliance  
- Failed Login Attempts
- Backup Success Rate
- Data Encryption Coverage
- Multi-Factor Authentication Coverage

## 🛠️ Technical Implementation

### Architecture Principles
These implementations follow FedRAMP's modernization goals:

1. **Machine-Readable First**: All outputs are JSON-based for automation
2. **Security-Based Assessment**: Focus on actual security metrics vs. compliance documentation
3. **Continuous Monitoring**: Real-time data streams replace periodic reports
4. **Automation-Friendly**: Designed for integration with CI/CD and monitoring tools

### Integration Points
- **OSCAL Compatibility**: Works with existing OSCAL SSP processing
- **FedRAMP Baselines**: Integrates with Low/Moderate/High baseline catalogs
- **CLI Integration**: Extends existing gocomply_fedramp command structure

## 🧪 Demo & Testing

Run the demonstration:
```bash
go run examples/r5_balance_demo.go
```

This demo shows:
- SCN creation and validation workflow
- Continuous reporting with standard metrics
- JSON export capabilities
- Integration between SCN and CRS systems

## 📈 R5 Balance Timeline Support

Based on the [roadmap issue #16](https://github.com/FedRAMP/roadmap/issues/16):

- ✅ **Preparation Phase** (June 2-27): Core implementations ready
- 🔄 **Closed Beta** (July 7 - August 29): Ready for beta testing
- 🔄 **Open Beta** (September 1 - October 31): Community feedback integration
- 🔄 **Wide Release** (November 3): Production readiness

## 🤝 Contributing to R5 Balance

### Immediate Contribution Opportunities:

1. **Beta Testing**: Test the SCN and CRS implementations with real-world scenarios
2. **Integration Examples**: Build connectors to popular monitoring tools
3. **Metric Expansion**: Add additional Key Security Metrics
4. **Validation Enhancement**: Improve SCN classification logic
5. **Dashboard Development**: Create visualization tools for CRS data

### 3. R5.MAS - Minimum Assessment Standard  
**Status**: ✅ Complete Implementation  
**GitHub Issue**: [#19](https://github.com/FedRAMP/roadmap/issues/19)

**Features:**
- Complete assessment framework following MAS requirements
- Support for all assessment types (initial, annual, significant change, incident)
- Assessment method tracking (examine, interview, test)
- Evidence collection and management
- Automated vs manual assessment tracking
- 3PAO attestation support
- Findings management with severity levels
- Validation of MAS compliance

**Key Files:**
- `pkg/fedramp/mas.go` - MAS data structures and validation
- `cli/cmd/mas.go` - CLI commands for MAS operations
- `examples/mas_demo.go` - Complete MAS assessment demo

### 4. R5.SSAD - Storing and Sharing Authorization Data  
**Status**: ✅ Complete Implementation  
**GitHub Issue**: [#28](https://github.com/FedRAMP/roadmap/issues/28)

**Features:**
- Standardized authorization data packages
- Document management (SSP, SAP, SAR, POA&M, ConMon, KSI)
- Access control with role-based permissions
- Distribution tracking and audit logs
- Package integrity with cryptographic hashing
- Repository management with search capabilities
- Support for FedRAMP 20x KSI reports
- Automated metadata extraction

**Key Files:**
- `pkg/fedramp/ssad.go` - SSAD package and repository management
- `cli/cmd/ssad.go` - CLI commands for SSAD operations
- `examples/ssad_demo.go` - Complete SSAD workflow demo

## 🔗 Related Resources

- [RFC-0007 Significant Change Notification Standard](https://www.fedramp.gov/rfcs/0007/)
- [RFC-0008 Continuous Reporting Standard](https://www.fedramp.gov/rfcs/0008/)
- [FedRAMP 20x Phase One Key Security Indicators](https://www.fedramp.gov/rfcs/0006/)
- [FedRAMP Roadmap Repository](https://github.com/FedRAMP/roadmap)

## 💡 Why This Matters

These R5 Balance implementations represent FedRAMP's shift from:
- **Manual → Automated** assessment processes
- **Compliance → Security** focused evaluation  
- **Point-in-time → Continuous** monitoring
- **Narrative → Machine-readable** documentation

By contributing to these implementations, you're helping modernize federal cloud security assessment for the entire government.

---

**Ready to contribute?** Join the [Rev5 Community Working Group](https://github.com/FedRAMP/roadmap/issues/16) and help shape the future of FedRAMP automation!

## 📚 Complete Usage Guide

### FedRAMP 20x Phase One Operations

```bash
# 1. Validate Key Security Indicators
gocomply_fedramp ksi validate CSO-EXAMPLE-001 --output ksi-report.json

# 2. Generate continuous reporting proposal
gocomply_fedramp ksi proposal --service-id CSO-EXAMPLE-001 --output proposal.json

# 3. Generate continuous monitoring report
gocomply_fedramp ksi report --service-id CSO-EXAMPLE-001 --output continuous-report.json

# 4. Run complete 20x demo (generates all required files)
go run examples/fedramp_20x_demo.go
```

### Submitting to FedRAMP 20x Pilot

1. **Run the complete demo** to generate all required files
2. **Review generated files**:
   - `ksi-report.json` - KSI validation results with evidence
   - `continuous-reporting-proposal.json` - Automated reporting plan
   - `fedramp-20x-submission.json` - Complete submission package
   - `trust-center.json` - Trust center configuration

3. **Send submission email** to 20x@fedramp.gov with:
   - Summary of your cloud service
   - Points of contact for CSP and 3PAO
   - Instructions to access the submission package (do NOT attach files)

4. **Share in community** (optional but encouraged):
   - Post draft submission to [FedRAMP Community Discussions](https://github.com/FedRAMP/community/discussions)
   - Get feedback from other participants
   - Learn from FedRAMP's responses to other submissions

### Building Custom KSI Validators

```go
// Example: Add custom KSI validation
evidence := []fedramp.KSIEvidence{
    {
        Type:        "api_check",
        Description: "Verified encryption via AWS KMS API",
        Reference:   "kms-check-20250110",
        Source:      "AWS API",
        Timestamp:   time.Now(),
    },
}

validation := fedramp.ValidateKSI("KSI-SC", evidence, true)
report.AddValidation(validation)
```

## 🔧 Development Setup

```bash
# Clone the repository
git clone https://github.com/gocomply/fedramp.git
cd fedramp

# Build all components
go build ./...

# Run tests
go test ./...

# Build CLI
go build -o gocomply_fedramp cli/gocomply_fedramp/main.go
```

## 🎯 Implementation Status

| Feature | R5.SCN | R5.CRS | R5.MAS | R5.SSAD | 20x KSI | 20x Continuous |
|---------|--------|--------|--------|---------|---------|----------------|
| Core Implementation | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| CLI Commands | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| JSON Export | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Validation | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| 3PAO Support | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Demo Applications | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| API Integration | 🔄 | 🔄 | 🔄 | 🔄 | ✅ | ✅ |
| OSCAL Integration | 🔄 | 🔄 | 🔄 | 🔄 | 🔄 | 🔄 |

Legend: ✅ Complete | 🔄 In Progress | ⏳ Planned

---

## FedRAMP Machine Readable (FRMR) Support

This implementation now includes support for the official FedRAMP Machine Readable (FRMR) document formats from the [FedRAMP/docs](https://github.com/FedRAMP/docs) repository. FRMR provides structured JSON schemas for FedRAMP requirements, making them easier to parse and validate programmatically.

### FRMR Features

- **Parser**: Full support for FRMR JSON schema including FRD (Definitions), FRR (Rules), FRA (Assistance), and KSI (Key Security Indicators)
- **Validation**: Validate KSI requirements against evidence with detailed reporting
- **CLI Integration**: Fetch, validate, and export FRMR documents
- **Official Alignment**: Direct compatibility with FedRAMP's official machine-readable formats

### FRMR CLI Commands

```bash
# Fetch official FRMR documents from FedRAMP repository
gocomply_fedramp frmr fetch ksi
gocomply_fedramp frmr fetch mas
gocomply_fedramp frmr fetch scn

# Display information about a FRMR document
gocomply_fedramp frmr info FRMR.KSI.key-security-indicators.json

# Validate KSI requirements against evidence
gocomply_fedramp frmr validate FRMR.KSI.key-security-indicators.json evidence.json

# Export FRMR data in different formats
gocomply_fedramp frmr export FRMR.KSI.key-security-indicators.json --format markdown --output ksi.md
```

### FRMR Tools and Capabilities

The FRMR implementation includes comprehensive tools aligned with the official FedRAMP/docs repository:

#### Document Management
- **Fetch**: Download official FRMR documents from FedRAMP/docs
- **Combine**: Merge multiple FRMR documents (useful for combined baselines)
- **Filter**: Extract specific requirements by impact level, type, or ID
- **Validate**: Check documents against FedRAMP schema requirements

#### Assessment Support
- **Evidence Templates**: Generate starter evidence files for KSI validation
- **Validation Reports**: Produce detailed compliance assessments
- **Markdown Export**: Convert FRMR to human-readable documentation

#### Advanced CLI Commands

```bash
# Combine multiple documents
gocomply_fedramp frmr combine FRMR.KSI.*.json --output combined-20x.json

# Filter for specific impact levels
gocomply_fedramp frmr filter FRMR.KSI.*.json --impact Low --output low-ksi.json

# Generate evidence template
gocomply_fedramp frmr evidence-template FRMR.KSI.*.json --output my-evidence.json

# Validate document schema
gocomply_fedramp frmr schema-validate my-frmr-doc.json

# Advanced filtering
gocomply_fedramp frmr filter doc.json --ksi KSI-IAM,KSI-MLA --type KSI,FRR
```

### FRMR Demos

Run the comprehensive demos to see all features:

```bash
# Basic FRMR operations
go run examples/frmr_demo.go

# Advanced tools demonstration
go run examples/frmr_tools_demo.go
```

These demonstrate:
- Fetching and parsing official documents
- Schema validation
- Document filtering and combining
- Evidence template generation
- Markdown export with templates
- Compliance scoring and reporting

---

**Questions?** Open an issue or join the discussion at [FedRAMP Community](https://github.com/FedRAMP/community/discussions) 