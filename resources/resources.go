package resources

type Resource interface {
	Import() (Resource, error)
}

type AWSResourceId struct {
	Id *string
}

// c *client
