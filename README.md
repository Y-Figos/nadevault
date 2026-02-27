# Nadevault: A CS2 lineups knowledge base [WIP]

## Project Goals
Nadevault is made to be a fast and simple database for CS2 lineups. The main goal is: players can quickly search, view, and learn a lineup in 15 seconds (the Freeze Time period at the beginning of each round). No bloat with heavy JS interactive maps or video players just to see crosshair alignment.

> Above all things, Nadevault is a learning experience to apply real-world architecture and cloud-native development.

## Stack
The thought process for the main stack choice was straightforward. It needs to be performant and simple. Go + Templ + HTMX fits the description with ease. 

### Infrastructure
The infrastructure needed for running Nadevault is PostgreSQL + Minio + Redis. 
- **PostgreSQL:** Main database for lineups.
- **Minio:** Stores all the lineup screenshots as public indexable URLs.
- **Redis:** A little spicy to reduce load on the main Nadevault API, used for queueing a WebP conversion worker. 

## Architecture
Nadevault has a simple architecture aimed to be easily run in a cloud-native environment. 
A single API that serves SSR pages with HTMX and Templ.
Minio and Redis are used as a way to reduce load from the main API.

Images are first uploaded to Minio with pre-signed PUTs. A job is sent to the Redis stream and a Go worker is responsible for the conversion to WebP and saving it in the public bucket, reducing load times on the frontend and storage space/cost.

## Roadmap v1
- [x] Basic Nadevault API
- [ ] Structured Logging
- [ ] Minio Upload Cycle
- [ ] Redis Streams and Go webp Worker
- [ ] Lineups visualizations Page
- [ ] Add Auth and Ratelimiting
- [ ] Observability Stack
> Note: Error handling and context propagation need to be improved.
