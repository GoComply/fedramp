# Key Security Indicators (KSI) Guide

## Overview

Key Security Indicators (KSIs) are the foundation of FedRAMP 20x Phase One, representing a paradigm shift from control-based to outcome-based security assessments. This guide provides comprehensive information on implementing, validating, and maintaining KSI compliance.

## What are KSIs?

Key Security Indicators are measurable security outcomes that demonstrate effective implementation of critical security capabilities. Unlike traditional control assessments, KSIs focus on:

- **Outcomes over Process**: What security is achieved, not how it's documented
- **Automation**: Machine-readable validation enabling continuous assessment
- **Binary Results**: Clear true/false determinations with evidence
- **Risk Focus**: Addressing the most critical security concerns

## The 11 Official KSIs (FedRAMP 25.05C)

### 1. KSI-CED: Cybersecurity Education
Ensures personnel have appropriate security awareness and training.

**Requirements:**
- `KSI-CED-01`: Security awareness training for all users
- `KSI-CED-02`: Role-based security training for privileged users

### 2. KSI-CMT: Change Management
Manages system changes to maintain security posture.

**Requirements:**
- `KSI-CMT-01`: Change control board established
- `KSI-CMT-02`: Change approval process documented
- `KSI-CMT-03`: Emergency change procedures defined
- `KSI-CMT-04`: Change testing requirements
- `KSI-CMT-05`: Rollback procedures available

### 3. KSI-CNA: Cloud Native Architecture
Leverages cloud-native security capabilities.

**Requirements:**
- `KSI-CNA-01`: Infrastructure as Code (IaC) implementation
- `KSI-CNA-02`: Container security controls
- `KSI-CNA-03`: Microservices architecture security
- `KSI-CNA-04`: API security gateway
- `KSI-CNA-05`: Service mesh security
- `KSI-CNA-06`: Serverless security controls
- `KSI-CNA-07`: Cloud-native monitoring

### 4. KSI-IAM: Identity and Access Management
Controls access to systems and data.

**Requirements:**
- `KSI-IAM-01`: Multi-factor authentication (MFA) for all users
- `KSI-IAM-02`: Privileged access management (PAM)
- `KSI-IAM-03`: Regular access reviews
- `KSI-IAM-04`: Automated provisioning/deprovisioning
- `KSI-IAM-05`: Single sign-on (SSO) implementation
- `KSI-IAM-06`: Zero trust principles applied

### 5. KSI-INR: Incident Reporting
Ensures timely detection and reporting of security incidents.

**Requirements:**
- `KSI-INR-01`: 24/7 incident response capability
- `KSI-INR-02`: Automated incident detection
- `KSI-INR-03`: US-CERT reporting integration

### 6. KSI-MLA: Monitoring, Logging, and Auditing
Provides comprehensive visibility into system activities.

**Requirements:**
- `KSI-MLA-01`: Centralized log management
- `KSI-MLA-02`: Real-time security monitoring
- `KSI-MLA-03`: Log retention per requirements
- `KSI-MLA-04`: Automated anomaly detection
- `KSI-MLA-05`: Audit trail protection
- `KSI-MLA-06`: Performance monitoring

### 7. KSI-PIY: Policy and Inventory
Maintains accurate system documentation and policies.

**Requirements:**
- `KSI-PIY-01`: Complete asset inventory
- `KSI-PIY-02`: Automated inventory updates
- `KSI-PIY-03`: Security policy documentation
- `KSI-PIY-04`: Policy enforcement automation
- `KSI-PIY-05`: Configuration baselines
- `KSI-PIY-06`: Software bill of materials (SBOM)
- `KSI-PIY-07`: Data classification implemented

### 8. KSI-RPL: Recovery Planning
Ensures business continuity and disaster recovery.

**Requirements:**
- `KSI-RPL-01`: Tested backup procedures
- `KSI-RPL-02`: Documented recovery time objectives (RTO)
- `KSI-RPL-03`: Documented recovery point objectives (RPO)
- `KSI-RPL-04`: Annual DR testing

### 9. KSI-SVC: Service Configuration
Maintains secure system configurations.

**Requirements:**
- `KSI-SVC-01`: Hardened configurations
- `KSI-SVC-02`: Automated configuration management
- `KSI-SVC-03`: Vulnerability scanning
- `KSI-SVC-04`: Patch management automation
- `KSI-SVC-05`: Secure baseline enforcement
- `KSI-SVC-06`: Configuration drift detection
- `KSI-SVC-07`: Compliance scanning automation

### 10. KSI-TPR: Third-Party Information Resources
Manages supply chain and third-party risks.

**Requirements:**
- `KSI-TPR-01`: Third-party risk assessments
- `KSI-TPR-02`: Vendor security requirements
- `KSI-TPR-03`: Supply chain visibility
- `KSI-TPR-04`: Continuous vendor monitoring

### 11. KSI-VUL: Vulnerability Management
(Note: Combined with KSI-SVC in 25.05C release)

## Implementation Guide

### Step 1: Assessment
```bash
# Generate an evidence template
gocomply_fedramp frmr evidence-template FRMR.KSI.key-security-indicators.json \
  --output ksi-evidence.json

# Review each requirement
cat ksi-evidence.json | jq '.requirements'
```

