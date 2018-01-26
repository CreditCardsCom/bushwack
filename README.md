# bushwack

*Why this name? Well, it seemed right at the time, and all logging related utilities seem to be named after the logging industry.*

This uses the new AWS provided [go lambda toolkit](https://github.com/aws/aws-lambda-go).

## Running


## Building

Building the lambda per aws-lambda-go documentation:


```bash
GOOS=linux go build -o main main.go
zip main.zip main
```

## Configuration

The following environment variables are needed to run.

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
