# Particle Events Prometheus Exporter

This will accept default webhooks from Particle Cloud and make them available for scraping by a Prometheus server.

## Setup

Set an API Key with the env var `API_KEY` .  This allows you to add an Authorization Header to your webhooks.
Optionally provide a Particle Cloud access token to resolve device ids to device names.  Set this by adding a `PARTICLE_ACCESS_TOKEN` env var.

## Current Metrics

Set your event names to one of the following metrics and it will automatically be exported.  Values will be converted to float64.

```
temperature
humidity
dust
sound
```

A label will be added to each metric of the device name if a access token was provided or the device id if it wasn't.