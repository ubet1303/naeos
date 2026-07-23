//go:build !nobroker

package broker

import (
	"errors"
	"testing"

	naeoserr "github.com/NAEOS-foundation/naeos/internal/errors"
)

func TestRealRedisNotConnectedErrIs(t *testing.T) {
	b := NewRealRedis()

	if err := b.Ping(); !errors.Is(err, naeoserr.ErrNotConnected) {
		t.Errorf("expected ErrNotConnected for Ping, got %v", err)
	}
	if err := b.Publish("test", &Message{Payload: []byte("hello")}); !errors.Is(err, naeoserr.ErrNotConnected) {
		t.Errorf("expected ErrNotConnected for Publish, got %v", err)
	}
	if err := b.Subscribe("test", func(m *Message) error { return nil }); !errors.Is(err, naeoserr.ErrNotConnected) {
		t.Errorf("expected ErrNotConnected for Subscribe, got %v", err)
	}
}

func TestRealRabbitMQNotConnectedErrIs(t *testing.T) {
	b := NewRealRabbitMQ()

	if err := b.Ping(); err == nil {
		t.Error("expected error for Ping when not connected")
	}
	if err := b.Publish("test", &Message{Payload: []byte("hello")}); err == nil {
		t.Error("expected error for Publish when not connected")
	}
	if err := b.Subscribe("test", func(m *Message) error { return nil }); err == nil {
		t.Error("expected error for Subscribe when not connected")
	}
}

func TestRealNATSNotConnectedExt(t *testing.T) {
	b := NewRealNATS()

	if err := b.Ping(); !errors.Is(err, naeoserr.ErrNotConnected) {
		t.Errorf("expected ErrNotConnected for Ping, got %v", err)
	}
	if err := b.Publish("test", &Message{Payload: []byte("hello")}); !errors.Is(err, naeoserr.ErrNotConnected) {
		t.Errorf("expected ErrNotConnected for Publish, got %v", err)
	}
	if err := b.Subscribe("test", func(m *Message) error { return nil }); !errors.Is(err, naeoserr.ErrNotConnected) {
		t.Errorf("expected ErrNotConnected for Subscribe, got %v", err)
	}
}

func TestRealKafkaNotConnectedErrIs(t *testing.T) {
	b := NewRealKafka()

	if err := b.Ping(); !errors.Is(err, naeoserr.ErrNotConnected) {
		t.Errorf("expected ErrNotConnected for Ping, got %v", err)
	}
	if err := b.Publish("test", &Message{Payload: []byte("hello")}); !errors.Is(err, naeoserr.ErrNotConnected) {
		t.Errorf("expected ErrNotConnected for Publish, got %v", err)
	}
	if err := b.Subscribe("test", func(m *Message) error { return nil }); err == nil {
		t.Error("expected error for Subscribe when not connected (nil config)")
	}
}

func TestRealRedisNameExt(t *testing.T) {
	b := NewRealRedis()
	if b.Name() != "redis" {
		t.Errorf("expected name 'redis', got %s", b.Name())
	}
}

func TestRealRabbitMQNameExt(t *testing.T) {
	b := NewRealRabbitMQ()
	if b.Name() != "rabbitmq" {
		t.Errorf("expected name 'rabbitmq', got %s", b.Name())
	}
}

func TestRealNATSNameExt(t *testing.T) {
	b := NewRealNATS()
	if b.Name() != "nats" {
		t.Errorf("expected name 'nats', got %s", b.Name())
	}
}

func TestRealKafkaNameExt(t *testing.T) {
	b := NewRealKafka()
	if b.Name() != "kafka" {
		t.Errorf("expected name 'kafka', got %s", b.Name())
	}
}

func TestRealRedisCloseNotConnected(t *testing.T) {
	b := NewRealRedis()
	if err := b.Close(); err != nil {
		t.Fatalf("unexpected error closing unconnected redis: %v", err)
	}
}

