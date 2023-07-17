## svc-temporal-frontend.hcl
service {
  name = "temporal-frontend"
  id = "temporal-frontend-1"
  tags = ["v1"]
  port = 7233

  check {
    id =  "check-temporal-frontend",
    name = "Product temporal-frontend status check",
    service_id = "temporal-frontend-1",
    args  = ["/usr/local/bin/temporal-healthchecker_linux_amd64", "--config", "/ops/healthchecker.json"],
    interval = "5s",
    timeout = "10s"
  }
}

## svc-temporal-history.hcl
service {
  name = "temporal-history"
  id = "temporal-history-1"
  tags = ["v1"]
  port = 7233

  check {
    id =  "check-temporal-history",
    name = "Product temporal-history status check",
    service_id = "temporal-history-1",
    args  = ["/usr/local/bin/temporal-healthchecker_linux_amd64", "--config", "/ops/healthchecker.json"],
    interval = "5s",
    timeout = "10s"
  }
}

## svc-temporal-matching.hcl
service {
  name = "temporal-matching"
  id = "temporal-matching-1"
  tags = ["v1"]
  port = 7233

  check {
    id =  "check-temporal-matching",
    name = "Product temporal-matching status check",
    service_id = "temporal-matching-1",
    args  = ["/usr/local/bin/temporal-healthchecker_linux_amd64", "--config", "/ops/healthchecker.json"],
    interval = "5s",
    timeout = "10s"
  }
}