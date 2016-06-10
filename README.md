# goni
![beta-badge](https://img.shields.io/badge/release-beta-yellow.svg) [![wiki-badge](https://img.shields.io/badge/github-wiki-blue.svg)](https://github.com/goniapm/goni/wiki)

![logo](./resource/logo.png)

Goni는 Go언어를 위한 오픈소스 APM(Application Performance Management) 툴입니다.

## Overview
- **Dashboard** : 시간대별 Instance의 CPU (*Linux만 지원*) / Heap (*추후 지원 예정*) 사용량 히트맵을 클릭하면, 그 시간대의 User, Top 5 Instance, Top 5 Transaction, Transaction Detail을 보여줍니다.

- **Transaction Trace** : Go언어에서 `CallStack`을 가져오는데 한계가 있어, `CallStack`을 대체할 Transaction Trace 기능을 제공합니다. 이 기능을 통해 요청이 들어왔을 때 부터, 요청이 완료될 때 까지의 상태를 `Sankey Chart`로 보여줍니다.

- **Metric View** : Expvar / Runtime Metric 뿐만 아니라, 요청에 대한 다양한 Metric을 제공합니다.

## Architecture
* Client - Collector간의 통신 포맷으로 [protobuf](https://github.com/google/protobuf)를 사용합니다.
* Metric Data는 `timeseries data` 저장에 최적화된 [InfluxDB](https://influxdata.com/)에 저장합니다.
* 일반적인 정보(회원 데이터 / 프로젝트 설정)는 [MySQL](https://www.mysql.com/)에 저장합니다.
* Frontend는 [React](https://facebook.github.io/react/)를 사용합니다.

## Quickstart
[Wiki](https://github.com/goniapm/goni/wiki)를 참고해주세요 :D

## Issue
사용시 문제점 / 궁금하신 점이 있으시면 [여기](https://github.com/goniapm/goni/issues)에 이슈를 남겨주세요.

## Contribution
Goni를 개선해주세요! Contribution은 언제나 환영합니다 :D

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
