# jx logging

[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x/jx-logging)](https://goreportcard.com/report/github.com/jenkins-x/jx-logging)
[![Downloads](https://img.shields.io/github/downloads/jenkins-x/jx-logging/total.svg)](https://github.com/jenkins-x/jx-logging/releases)

A simple logging wrapper around logrus used in all JX components.  The aim is to have a single place to configure log formats and levels.

# Configuration

Configuration can be done with the following environment variables.

| Variable      | Options       | Default |
| ------------- | ------------- |---------|
| JX_LOG_LEVEL  | debug, info, warn | info |
| JX_LOG_FORMAT | text, json, stackdriver | text |
| JX_LOG_FILE    | a location to send debug logs to | |
| JX_LOG_SERVICE | the service name | |
| JX_LOG_SERVICE_VERSION | the service version | |
| JX_LOG_STACK_SKIP | the comma separated stack frames to skip in the logs | |

