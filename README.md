# prometheus-multi-tenant-proxy-server

Tiny reverse proxy server to make Prometheus multi-tenant. This reverse proxy provides a little Authentication/Authorization layer on top of Prometheus and injects the project label into all user queries.

## Architecture

Here is the big picture of the architecture to monitor Kubernetes cluster in a multi-tenant model.

![prometheus-multi-tenant-proxy-server](diagram.png)

Copyright 2023 Saeid Bostandoust <ssbostan@yahoo.com>