### Step 2: Evidence Collection

For each KSI requirement, collect supporting evidence:

```json
{
  "KSI-IAM-01": {
    "status": true,
    "evidence": {
      "description": "MFA enabled for all users via Azure AD",
      "artifacts": [
        "screenshots/azure-mfa-policy.png",
        "reports/mfa-coverage-report.pdf"
      ],
      "lastValidated": "2024-01-15T10:00:00Z"
    }
  }
}
```

### Step 3: Validation

```bash
# Validate your evidence
gocomply_fedramp frmr validate FRMR.KSI.key-security-indicators.json \
  ksi-evidence.json \
  --output validation-report.json

# Check your score
cat validation-report.json | jq '.summary'
```

### Step 4: Continuous Monitoring

```bash
# Set up automated validation
cat > ksi-monitor.sh << 'EOF'
#!/bin/bash
# Run daily KSI validation
gocomply_fedramp ksi validate CSO-001 \
  --evidence /path/to/evidence.json \
  --output /var/log/ksi/$(date +%Y%m%d).json

# Alert on failures
if [ $? -ne 0 ]; then
  mail -s "KSI Validation Failed" security@company.com < /var/log/ksi/$(date +%Y%m%d).json
fi
EOF

# Add to crontab
crontab -e
# 0 2 * * * /path/to/ksi-monitor.sh
```

## Best Practices

### 1. Evidence Management

- **Automate Collection**: Use APIs to gather evidence automatically
- **Version Control**: Track changes to evidence over time
- **Regular Updates**: Refresh evidence at least monthly
- **Clear Documentation**: Explain how each control is met

### 2. Implementation Priorities

Start with high-impact, easy-to-implement KSIs:

1. **KSI-IAM-01** (MFA): Critical for security, straightforward to implement
2. **KSI-MLA-01** (Centralized Logging): Foundation for other monitoring
3. **KSI-SVC-03** (Vulnerability Scanning): Essential for risk management
4. **KSI-INR-01** (Incident Response): Required for all systems

### 3. Common Pitfalls

- **Partial Implementation**: Each sub-requirement must be fully met
- **Stale Evidence**: Keep evidence current (< 30 days old)
- **Missing Documentation**: Document HOW requirements are met
- **Manual Processes**: Automate validation wherever possible

## Integration Examples

### Example 1: CI/CD Pipeline Integration

```yaml
# .github/workflows/ksi-validation.yml
name: KSI Validation
on:
  push:
    branches: [main]
  schedule:
    - cron: '0 0 * * *'  # Daily

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.19
          
      - name: Install FedRAMP Tools
        run: go install github.com/gocomply/fedramp/cli/gocomply_fedramp@latest
        
      - name: Validate KSIs
        run: |
          gocomply_fedramp ksi validate ${{ secrets.CSO_ID }} \
            --evidence evidence/ksi.json \
            --output results.json
            
      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: ksi-validation-results
          path: results.json
```

### Example 2: Monitoring Dashboard

```go
// ksi_dashboard.go
package main

import (
    "github.com/gocomply/fedramp/pkg/fedramp"
    "net/http"
)

func main() {
    http.HandleFunc("/api/ksi/status", func(w http.ResponseWriter, r *http.Request) {
        validation, err := fedramp.ValidateKSI("CSO-001", "evidence.json")
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(validation)
    })
    
    http.ListenAndServe(":8080", nil)
}
```

## Compliance Mapping

| KSI | Primary NIST Controls | FedRAMP Baseline |
|-----|----------------------|------------------|
| KSI-CED | AT-2, AT-3, AT-4 | All |
| KSI-CMT | CM-3, CM-4, CM-5 | All |
| KSI-CNA | SC-7, SC-8, SC-13 | Moderate+ |
| KSI-IAM | IA-2, IA-5, AC-2 | All |
| KSI-INR | IR-4, IR-5, IR-6 | All |
| KSI-MLA | AU-2, AU-3, AU-4 | All |
| KSI-PIY | CM-8, PL-2, RA-2 | All |
| KSI-RPL | CP-9, CP-10, CP-2 | All |
| KSI-SVC | CM-6, CM-7, RA-5 | All |
| KSI-TPR | SA-4, SA-9, SA-12 | Moderate+ |

## Frequently Asked Questions

### Q: What happens if a KSI fails validation?
A: Failed KSIs must be remediated within 30 days. The system remains authorized but enters enhanced monitoring status.

### Q: Can partial credit be given for KSIs?
A: No. Each sub-requirement is binary (true/false). All sub-requirements must be met for the KSI to pass.

### Q: How often should KSIs be validated?
A: Minimum monthly, but continuous automated validation is recommended.

### Q: What evidence is required?
A: Evidence must be:
- Current (< 30 days)
- Specific to the requirement
- Independently verifiable
- Machine-readable where possible

## Resources

- [Official FedRAMP KSI Documentation](https://github.com/FedRAMP/docs)
- [NIST SP 800-53 Control Catalog](https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final)
- [FedRAMP 20x Pilot Information](https://www.fedramp.gov/20x)
- [Community Discussion Forum](https://github.com/FedRAMP/community/discussions) 