# Temporal Healthchecker

This project is a health check tool for Temporal services. It provides functionality to check the health of the
FrontendService, HistoryService, and MatchingService of a Temporal server.

## Usage

- Download release binary from [here](https://github.com/saga420/temporal-healthchecker/releases).

- Check sha256 checksum of the downloaded binary.

- Create a config.json file following the structure described in the Configuration section.

- Run the healthchecker:

```bash
./temporal-healthchecker --config=path_to_your_config.json
```

---------

You can also build the binary from source:

```bash
git clone git@github.com:saga420/temporal-healthchecker.git
cd temporal-healthchecker
go build
./temporal-healthchecker --config=path_to_your_config.json
```

## Configuration

The configuration is done via a config.json file. Here is an example of the file:

```json 
{
  "FrontendService": {
    "IsEnabled": true,
    "Address": "127.0.0.1:7233",
    "TimeOut": 10
  },
  "HistoryService": {
    "IsEnabled": true,
    "Address": "127.0.0.1:7234",
    "TimeOut": 10
  },
  "MatchingService": {
    "IsEnabled": true,
    "Address": "127.0.0.1:7235",
    "TimeOut": 10
  }
}
```

Results of the health check will be logged to the console.

``` 
2023/07/17 06:06:52 Health Checker GitRevision 6943229
2023/07/17 06:06:52 Checking Temporal frontend service health
2023/07/17 06:06:52 Checking Temporal frontend service health done. Time elapsed: 1.474429ms
2023/07/17 06:06:52 Checking Temporal history service health
2023/07/17 06:06:52 Checking Temporal history service health done. Time elapsed: 379.883µs
2023/07/17 06:06:52 Checking Temporal matching service health
2023/07/17 06:06:52 Checking Temporal matching service health done. Time elapsed: 342.634µs
2023/07/17 06:06:52 full check started at 2023-07-17 06:06:52.983942753 +0000 UTC m=+0.004455329
2023/07/17 06:06:52 initialFullCheckClients took 74.801µs
2023/07/17 06:06:52 Cluster Id: 0096a900-3ad0-4de7-83f6-c4f0b29bf591
2023/07/17 06:06:52 Version Info: nil
2023/07/17 06:06:52 Cluster Name: active
2023/07/17 06:06:52 History Shard Count: 300
2023/07/17 06:06:52 checkClusterInfo took 2.366526ms
2023/07/17 06:06:52 Capabilities ActivityFailureIncludeHeartbeat: true
2023/07/17 06:06:52 Capabilities SdkMetadata: true
2023/07/17 06:06:52 Capabilities BuildIdBasedVersioning: true
2023/07/17 06:06:52 Capabilities UpsertMemo: true
2023/07/17 06:06:52 ServerVersion: 1.22.0
2023/07/17 06:06:52 checkSystemInfo took 2.834694ms
2023/07/17 06:06:52 Namespace: testnamespace, State: Registered, Description: This is my namespace
2023/07/17 06:06:52 Namespace: temporal-system, State: Registered, Description: Temporal internal system namespace
2023/07/17 06:06:52 Namespace: testnamespace1, State: Registered, Description: This is my namespace
2023/07/17 06:06:52 Namespace: testnamespace2, State: Registered, Description: This is my namespace
2023/07/17 06:06:52 Namespace: testnamespace3, State: Registered, Description: This is my namespace
2023/07/17 06:06:52 checkNamespaces took 4.086264ms
2023/07/17 06:06:52 Cluster 0: active, 0096a900-3ad0-4de7-83f6-c4f0b29bf591
2023/07/17 06:06:52 Cluster 0 is connected
2023/07/17 06:06:52 checkListClusters took 5.094278ms
2023/07/17 06:06:52 full check finished at 2023-07-17 06:06:52.989041249 +0000 UTC m=+0.009553825
2023/07/17 06:06:52 All services are healthy.
```

## Documentation

This tool provides two main checks: BasicCheck and FullCheck.

If any of these checks fail, the program will log the errors and exit with a status code of

1. If all checks pass, the program logs that all services are healthy and exits with a status code of 0.

The BasicCheck checks if the configured services are reachable, while the FullCheck performs a deeper health check (
implementation depends on your actual requirements).

The tool is also designed to be compatible with Docker and Consul health check requirements.

### Consul (see examples/*)

```hcl
## svc-temporal-frontend.hcl
service {
  name = "temporal-frontend"
  id   = "temporal-frontend-1"
  tags = ["v1"]
  port = 7233

  check {
    id         = "check-temporal-frontend",
    name       = "Product temporal-frontend status check",
    service_id = "temporal-frontend-1",
    args       = ["/usr/local/bin/temporal-healthchecker_linux_amd64", "--config", "/ops/healthchecker.json"],
    interval   = "5s",
    timeout    = "10s"
  }
}
```

### Docker (see examples/*)

```yaml
healthcheck:
  test: [ "CMD", "/usr/local/bin/temporal-healthchecker_linux_amd64", "--config", "/ops/healthchecker.json" ]
  interval: 5s
  timeout: 10s
  retries: 12
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT