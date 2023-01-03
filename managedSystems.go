package hmc

type ManagedSystems struct {
	client *Client
	uuid   string
}

func (c *Client) ManagedSystems() *ManagedSystems {
	c.SetBaseURL(c.GetBaseURL() + "/ManagedSystem")
	return &ManagedSystems{client: c}
}

func (m *ManagedSystems) UUID(id string) *ManagedSystems {
	m.client.SetBaseURL(m.client.GetBaseURL() + "/" + id)
	return m
}

func (m *ManagedSystems) GET() (*Feed, *DetailedResponse, error) {
	return m.client.GET()
}
