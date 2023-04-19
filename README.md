## Usage

Define env variables to run exporter

```bash
export UNITD_CONTROL_NETWORK="tcp"
export UNITD_CONTROL_ADDRESS=":8081"
export METRICS_LISTEN_ADDRESS=":9094"
```

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

## Prometheus mertics prototipe
```bash
unit_instance_connections_accepted
unit_instance_connections_active
unit_instance_connections_idle
unit_instance_connections_closed

unit_instance_requests_total

unit_application_processes_running{application= }
unit_application_processes_starting{application= }
unit_application_processes_idle{application= }
unit_application_processes_requests_active{application= }
```

## Result metrics

```bash
# TYPE unit_application_processes gauge
unit_application_processes{application="laravel",instance="unit",state="idle"} 1
unit_application_processes{application="laravel",instance="unit",state="running"} 1
unit_application_processes{application="laravel",instance="unit",state="starting"} 0
unit_application_processes{application="myblog",instance="unit",state="idle"} 3
unit_application_processes{application="myblog",instance="unit",state="running"} 3
unit_application_processes{application="myblog",instance="unit",state="starting"} 0
# HELP unit_application_requests_active Similar to /status/requests, but includes only the data for a specific app.
# TYPE unit_application_requests_active gauge
unit_application_requests_active{application="unit",instance="laravel"} 0
unit_application_requests_active{application="unit",instance="myblog"} 0
# HELP unit_instance_connections_accepted Total accepted connections during the instance’s lifetime.
# TYPE unit_instance_connections_accepted counter
unit_instance_connections_accepted{application="",instance="unit"} 7
# HELP unit_instance_requests_total Total non-API requests during the instance’s lifetime.
# TYPE unit_instance_requests_total counter
unit_instance_requests_total{application="",instance="unit"} 7

```