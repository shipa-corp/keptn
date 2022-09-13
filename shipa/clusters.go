package shipa

import "context"

// Cluster - represents Shipa cluster
type Cluster struct {
	Name      string            `json:"name"`
	Endpoint  *ClusterEndpoint  `json:"endpoint"`
	Resources *ClusterResources `json:"resources,omitempty"`
}

// ClusterEndpoint - part of Cluster object
type ClusterEndpoint struct {
	Addresses         []string `json:"addresses,omitempty"`
	Certificate       string   `json:"caCert,omitempty"`
	ClientCertificate string   `json:"clientCert,omitempty"`
	ClientKey         string   `json:"clientKey,omitempty"`
	Token             string   `json:"token,omitempty"`
	Username          string   `json:"username,omitempty"`
	Password          string   `json:"password,omitempty"`
}

// ClusterResources - part of Cluster object
type ClusterResources struct {
	Frameworks         []*Framework         `json:"frameworks,omitempty"`
	IngressControllers []*IngressController `json:"ingressControllers,omitempty"`
}

// IngressController - part of ClusterResources object
type IngressController struct {
	IngressIP     string `json:"ingressIp,omitempty"`
	ServiceType   string `json:"serviceType,omitempty"`
	Type          string `json:"type,omitempty"`
	HTTPPort      int64  `json:"httpPort,omitempty"`
	HTTPSPort     int64  `json:"httpsPort,omitempty"`
	ProtectedPort int64  `json:"protectedPort,omitempty"`
	Debug         bool   `json:"debug"`
	AcmeEmail     string `json:"acmeEmail,omitempty"`
	AcmeServer    string `json:"acmeServer,omitempty"`
}

// Framework - part of ClusterResources object
type Framework struct {
	Name string `json:"name,omitempty"`
}

// GetCluster - retrives cluster
func (c *Client) GetCluster(ctx context.Context, name string) (*Cluster, error) {
	cluster := &Cluster{}
	err := c.get(ctx, &cluster, apiClusters, name)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// CreateCluster - creates cluster
func (c *Client) CreateCluster(ctx context.Context, req *Cluster) error {
	return c.post(ctx, req, apiClusters)
}

// UpdateCluster - updates cluster
func (c *Client) UpdateCluster(ctx context.Context, req *Cluster) error {
	return c.put(ctx, req, apiClusters, req.Name)
}

// DeleteCluster - deletes cluster
func (c *Client) DeleteCluster(ctx context.Context, name string) error {
	return c.delete(ctx, apiClusters, name)
}
