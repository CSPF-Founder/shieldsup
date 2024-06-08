package enums

// Type definition for BugTrackSeverity
type BugTrackSeverity int

// BugTrackSeverity enum values
const (
	BTSeverityCritical       BugTrackSeverity = 1
	BTSeverityHigh           BugTrackSeverity = 2
	BTSeverityMedium         BugTrackSeverity = 3
	BTSeverityLow            BugTrackSeverity = 4
	BtSeverityRecommendation BugTrackSeverity = 5
)

// Type alias with underlying type of IntEnumMap[BugTrackSeverity]
type BTSeverityMapType = IntEnumMap[BugTrackSeverity]

// BTSeverityMap is the map of BugTrackSeverity to string
var BTSeverityMap = BTSeverityMapType{
	BTSeverityCritical:       "Critical",
	BTSeverityHigh:           "High",
	BTSeverityMedium:         "Medium",
	BTSeverityLow:            "Low",
	BtSeverityRecommendation: "Recommendation",
}

// Type definition for BugTrackStatus
type BugTrackStatus int

const (
	BTStatusUnfixed BugTrackStatus = 1
	BTStatusFixed   BugTrackStatus = 3
)

// Type alias with underlying type of IntEnumMap[BugTrackStatus]
type BTStatusMapType = IntEnumMap[BugTrackStatus]

// BTStatusMap is the map of BugTrackStatus to string
var BTStatusMap = BTStatusMapType{
	BTStatusUnfixed: "Unfixed",
	BTStatusFixed:   "Fixed",
}

// Type definition for Prioritization
type Prioritization int

const (
	PrioritizationLow    Prioritization = 1
	PrioritizationMedium Prioritization = 2
	PrioritizationHigh   Prioritization = 3
)

// Type alias with underlying type of IntEnumMap[Prioritization]
type PrioritizationMapType = IntEnumMap[Prioritization]

// PrioritizationMap is the map of Prioritization to string
var PrioritizationMap = PrioritizationMapType{
	PrioritizationLow:    "Low",
	PrioritizationMedium: "Medium",
	PrioritizationHigh:   "High",
}

// Type definition for CanWafStop
type CanWafStop int

// CanWafStop enum values
const (
	CanWafStopNo            CanWafStop = 0
	CanWafStopYes           CanWafStop = 1
	CanWafStopNotApplicable CanWafStop = 2
)

// Type alias with underlying type of IntEnumMap[CanWafStop]
type CanWafStopMapType = IntEnumMap[CanWafStop]

// CanWafStopMap is the map of CanWafStop to string
var CanWafStopMap = CanWafStopMapType{
	CanWafStopNo:            "No",
	CanWafStopYes:           "Yes",
	CanWafStopNotApplicable: "Not Applicable",
}

// Type definition for DataLeakage
type DataLeakage int

// DataLeakage enum values
const (
	DataLeakageNo            DataLeakage = 0
	DataLeakageYes           DataLeakage = 1
	DataLeakageNotApplicable DataLeakage = 2
)

// Type alias with underlying type of IntEnumMap[DataLeakage]
type DataLeakageMapType = IntEnumMap[DataLeakage]

// DataLeakageMap is the map of DataLeakage to string
var DataLeakageMap = DataLeakageMapType{
	DataLeakageNo:            "No",
	DataLeakageYes:           "Yes",
	DataLeakageNotApplicable: "Not Applicable",
}

// Type definition for TestingMethod
type TestingMethod int

const (
	TestingMethodAutomatic TestingMethod = 1
	TestingMethodRedTeam   TestingMethod = 2
)

// Type alias with underlying type of IntEnumMap[TestingMethod]
type TestingMethodMapType = IntEnumMap[TestingMethod]

// TestingMethodMap is the map of TestingMethod to string
var TestingMethodMap = TestingMethodMapType{
	TestingMethodAutomatic: "Automatic",
	TestingMethodRedTeam:   "Red Teaming",
}

// Type definition for EffortsToExploit
type EffortsToExploit int

const (
	EffortsToExploitEasy          EffortsToExploit = 1
	EffortsToExploitModerate      EffortsToExploit = 2
	EffortsToExploitHard          EffortsToExploit = 3
	EffortsToExploitVeryHard      EffortsToExploit = 4
	EffortsToExploitNotApplicable EffortsToExploit = 5
)

// Type alias with underlying type of IntEnumMap[EffortsToExploit]
type EffortsToExploitMapType = IntEnumMap[EffortsToExploit]

// EffortsToExploitMap is the map of EffortsToExploit to string
var EffortsToExploitMap = EffortsToExploitMapType{
	EffortsToExploitEasy:          "Easy",
	EffortsToExploitModerate:      "Moderate",
	EffortsToExploitHard:          "Hard",
	EffortsToExploitVeryHard:      "Very Hard",
	EffortsToExploitNotApplicable: "Not Applicable",
}

// Type definition for Likelihood
type Likelihood int

const (
	LikelihoodVeryLow       Likelihood = 1
	LikelihoodLow           Likelihood = 2
	LikelihoodProbable      Likelihood = 3
	LikelihoodHigh          Likelihood = 4
	LikelihoodVeryHigh      Likelihood = 5
	LikelihoodNotApplicable Likelihood = 6
)

// Type alias with underlying type of IntEnumMap[Likelihood]
type LikelihoodMapType = IntEnumMap[Likelihood]

// LikelihoodMap is the map of Likelihood to string
var LikelihoodMap = LikelihoodMapType{
	LikelihoodVeryLow:       "Very Low",
	LikelihoodLow:           "Low",
	LikelihoodProbable:      "Probable",
	LikelihoodHigh:          "High",
	LikelihoodVeryHigh:      "Very High",
	LikelihoodNotApplicable: "Not Applicable",
}
