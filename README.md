# jx logging

A simple logging wrapper around logrus used in all JX components.  The aim is to have a single place to configure log formats and levels.

# Configuration

Configuration can be done with the following environment variables.

| Variable      | Options       | Default |
| ------------- | ------------- |---------|
| JX_LOG_LEVEL  | debug, info, warn | info |
| JX_LOG_FORMAT  | text, json, stackdriver | text |
