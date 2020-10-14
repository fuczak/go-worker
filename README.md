## Building

`go build src/publisher/publisher.go`

`go build src/consumer/consumer.go`

## Running

Important! Currently, the RabbitMQ connection string is hard-coded and config values are ignored.

### Publisher

`publisher` params: 
```text
  -config string
        Path to the config file (default "./config.json")
  -data string
        Data of the message to publish (default "....")
  -queue string
        Name of the queue to publish the message to (default "go_test")
  -type string
        Type of the message to publish (default "Compute:Dots")
```

Example: `./publisher -queue go_test -type Compute:Dots -data "...."`

### Consumer

`consumer` params:
```text
  -config string
        Path to the config file (default "./config.json")
```

Example: `./consumer`