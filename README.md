# bushwack

*Why this name? Well, it seemed right at the time, and all logging related utilities seem to be named after the logging industry.*

## Overview

After using an S3 input in logstash it seemed be to having constant issues attempting to ingest logs from our S3 bucket we use for access logs. Additionality we didn't really want to shuffle around the data or rename all keys to satisfy an exclude pattern within logstash.

We realized lambda could accomplish the large task of indexing all logs into ElasticSearch in a *near* real time fashion (we want to index the logs as close to when AWS writes it into the bucket as possible).

#### Caveats

- We currently don't have any feature to retry indexing a log file if any type of failure happens.
- We are only parsing ALB logs, not legacy ELB logs.

## Setup

Ensure you have [dep](https://github.com/golang/dep) installed.

- `go get -d github.com/CreditCardsCom/bushwack`
- `dep ensure`
- `go test bushwack/`

## Building

This uses the new AWS provided [go lambda toolkit](https://github.com/aws/aws-lambda-go). Building the lambda per aws-lambda-go documentation:


```bash
GOOS=linux go build -o main main.go
zip main.zip main
```

## Configuration

At this point the project only supports being run on [AWS Lambda](https://aws.amazon.com/lambda/). The following environment variables are needed to run.

#### `ES_HOST`

Configures the Elasticsearch host that will be used to index the log output.

## License (MIT)

Copyright (c) 2018 CreditCards.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
