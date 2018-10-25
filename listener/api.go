// Copyright 2015 ThoughtWorks, Inc.

// This file is part of getgauge/flash.

// getgauge/flash is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// getgauge/flash is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with getgauge/flash.  If not, see <http://www.gnu.org/licenses/>.

package listener

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"

	"reportportal-agent-gauge/event"
	m "reportportal-agent-gauge/gauge_messages"

	"github.com/golang/protobuf/proto"
)

type handlerFn func(*m.Message)

type apiListener struct {
	connection net.Conn
	handlers   map[m.Message_MessageType]handlerFn
	event      chan event.Event
}

func NewApiListener(host string, port string, e chan event.Event) (Listener, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}
	return &apiListener{connection: conn, handlers: map[m.Message_MessageType]handlerFn{
		m.Message_SuiteExecutionResult: func(msg *m.Message) {
			e <- event.NewEndEvent(msg.SuiteExecutionResult.SuiteResult.GetFailed())
		},
		m.Message_SpecExecutionStarting: func(msg *m.Message) {
			e <- event.NewSpecEvent(msg.SpecExecutionStartingRequest.CurrentExecutionInfo, true)
		},
		m.Message_SpecExecutionEnding: func(msg *m.Message) {
			e <- event.NewSpecEvent(msg.SpecExecutionEndingRequest.CurrentExecutionInfo, false)
		},
		m.Message_ScenarioExecutionStarting: func(msg *m.Message) {
			e <- event.NewScenarioEvent(msg.ScenarioExecutionStartingRequest.CurrentExecutionInfo, true)
		},
		m.Message_ScenarioExecutionEnding: func(msg *m.Message) {
			e <- event.NewScenarioEvent(msg.ScenarioExecutionEndingRequest.CurrentExecutionInfo, false)
		},
		m.Message_StepExecutionStarting: func(msg *m.Message) {
			e <- event.NewStepEvent(msg.StepExecutionStartingRequest.CurrentExecutionInfo, true)
		},
		m.Message_StepExecutionEnding: func(msg *m.Message) {
			e <- event.NewStepEvent(msg.StepExecutionEndingRequest.CurrentExecutionInfo, false)
		},
	}, event: e}, nil
}

func (l *apiListener) Start() {
	buffer := new(bytes.Buffer)
	data := make([]byte, 8192)
	for {
		n, err := l.connection.Read(data)
		if err != nil {
			return
		}
		buffer.Write(data[0:n])
		l.processMessages(buffer)
	}
}

