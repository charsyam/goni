# goniplus-worker
Data handler / worker for Goni+

[![Go Report Card](https://goreportcard.com/badge/github.com/goniapm/goniplus-worker)](https://goreportcard.com/report/github.com/goniapm/goniplus-worker)

## Before you start
goniplus-worker require databases with settings. For more information, please visit [Wiki](https://github.com/goniapm/goni/wiki)

- InfluxDB (v0.9.6.1)
- MySQL or MariaDB

Also, You **must** set ENV variables listed below.

- `GONI_MYSQL_HOST`
  - e.g. `127.0.0.1`
- `GONI_MYSQL_PORT`
  - e.g. `3306`
- `GONI_MYSQL_USER`
  - e.g. `username`
- `GONI_MYSQL_PASS`
  - e.g. `password`
- `GONI_INFLUX_HOST`
  - e.g. `http://127.0.0.1:8086`
  - **Protocol and Port must be included**
- `GONI_INFLUX_USER`
  - e.g. `username`
- `GONI_INFLUX_PASS`
  - e.g. `password`


  ## License
  ```
  The MIT License (MIT)

  Copyright (c) 2016 Goni

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
  ```
