package listener

import (
	"log"
	"os"
	"reportportal-agent-gauge/reportportal"
)

func reportPortal() *reportportal.Client {
	reportPortalServer := os.Getenv("REPORTPORTAL_SERVER")
	reportPortalUUID := os.Getenv("REPORTPORTAL_UUID")
	client := reportportal.NewClient(reportPortalServer, reportPortalUUID)
	return client
}

// getListOfProjectLaunchesByFilter
func getListOfProjectLaunchesByFilter(projectName string, filterQuery string) string {

	var filterdLaunchID = "0"
	launchFilter, err := reportPortal().GetListOfProjectLaunchesByFilter(projectName, filterQuery)

	if err != nil {
		log.Fatalf(err.Error())
	}

	if launchFilter.Page.TotalElements != 0 {
		filterdLaunchID = launchFilter.Content[0].ID
	}

	return filterdLaunchID
}

// findTestItemsBySpecifiedFilter specific tests by filter
func findTestItemsBySpecifiedFilter(projectName string, filterQuery string) string {
	var filterdTestID = "0"
	testFilter, err := reportPortal().FindTestItemsBySpecifiedFilter(projectName, filterQuery)

	if err != nil {
		log.Fatalf(err.Error())
	}

	if testFilter.Page.TotalElements != 0 {
		filterdTestID = testFilter.Content[0].ID
	}

	return filterdTestID
}

// startsLaunchForSpecifiedProject
func startsLaunchForSpecifiedProject(projectName string, launchName string, tags []string) string {

	launch, err := reportPortal().StartsLaunchForSpecifiedProject(projectName, launchName, tags)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return launch.ID
}

// finishLaunchForSpecifiedProject
func finishLaunchForSpecifiedProject(projectName string, launchID string, status string, tags []string) string {
	launch, err := reportPortal().FinishLaunchForSpecifiedProject(projectName, launchID, status, tags)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return launch.ID
}

// startARootTestItem
func startARootTestItem(projectName string, launchID string, testItemName string, tags []string, testType string) string {

	launch, err := reportPortal().StartARootTestItem(projectName, launchID, testItemName, tags, testType)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return launch.ID

}

// finishTestItem
func finishTestItem(projectName string, testItemID string, testItemName string, status string, tags []string) string {
	launch, err := reportPortal().FinishTestItem(projectName, testItemID, testItemName, status, tags)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return launch.ID

}

// startAchildTestItem
func startAchildTestItem(projectName string, launchID string, parentItem string, testItemName string, tags []string, testType string) string {

	launch, err := reportPortal().StartAchildTestItem(projectName, launchID, parentItem, testItemName, tags, testType)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return launch.ID

}

// createLog
func createLog(projectName string, fileName string, itemID string, logLevel string, message string) string {

	logRP, err := reportPortal().CreateLog(projectName, fileName, itemID, logLevel, message)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return logRP.ID

}
