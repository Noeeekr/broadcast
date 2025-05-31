package instance

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
)

type Logger struct{}

// NewLogger returns a Logger struct with methods to output
func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) TrackedLog(LogType, LogReason, Description string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "Unknown"
		line = 0
	}

	fmt.Printf("(%s) Was caused because of %s.\n", LogType, LogReason)
	fmt.Printf("Path: %s:%d \n", file, line)
	fmt.Printf("Decription: %s \n", Description)
}

func (logger *Logger) SimpleLog(LogType, LogEvent, Data string) {
	fmt.Printf("(%s) %s.\n", LogType, LogEvent)
	if Data != "" {
		fmt.Printf("Data: %s \n", Data)
	}
}

func (logger *Logger) ErrorLog(reason, description string) {
	fmt.Printf("########## Start of error log ##########\n")
	logger.TrackedLog("Error", reason, description)
	fmt.Printf("########## Ending of error log ##########\n")
}

// LogStep is called by LogSteps internally to log a singular job step
func (logger *Logger) logStep(jobName, jobStageDescription, jobID string, jobStage int8) {
	fmt.Printf("(JOB) NAME: %s - ID: %s\n(JOB) Name: %s - Current stage: %d\n (JOB) Name: %s - Description: %s\n",
		jobName,
		jobID,
		jobName,
		jobStage,
		jobName,
		jobStageDescription,
	)
}

// NewStepLogger returns LogStep.
// When LogStep is called, a message is sent to the specified stdout with a specific identifier
func (logger *Logger) NewStepLogger(jobName, jobStageDescription string) (LogStep func()) {
	var stepCounter int8 = 0

	var ID string = ""
	for range 10 {
		n, err := rand.Int(nil, big.NewInt(127))
		if err != nil {
			ID += string(n.Bytes())
		} else {
			ID += "X"
		}
	}

	fmt.Printf("(JOB) Starting %s\n", jobName)
	return func() {
		stepCounter++

		logger.logStep(jobName, jobStageDescription, ID, stepCounter)
	}
}