func (l *apiListener) processMessages(buffer *bytes.Buffer) {
	var projectName = os.Getenv("REPORTPORTAL_PROJECT_NAME")
	var launchName = os.Getenv("REPORTPORTAL_LAUNCH_NAME")
	var launchID = ""
	var testItemID = ""
	data := os.Getenv("REPORTPORTAL_TAGS")
	var tags = strings.Split(data, ",")

	for {
		messageLength, bytesRead := proto.DecodeVarint(buffer.Bytes())
		if messageLength > 0 && messageLength < uint64(buffer.Len()) {
			message := &m.Message{}
			messageBoundary := int(messageLength) + bytesRead
			err := proto.Unmarshal(buffer.Bytes()[bytesRead:messageBoundary], message)
			if err != nil {
				log.Printf("Failed to read proto message: %s\n", err.Error())
			} else {
				if message.MessageType == m.Message_KillProcessRequest {
					inprogressLauchID := getListOfProjectLaunchesByFilter(projectName, "filter.eq.name="+url.QueryEscape(launchName)+"&filter.eq.status=IN_PROGRESS")
					finishLaunchForSpecifiedProject(projectName, inprogressLauchID, "FAILED", tags)

					l.connection.Close()
					<-l.event
					os.Exit(0)
				} else {
					h := l.handlers[message.MessageType]
					if h != nil {
						h(message)

						if message.GetExecutionStartingRequest() != nil {

						}

						if message.GetSpecExecutionStartingRequest() != nil {
							// Start a Launch for Specified Project in Report Portal
							inprogressLauchID := getListOfProjectLaunchesByFilter(projectName, "filter.eq.name="+url.QueryEscape(launchName)+"&filter.eq.status=IN_PROGRESS")
							// If no laucnch id then startsLaunchForSpecifiedProject
							if string(inprogressLauchID) == "0" {
								launchID = startsLaunchForSpecifiedProject(projectName, launchName, tags)
							} else {
								launchID = inprogressLauchID
							}

							// Gauge .spec data
							var sections []string
							fileNamePath := message.GetSpecExecutionStartingRequest().CurrentExecutionInfo.CurrentSpec.FileName
							if strings.Contains(fileNamePath, "\\") == true {
								// For Windows
								sections = strings.Split(fileNamePath, "\\")
							} else {
								// For Unix like env.
								sections = strings.Split(fileNamePath, "/")
							}

							var specName = sections[len(sections)-1]
							var specTags = message.GetSpecExecutionStartingRequest().CurrentExecutionInfo.CurrentSpec.Tags
							//var tag1 = strings.Join(specTags, ",")
							tags = specTags

							// Start a Root Test Item in Report Portal.
							testItemID = startARootTestItem(projectName, launchID, specName, tags, "SUITE")
						}

						if message.GetScenarioExecutionStartingRequest() != nil {
							// Gauge .spec data
							var sections []string
							fileNamePath := message.GetScenarioExecutionStartingRequest().CurrentExecutionInfo.CurrentSpec.FileName
							if strings.Contains(fileNamePath, "\\") == true {
								// For Windows
								sections = strings.Split(fileNamePath, "\\")
							} else {
								// For Unix like env.
								sections = strings.Split(fileNamePath, "/")
							}

							var specName = sections[len(sections)-1]
							var scenarioName = message.GetScenarioExecutionStartingRequest().CurrentExecutionInfo.CurrentScenario.Name

							testItemID = findTestItemsBySpecifiedFilter(projectName, "filter.eq.name="+url.QueryEscape(specName)+"&filter.eq.status=IN_PROGRESS")

							// Start a Child Test Item in Report Portal.
							startAchildTestItem(projectName, launchID, testItemID, scenarioName, tags, "STEP")

						}

						if message.GetStepExecutionStartingRequest() != nil {
							/*
								var currentStepName = message.GetStepExecutionStartingRequest().CurrentExecutionInfo.CurrentStep.Step.ActualStepText
								var tgs = message.GetStepExecutionStartingRequest().CurrentExecutionInfo.CurrentScenario.Tags
								var tags1 = strings.Join(tgs, ",")
								tags = launchName

								fmt.Print("getListOfProjectLaunchesByFilter | ")
								launchID := getListOfProjectLaunchesByFilter(projectName, "filter.eq.name="+url.QueryEscape(launchName)+"&filter.eq.status=IN_PROGRESS")
								//TODO: getTestItemID from Report
								fmt.Print("startAchildTestItem | ")
								var testItemIDChild = startAchildTestItem(projectName, launchID, testItemID, currentStepName, tags, "TEST")
								print(testItemIDChild)
							*/

						}

						if message.GetStepExecutionEndingRequest() != nil {
							var actualStepName = message.GetStepExecutionEndingRequest().CurrentExecutionInfo.CurrentStep.Step.ActualStepText
							//actualStepName = strings.Replace(actualStepName, "\"", "|", -1)
							//var currentStepName = message.GetStepExecutionEndingRequest().CurrentExecutionInfo.CurrentStep.Step.ParsedStepText
							var currentScenarioName = message.GetStepExecutionEndingRequest().CurrentExecutionInfo.CurrentScenario.Name
							//var tgs = message.GetStepExecutionEndingRequest().CurrentExecutionInfo.CurrentScenario.Tags
							//var tags1 = strings.Join(tgs, ",")
							//tags = launchName

							testItemIDChild := findTestItemsBySpecifiedFilter(projectName, "filter.eq.name="+url.QueryEscape(currentScenarioName)+"&filter.eq.status=IN_PROGRESS")

							// StatusID
							var logLevel = "info"
							if message.GetStepExecutionEndingRequest().CurrentExecutionInfo.CurrentStep.IsFailed {
								logLevel = "error"
								//var stepStackTrace = message.GetStepExecutionEndingRequest().CurrentExecutionInfo.
								// fmt.Print("logFileGauge | ")
								// logFileGauge := "/Users/gemunu/Documents/GitHub/gauge-python-init/logs/gauge.log"

								// fmt.Print("createLog | ")
								// var logID = createLog(projectName, "gauge.log", testItemIDChild, "debug", str)
								// fmt.Println(logID)

							}

							fileName := currentScenarioName
							createLog(projectName, fileName, testItemIDChild, logLevel, actualStepName)

						}

						if message.GetScenarioExecutionEndingRequest() != nil {
							var scenarioName = message.GetScenarioExecutionEndingRequest().CurrentExecutionInfo.CurrentScenario.Name
							scenarioTags := message.GetScenarioExecutionEndingRequest().CurrentExecutionInfo.CurrentScenario.Tags
							//tags = strings.Join(tgs, ",")
							tags = scenarioTags

							// StatusID
							var status = "IN_PROGRESS"
							if message.GetScenarioExecutionEndingRequest().CurrentExecutionInfo.CurrentScenario.IsFailed {
								status = "FAILED"
							} else {
								status = "PASSED"
							}

							testItemID := findTestItemsBySpecifiedFilter(projectName, "filter.eq.name="+url.QueryEscape(scenarioName)+"&filter.eq.status=IN_PROGRESS")
							finishTestItem(projectName, testItemID, scenarioName, status, tags)

						}

						if message.GetExecutionEndingRequest() != nil {

						}

						if message.GetSpecExecutionEndingRequest() != nil {
							// Gauge .spec data
							var sections []string
							fileNamePath := message.GetSpecExecutionEndingRequest().CurrentExecutionInfo.CurrentSpec.FileName
							if strings.Contains(fileNamePath, "\\") == true {
								// For Windows
								sections = strings.Split(fileNamePath, "\\")
							} else {
								// For Unix like env.
								sections = strings.Split(fileNamePath, "/")
							}

							var specName = sections[len(sections)-1]

							specTags := message.GetSpecExecutionEndingRequest().CurrentExecutionInfo.CurrentSpec.Tags
							//var tags1 = strings.Join(tgs, ",")
							tags = specTags

							// StatusID
							var status = "IN_PROGRESS"
							if message.GetSpecExecutionEndingRequest().CurrentExecutionInfo.CurrentSpec.IsFailed {
								status = "FAILED"
							} else {
								status = "PASSED"
							}
							testItemID = findTestItemsBySpecifiedFilter(projectName, "filter.eq.name="+url.QueryEscape(specName)+"&filter.eq.status=IN_PROGRESS")
							finishTestItem(projectName, testItemID, specName, status, tags)
						}

					}
				}
				buffer.Next(messageBoundary)
				if buffer.Len() == 0 {
					return
				}
			}
		} else {
			return
		}
	}
}
