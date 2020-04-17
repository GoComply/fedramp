# FedRAMP - Automate Authorization Process
This is open source tool that manipulates official FedRAMP assets. Everyone is welcome to contribute!

## Features
 - take FedRAMP/OSCAL formatted System Security Plan and outputs FedRAMP document
 - take opencontrol repository and produce FedRAMP/OSCAL formatted System Security Plans

## User Resources
 - [FedRAMP adopts OSCAL Announcement](https://www.fedramp.gov/FedRAMP-moves-to-automate-the-authorization-process/)

## Developer Resources
 - [Guide to OSCAL-based FedRAMP System Security Plans](https://github.com/GSA/fedramp-automation/blob/master/documents/Guide_to_OSCAL-based_FedRAMP_System_Security_Plans.pdf)

## Exemplary usage

Build project (install golang as prerequisite)

```
make
```

Explore command-line UI

```
./fedramp --help
./fedramp opencontrol --hep
./fedramp convert --help
```

Covert [Open Control](https://open-control.org/) SSPs (in form of [masonry repository](https://github.com/opencontrol/compliance-masonry)) to OSCAL SSPs

```
./fedramp opencontrol https://github.com/ComplianceAsCode/redhat test_output/
```

Covert OSCAL SSP to DOCX Document

```
wget https://github.com/ComplianceAsCode/oscal/blob/master/xml/openshift-container-platform-4-fedramp-Low.xml
./fedramp convert ./openshift-container-platform-4-fedramp-Low.xml FedRAMP-Low.docx
```

This latest step is not fully complete as you can see, some of the fields in the DOCX being blank. This is work in progress.
