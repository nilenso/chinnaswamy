# Spec

## Overview

Our goal is to build a URL shortener that can handle some amount of burst-y
traffic. This is a toy project for us to get some recent Go experience under
our belts.

## Functional Requirements

1. Users of the API must be able to create short urls when they provide a long URL and an expiry time
2. A request to an unexpired short URL should redirect
3. Must send notifications whenever a shorten url is visited over a queue mechanism
4. Short URLs can be in-expirable

## Non-functional Requirements

__Latency__

- 99.9%ile url shorten requests must respond within 10ms at up to 2000 rps
  - Loosely defined
  - Assuming the requests are spread evenly across the 1s period
  - This should mean that about 20 requests are made every 10ms

__Reliability__

- Should be resilient to any one machine going down

__Timing & Delays__

- The short url returned in a successful response to a url shorten call must be
  immediately usable. Note that it would be better design to allow for some
  async processing to happen and define a time period after which the link will
  be usable (f.ex. 1s). But we're not doing that for the simplicity it affords
  the consumer.
- The notification may be sent upto 30 seconds after the corresponding short
  URL is visited.
