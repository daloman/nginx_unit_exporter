
## Unit stats meaning

```json
{
  "connections": {
    "accepted": 317891, # Counter - _total Integer; total accepted connections during the instance’s lifetime.
    "active": 0,        # Gauge   - Integer; current active connections for the instance.
    "idle": 2,          # Gauge   - Integer; current idle connections for the instance.
    "closed": 317889    # Counter - Integer; total closed connections during the instance’s lifetime.
  },
  "requests": {
    "total": 381513     # Counter - Integer; total non-API requests during the instance’s lifetime.
  },
  "applications": {
    "laravel": {
      "processes": {
        "running": 1,   # Gauge   - Integer; current running app processes.
        "starting": 0,  # Gauge   - Integer; current starting app processes.
        "idle": 1       # Gauge   - Integer; current idle app processes.
      },
      "requests": {
        "active": 0     # Gauge   - Integer; similar to /status/requests, but includes only the data for a specific app.
      }
    }
  }
}
```

unit_instance_connections_ 
unit_instance_requests_total

unit_application_ process_ __
unit_application_ requests  active