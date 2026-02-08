# connector-sdk-go

SDK base para construir conectores em Go.

## O que inclui

- Interface `Connector` com operações de metadata, validação, descoberta, start/stop e execução de comando.
- `KafkaPublisher` com serialização de evento (`gen/go/platformschemas`), headers padrão e retry/backoff simples.
- Servidor HTTP embutido opcional com:
  - `GET /healthz`
  - `GET /metrics` (Prometheus client)
- Logging estruturado via `log/slog`.
- Configuração via variáveis de ambiente.
- Diretório `template/minimal` com conector mínimo compilável.
- GitHub Actions CI com checks de formatação, testes, lint e publicação opcional da imagem builder no GHCR.

## Instalação

```bash
go mod tidy
```

## Variáveis de ambiente

- `CONNECTOR_LOG_LEVEL` (default: `INFO`)
- `CONNECTOR_KAFKA_BROKERS` (default: `localhost:9092`)
- `CONNECTOR_KAFKA_TOPIC` (default: `connector.events`)
- `CONNECTOR_HTTP_ADDR` (default: `:8080`)
- `CONNECTOR_KAFKA_MAX_RETRIES` (default: `3`)
- `CONNECTOR_KAFKA_BACKOFF` (default: `500ms`)

## Exemplo rápido

```go
cfg, _ := config.FromEnv()
logger := logging.NewLogger(cfg.LogLevel)

publisher, _ := kafka.NewPublisher(kafka.Config{
    Brokers: cfg.KafkaBrokers,
    Topic: cfg.KafkaTopic,
    MaxRetries: cfg.KafkaMaxRetries,
    Backoff: cfg.KafkaBackoff,
})
defer publisher.Close()

server := httpserver.New(cfg.HTTPAddr)
_ = server.Start()

event := platformschemas.Event{Type: "device.status", Timestamp: time.Now().UnixMilli(), Payload: []byte(`{"status":"ok"}`)}
_ = publisher.Publish(context.Background(), kafka.MessageContext{
    TenantID: "tenant-1",
    SiteID: "site-1",
    DeviceID: "device-1",
    SchemaVersion: "v1",
}, event)

logger.Info("connector started")
```

## Template mínimo

```bash
cd template/minimal
go run .
```

## CI

Pipeline em `.github/workflows/ci.yml` executa:

- Em `pull_request`: validação de formatação (`gofmt -l`), `go test ./...` e `golangci-lint`.
- Em `push` para `main`: os mesmos checks e, se existir `Dockerfile`, publica a imagem `ghcr.io/<owner>/connector-sdk-builder` com as tags `:main` e `:sha`.
