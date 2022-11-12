# glog2kafka

Exporting logs to kafka

```go
package main

import (
	"github.com/Shopify/sarama"
	log "github.com/ml444/glog"
	"github.com/ml444/glog/config"
	"github.com/ml444/glog2kafka"
)

func main() {
	cfg := sarama.NewConfig()
	endpoint, err := glog2kafka.NewKafkaEndpoint([]string{""}, "topic_name", cfg)
	if err != nil {
		return
	}
	err = log.InitLog(config.SetStreamer2Report(endpoint))
	if err != nil {
		return
	}
	defer log.Exit()

	log.Report("hello glog")
}
```

## Modified sarama.Config

```go
package main

import (
	"time"

	"github.com/Shopify/sarama"
	log "github.com/ml444/glog"
	"github.com/ml444/glog/config"
	"github.com/ml444/glog2kafka"
)

func GetKafkaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.ClientID = "glog2kafka"
	cfg.ChannelBufferSize = 1024

	//cfg.Producer.MaxMessageBytes = 1024 * 1024 * 5
	cfg.Producer.Return.Errors = true
	cfg.Producer.Flush.Bytes = 1024 * 1024 * 2
	cfg.Producer.Flush.Messages = 100
	cfg.Producer.Flush.MaxMessages = 1024
	cfg.Producer.Flush.Frequency = time.Millisecond * 100
	cfg.Producer.Retry.Max = 3
	return cfg
}

func main() {
	cfg := GetKafkaConfig()
	endpoint, err := glog2kafka.NewKafkaEndpoint([]string{""}, "topic_name", cfg)
	if err != nil {
		return
	}
	err = log.InitLog(config.SetStreamer2Report(endpoint))
	if err != nil {
		return
	}
	defer log.Exit()

	log.Report("hello glog")
}
```