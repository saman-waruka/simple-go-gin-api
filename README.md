# This project for demonstrate implement simple GO GIN API
## with integrate monitoring tools for e.g. Sentry.io, Prometheus, Grafana, Zap logging

### To start this project on local

1. Start docker compose for run Prometheus and Grafana in local via Docker
    ```bash
    docker compose up -d
    ```

2. Install dependency 
    ```bash
    go mod download
    ```

3. Create environment variable in ```.env``` 
    ```bash
    cp .env.example .env
    ```
    and Fill in data.
    <b>Note:</b> for <b>SentryDNS</b> you can receive by register at Sentry.io for trial for test.
<br/>


4. Start GO project
    ```bash
    go run main.go
    ```

[Optional] Run in develop mode with ```Air``` for live reload when files in project changes
1. Install Air
    ```bash 
    go install github.com/air-verse/air@latest
    ```

2. Start project with live reload
    ```bash
    air
    ```

---
Swagger 
To access swagger url  http://localhost:8216/swagger/index.html

To generate swagger description 
1. Add comment at handler (example see ```user_handler.go``` on function GetUsers and CreateUser)
2. run command to generate swagger file (JSON, yml)
    ```bash
    swag init
    ```
    This process will re-create files at folder ```docs```  (```docs.go```, ```swagger.json```, ```swagger.yaml```)