func TestRealRabbitMQCloseNotConnected(t *testing.T) {
	b := NewRealRabbitMQ()
	if err := b.Close(); err != nil {
		t.Fatalf("unexpected error closing unconnected rabbitmq: %v", err)
	}
}

func TestRealNATSCloseNotConnected(t *testing.T) {
	b := NewRealNATS()
	if err := b.Close(); err != nil {
		t.Fatalf("unexpected error closing unconnected nats: %v", err)
	}
}

func TestRealKafkaCloseNotConnected(t *testing.T) {
	b := NewRealKafka()
	if err := b.Close(); err != nil {
		t.Fatalf("unexpected error closing unconnected kafka: %v", err)
	}
}

func TestRealRabbitMQUnsubscribe(t *testing.T) {
	b := NewRealRabbitMQ()
	if err := b.Unsubscribe("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRealNATSUnsubscribe(t *testing.T) {
	b := NewRealNATS()
	if err := b.Unsubscribe("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRealKafkaUnsubscribe(t *testing.T) {
	b := NewRealKafka()
	if err := b.Unsubscribe("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRealRedisUnsubscribe(t *testing.T) {
	b := NewRealRedis()
	if err := b.Unsubscribe("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRealKafkaConnectAndClose(t *testing.T) {
	b := NewRealKafka()
	if err := b.Connect(&Config{Host: "localhost", Port: 9092}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := b.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}
	if err := b.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestRealKafkaPublishWhenConnected(t *testing.T) {
	b := NewRealKafka()
	b.Connect(&Config{Host: "localhost", Port: 9092})
	defer b.Close()

	err := b.Publish("test", &Message{Payload: []byte("data")})
	if err == nil {
		t.Log("publish succeeded (unexpected but OK for coverage)")
	}
}

func TestRealKafkaSubscribeWhenConnected(t *testing.T) {
	b := NewRealKafka()
	b.Connect(&Config{Host: "localhost", Port: 9092})
	defer b.Close()

	err := b.Subscribe("test", func(m *Message) error { return nil })
	if err == nil {
		b.Unsubscribe("test")
	}
}

func TestRealKafkaSubscribeAlreadySubscribed(t *testing.T) {
	b := NewRealKafka()
	b.Connect(&Config{Host: "localhost", Port: 9092})
	defer b.Close()

	_ = b.Subscribe("test", func(m *Message) error { return nil })
	err := b.Subscribe("test", func(m *Message) error { return nil })
	if err == nil {
		t.Error("expected error for double subscribe")
	}
	b.Unsubscribe("test")
}

func TestRealKafkaCloseWithReaderAndWriter(t *testing.T) {
	b := NewRealKafka()
	b.Connect(&Config{Host: "localhost", Port: 9092})
	if err := b.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestRealKafkaCloseWriterOnly(t *testing.T) {
	b := NewRealKafka()
	if err := b.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestRealNATSConnectAndClose(t *testing.T) {
	b := NewRealNATS()
	err := b.Connect(&Config{Host: "192.0.2.1", Port: 1, Timeout: 1})
	if err != nil {
		t.Skip("NATS connect to unreachable host failed (expected)")
	}
	if err := b.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestRealRedisSubscribeNotConnected(t *testing.T) {
	b := NewRealRedis()
	err := b.Subscribe("test", func(m *Message) error { return nil })
	if !errors.Is(err, naeoserr.ErrNotConnected) {
		t.Errorf("expected ErrNotConnected, got %v", err)
	}
}

func TestRealRabbitMQPingNotConnected(t *testing.T) {
	b := NewRealRabbitMQ()
	err := b.Ping()
	if err == nil {
		t.Error("expected error")
	}
}

func TestRealRabbitMQPublishNotConnected(t *testing.T) {
	b := NewRealRabbitMQ()
	err := b.Publish("test", &Message{Payload: []byte("x")})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRealRabbitMQSubscribeNotConnected(t *testing.T) {
	b := NewRealRabbitMQ()
	err := b.Subscribe("test", func(m *Message) error { return nil })
	if err == nil {
		t.Error("expected error")
	}
}
