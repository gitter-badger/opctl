package nodeprovider

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ NodeProvider

import "github.com/opspec-io/opctl/pkg/node"

type NodeProvider interface {
	ListNodes() (nodes []*node.InfoView, err error)
	CreateNode() (nodeInfo *node.InfoView, err error)
	KillNodeIfExists(nodeId string) (err error)
}
