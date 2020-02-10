package common

type BaselineLevel int

const (
	LevelUnknown = iota
	LevelLow
	LevelModerate
	LevelHigh
)

var ProfileUrls = map[BaselineLevel]string{
	// As defined in APPENDIX A.
	LevelLow:      "https://raw.githubusercontent.com/usnistgov/OSCAL/master/content/fedramp.gov/xml/FedRAMP_LOW-baseline_profile.xml",
	LevelModerate: "https://raw.githubusercontent.com/usnistgov/OSCAL/master/content/fedramp.gov/xml/FedRAMP_MODERATE-baseline_profile.xml",
	LevelHigh:     "https://raw.githubusercontent.com/usnistgov/OSCAL/master/content/fedramp.gov/xml/FedRAMP_HIGH-baseline_profile.xml",
}

func (l BaselineLevel) Name() string {
	switch l {
	case LevelLow:
		return "Low"
	case LevelModerate:
		return "Moderate"
	case LevelHigh:
		return "High"
	}
	return "unknown"
}
