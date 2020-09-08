# FedRAMP - Automate Authorization Process ![Build CI](https://github.com/gocomply/fedramp/workflows/Build%20CI/badge.svg) [![Gitter](https://badges.gitter.im/GoComply/community.svg)](https://gitter.im/GoComply/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
This is open source tool that manipulates official FedRAMP assets. Everyone is welcome to contribute!

## Features
 - take FedRAMP/OSCAL formatted System Security Plan and outputs FedRAMP document
 - take opencontrol repository and produce FedRAMP/OSCAL formatted System Security Plans

## User Resources
 - [Additional FedRAMP OSCAL Resources and Templates](https://www.fedramp.gov/additional-fedramp-oscal-resources-and-templates/) (August 20, 2020)
 - [FedRAMP adopts OSCAL Announcement](https://www.fedramp.gov/FedRAMP-moves-to-automate-the-authorization-process/) (December 17, 2019)

## Developer Resources
 - [Guide to OSCAL-based FedRAMP System Security Plans](https://github.com/GSA/fedramp-automation/blob/master/documents/Guide_to_OSCAL-based_FedRAMP_System_Security_Plans.pdf)

## Exemplary usage

Build project (install golang as prerequisite)

```
go get -u -v github.com/gocomply/fedramp/cli/gocomply_fedramp
```

Explore command-line UI

```
gocomply_fedramp --help
gocomply_fedramp opencontrol --hep
gocomply_fedramp convert --help
```

Covert [Open Control](https://open-control.org/) SSPs (in form of [masonry repository](https://github.com/opencontrol/compliance-masonry)) to OSCAL SSPs

```
gocomply_fedramp opencontrol https://github.com/ComplianceAsCode/redhat test_output/
```

Covert OSCAL SSP to DOCX Document

```
wget https://raw.githubusercontent.com/ComplianceAsCode/oscal/master/xml/openshift-container-platform-4-fedramp-Low.xml
gocomply_fedramp convert ./openshift-container-platform-4-fedramp-Low.xml FedRAMP-Low.docx
```

This latest step is not fully complete as you can see, some of the fields in the DOCX being blank. This is work in progress.
