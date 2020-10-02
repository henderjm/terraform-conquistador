package resources

type Resource interface {
	Import(c *client) error
}

type AWSResourceId struct {
	Id *string
}

// c *client
