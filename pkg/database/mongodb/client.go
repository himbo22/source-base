package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/himbo22/source-base/pkg/settings"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	defaultMaxPoolSize     = uint64(100)
	defaultMinPoolSize     = uint64(5)
	defaultMaxConnIdleTime = uint64(300) // seconds
	defaultTimeout         = 10          // seconds
)

type Client struct {
	Client *mongo.Client
	Config *settings.MongoDB
}

func NewClient(cfg *settings.MongoDB) *Client {
	return &Client{Config: cfg}
}

func (c *Client) Connect() error {
	c.setDefaultConfig()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Config.Timeout)*time.Second)
	defer cancel()

	uri := c.buildURI()
	clientOpts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(c.Config.MaxPoolSize).
		SetMinPoolSize(c.Config.MinPoolSize).
		SetMaxConnIdleTime(time.Duration(c.Config.MaxConnIdleTime) * time.Second).
		SetServerSelectionTimeout(time.Duration(c.Config.Timeout) * time.Second).
		SetConnectTimeout(time.Duration(c.Config.Timeout) * time.Second)

	if c.Config.TLSEnabled {
		tlsConfig, err := options.BuildTLSConfig(map[string]any{
			"insecure": c.Config.TLSInsecure,
		})
		if err != nil {
			return fmt.Errorf("failed to build tls config: %w", err)
		}
		clientOpts.SetTLSConfig(tlsConfig)
	}

	if c.Config.AuthSource != "" {
		clientOpts.SetAuth(options.Credential{
			Username:   c.Config.Username,
			Password:   c.Config.Password,
			AuthSource: c.Config.AuthSource,
		})
	} else if c.Config.Username != "" && c.Config.Password != "" {
		clientOpts.SetAuth(options.Credential{
			Username: c.Config.Username,
			Password: c.Config.Password,
		})
	}

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create mongo client: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return fmt.Errorf("failed to ping mongo: %w", err)
	}

	c.Client = client
	return nil
}

func (c *Client) buildURI() string {
	if c.Config.URL != "" {
		return c.Config.URL
	}

	host := c.Config.Host
	if host == "" {
		host = "localhost"
	}
	port := c.Config.Port
	if port == 0 {
		port = 27017
	}

	return fmt.Sprintf("mongodb://%s:%d", host, port)
}

func (c *Client) Database(name string) *mongo.Database {
	if c.Client == nil {
		return nil
	}
	dbName := name
	if dbName == "" && c.Config.Database != "" {
		dbName = c.Config.Database
	}
	return c.Client.Database(dbName)
}

func (c *Client) Collection(dbName, collName string) *mongo.Collection {
	db := c.Database(dbName)
	if db == nil {
		return nil
	}
	return db.Collection(collName)
}

func (c *Client) Ping(ctx context.Context) error {
	if c.Client == nil {
		return fmt.Errorf("mongo client not initialized")
	}
	return c.Client.Ping(ctx, readpref.Primary())
}

func (c *Client) Disconnect(ctx context.Context) error {
	if c.Client == nil {
		return nil
	}
	return c.Client.Disconnect(ctx)
}

func (c *Client) StartSession(ctx context.Context) (*mongo.Session, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("mongo client not initialized")
	}
	return c.Client.StartSession()
}

func (c *Client) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	session, err := c.StartSession(ctx)
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, fn)
}

func (c *Client) setDefaultConfig() {
	if c.Config.MaxPoolSize == 0 {
		c.Config.MaxPoolSize = defaultMaxPoolSize
	}
	if c.Config.MinPoolSize == 0 {
		c.Config.MinPoolSize = defaultMinPoolSize
	}
	if c.Config.MaxConnIdleTime == 0 {
		c.Config.MaxConnIdleTime = defaultMaxConnIdleTime
	}
	if c.Config.Timeout == 0 {
		c.Config.Timeout = defaultTimeout
	}
}