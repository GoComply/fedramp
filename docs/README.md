# FedRAMP Automation Tools - R5 Balance & 20x Implementation

## Overview

This repository extends the [gocomply/fedramp](https://github.com/gocomply/fedramp) project with implementations for:

- **FedRAMP R5 Balance Initiatives** - Modern improvements to the FedRAMP process
- **FedRAMP 20x Phase One** - Key Security Indicators (KSI) validation framework
- **REST API Server** - Complete API for all FedRAMP operations
- **Continuous Monitoring** - Real-time compliance validation

## What's Implemented

### ✅ Complete Features

1. **R5 Balance Components**
   - Significant Change Notifications (SCN) - RFC-0007 compliant
   - Continuous Reporting Standard (CRS) - 6 key metrics
   - Minimum Assessment Standard (MAS) - Assessment framework
   - Storing & Sharing Authorization Data (SSAD) - Package management

2. **FedRAMP 20x Phase One**
   - All 11 Key Security Indicators per release 25.05C
   - 51 sub-requirements with validation logic
   - Continuous reporting proposal generation
   - Machine-readable JSON output

3. **Infrastructure**
   - REST API with all endpoints
   - Database schema and operations
   - Continuous monitoring framework
   - Alert management system
   - Web dashboard

### 🚧 Work In Progress

- Security Assessment Report (SAR) - Basic structure only
- Plan of Action & Milestones (POA&M) - Basic structure only
- Security Assessment Plan (SAP) - Basic structure only

### ❌ Not Implemented

Traditional FedRAMP documents like IRP, CP, CMP, PIA, etc. See `docs/MISSING_REV5_COMPONENTS.md` for full list.

## Quick Start

```bash
# Build the server
go build -mod=mod -o fedramp-server cmd/server/main.go

# Run the server
./fedramp-server

# Check health
curl http://localhost:8080/api/v1/health
```

## API Examples

```bash
# Validate KSIs
curl -X POST http://localhost:8080/api/v1/ksi/validate \
  -H "Content-Type: application/json" \
  -d '{"csoId": "CSO-001", "evidence": {}}'

# Create SCN
curl -X POST http://localhost:8080/api/v1/scn \
  -H "Content-Type: application/json" \
  -d '{
    "service_offering_id": "CSO-001",
    "change_type": "new functionality",
    "short_description": "Adding new API endpoint",
    "reason_for_change": "Customer requirement"
  }'

# Get metrics
curl http://localhost:8080/api/v1/crs/metrics/CSO-001
```

## CLI Usage

```bash
# Original functionality
gocomply_fedramp convert ssp.xml output.docx

# New R5 Balance commands
gocomply_fedramp scn create --service CSO-001 --type adaptive
gocomply_fedramp ksi validate --service CSO-001
gocomply_fedramp mas create --type initial
```

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌───────────────┐
│   CLI Tool  │────▶│  API Server  │────▶│   Database    │
└─────────────┘     └──────────────┘     └───────────────┘
                            │
                    ┌───────┴────────┐
                    ▼                ▼
            ┌──────────────┐  ┌──────────────┐
            │  Monitoring  │  │  Dashboard   │
            └──────────────┘  └──────────────┘
```

## Development

```bash
# Run tests
go test ./...

# Run with debug logging
./fedramp-server --debug

# Use Docker Compose for full stack
docker-compose up -d
```

## Contributing

1. Check `docs/IMPLEMENTATION_STATUS.md` for current state
2. Mark WIP code clearly with TODO comments
3. Update documentation when adding features
4. Follow existing patterns for consistency

## License

Same as parent project - see LICENSE.md

## Acknowledgments

Built on top of [gocomply/fedramp](https://github.com/gocomply/fedramp) by the GoComply team. 