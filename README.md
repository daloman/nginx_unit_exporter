# Unit Prometheus Exporter

This code aims to export in Prometheus format [usage statistics](https://unit.nginx.org/usagestats/)
available via the GET-only /status section of [Unit control API](https://unit.nginx.org/controlapi/#configuration-api).

## Usage

One can run nginx_unit_exporter as: 
 - sidecar container within same [Pod](https://kubernetes.io/docs/concepts/workloads/pods/) 
with Unit container.
 - separate Docker container
 - Systemd service
 - any other manner

Buid app.

Define env variables to run exporter

|Variable name             |Default value | Example                                                                  |
|--------------------------|--------------|--------------------------------------------------------------------------|
|UNITD_CONTROL_NETWORK     |  tcp         | `tcp` OR `unix` ... see [net.Dial examples](https://pkg.go.dev/net#Dial) |
|UNITD_CONTROL_ADDRESS     |  :8081       | `127.0.0.1:8081` OR `/var/run/control.unit.sock`                         |
|METRICS_LISTEN_ADDRESS    |  :9095       | `127.0.0.1:9095`                                                         |

Run nginx_unit_exporter.

Example for bash:
```bash
export UNITD_CONTROL_NETWORK="tcp"
export UNITD_CONTROL_ADDRESS=":8081"
export METRICS_LISTEN_ADDRESS=":9095"
```

## Unit stats meaning

Read the docs https://unit.nginx.org/usagestats/
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

## Result metrics example

```bash
# HELP unit_application_processes Current app processes.
# TYPE unit_application_processes gauge
unit_application_processes{application="laravel",instance="unit",state="idle"} 1
unit_application_processes{application="laravel",instance="unit",state="running"} 1
unit_application_processes{application="laravel",instance="unit",state="starting"} 0
unit_application_processes{application="myblog",instance="unit",state="idle"} 3
unit_application_processes{application="myblog",instance="unit",state="running"} 3
unit_application_processes{application="myblog",instance="unit",state="starting"} 0
# HELP unit_application_requests_active Similar to /status/requests, but includes only the data for a specific app.
# TYPE unit_application_requests_active gauge
unit_application_requests_active{application="laravel",instance="unit"} 0
unit_application_requests_active{application="myblog",instance="unit"} 0
# HELP unit_instance_connections_accepted_total Total accepted connections during the instance’s lifetime.
# TYPE unit_instance_connections_accepted_total counter
unit_instance_connections_accepted_total{application="",instance="unit"} 2
# HELP unit_instance_connections_active Current active connections for the instance
# TYPE unit_instance_connections_active gauge
unit_instance_connections_active{application="",instance="unit"} 0
# HELP unit_instance_connections_closed_total Total closed connections during the instance’s lifetime
# TYPE unit_instance_connections_closed_total counter
unit_instance_connections_closed_total{application="",instance="unit"} 2
# HELP unit_instance_connections_idle Current idle connections for the instance
# TYPE unit_instance_connections_idle gauge
unit_instance_connections_idle{application="",instance="unit"} 0
# HELP unit_instance_requests_total Total non-API requests during the instance’s lifetime.
# TYPE unit_instance_requests_total counter
unit_instance_requests_total{application="",instance="unit"} 2

```