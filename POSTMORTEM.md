# Postmortem

## Requirements for running the tests

* Install Golang 1.18.x (make sure its in your $PATH)
* Install Java 11.x
* Clone this repo
* In your terminal
* `$ cd /path/to/repo/tests`
* `$ go test`

## Discovering the CLI's boundaries

The goal here was to discover the boundaries of the CLI. With this in mind, I did the following.

* Use yaml file with 100s of OIDs
* Try to pass unexpected file types (json, xml, jpeg)
* Point to non-existent files
* Use bad yaml data
* Use invalid inputs

## Approach

* Learn enough Golang to be somehwat proficient 😅😅 (duration: one day)
* Exploratory testing and discovery (duration: 1 hour)
* Build automated test scripts using Golang (duration 3 hours)

## Issues with recommendations for resolution

* Using OID values in the yaml file that are not strings makes Java throw type error. I'd recommend catching this error and displaying a more user-readable error message.

## Things I wanted to do, but got cut for time

* Dockerize the project
* Benchmark the CLI's performance under varying conditions
* Enhance the test reporting
* Created GitHub Actions config to build, deploy, and test in a CI/CD environment
