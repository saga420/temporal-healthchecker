# Temporal Healthchecker

This project is a health check tool for Temporal services. It provides functionality to check the health of the
FrontendService, HistoryService, and MatchingService of a Temporal server.

## Usage

- Clone this repository to your local machine:

```bash 
git clone https://github.com/saga420/temopral-healthchecker.git
```

- Build the project:

```bash
cd temporal-healthchecker
go build
```

- Create a config.json file following the structure described in the Configuration section.

- Run the healthchecker:

```bash
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
2023/07/17 11:31:42 NewHealthChecker took 740.129µs
2023/07/17 11:31:42 BasicCheck is starting...
2023/07/17 11:31:43 BasicCheck is done.
2023/07/17 11:31:43 BasicCheck took 947.428265ms
2023/07/17 11:31:43 FullCheck is starting...
2023/07/17 11:31:43 Cluster Id: 0096a900-3ad0-4de7-83f6-c4f0b29bf591
2023/07/17 11:31:43 Version Info: nil
2023/07/17 11:31:43 Cluster Name: active
2023/07/17 11:31:43 History Shard Count: 300
2023/07/17 11:31:44 Capabilities ActivityFailureIncludeHeartbeat: true
2023/07/17 11:31:44 Capabilities SdkMetadata: true
2023/07/17 11:31:44 Capabilities BuildIdBasedVersioning: true
2023/07/17 11:31:44 Capabilities UpsertMemo: true
2023/07/17 11:31:44 ServerVersion: 1.22.0
2023/07/17 11:31:44 Namespace: testnamespace, State: Registered, Description: This is my namespace
2023/07/17 11:31:44 Namespace: temporal-system, State: Registered, Description: Temporal internal system namespace
2023/07/17 11:31:44 Namespace: testnamespace1, State: Registered, Description: This is my namespace
2023/07/17 11:31:44 Namespace: testnamespace2, State: Registered, Description: This is my namespace
2023/07/17 11:31:44 Namespace: testnamespace3, State: Registered, Description: This is my namespace
2023/07/17 11:31:44 Cluster 0: active, 0096a900-3ad0-4de7-83f6-c4f0b29bf591
2023/07/17 11:31:44 Cluster 0 is connected
2023/07/17 11:31:44 FullCheck is done.
2023/07/17 11:31:44 FullCheck took 2.171620989s
2023/07/17 11:31:44 All services are healthy.
```

## Documentation

This tool provides two main checks: BasicCheck and FullCheck.

If any of these checks fail, the program will log the errors and exit with a status code of

1. If all checks pass, the program logs that all services are healthy and exits with a status code of 0.

The BasicCheck checks if the configured services are reachable, while the FullCheck performs a deeper health check (
implementation depends on your actual requirements).

The tool is also designed to be compatible with Docker and Consul health check requirements.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT