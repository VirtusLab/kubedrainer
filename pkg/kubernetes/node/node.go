package node

import (
	"strings"

	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/VirtusLab/go-extended/pkg/matcher"
	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Node represents a Kubernetes node API
type Node struct {
	Client kubernetes.Interface
}

// GetNode returns error if the given node cannot be found
func (n *Node) GetNode(nodeName string) (v1.Node, error) {
	if len(nodeName) == 0 {
		return v1.Node{}, errors.New("node name cannot be empty")
	}

	nodes, err := n.Client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return v1.Node{}, errors.Wrap(err)
	}
	if len(nodes.Items) == 0 {
		return v1.Node{}, errors.New("no nodes found")
	}
	for _, n := range nodes.Items {
		if n.Name == nodeName {
			glog.V(1).Infof("Found node: '%s'", nodeName)
			return n, nil
		}
	}

	// inform about available nodes
	var nodeNames []string
	for _, n := range nodes.Items {
		nodeNames = append(nodeNames, n.Name)
	}
	return v1.Node{}, errors.Errorf("node '%s' not found, got: '%s'", nodeName, strings.Join(nodeNames, ", "))
}

// GetProviderID returns a cloud provider specific ID for the given Kubernetes node
func (n *Node) GetProviderID(nodeName string) (string, string, error) {
	node, err := n.Client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return "", "", errors.Wrap(err)
	}

	providerIDExpression := `^(?P<ProviderName>\S+)://(?P<ProviderSpecificNodeID>\S+)`
	results, ok := matcher.Must(providerIDExpression).MatchGroups(node.Spec.ProviderID)
	if !ok {
		return "", "", errors.Errorf("Can't match expression '%s' to '%s'",
			providerIDExpression, node.Spec.ProviderID)
	}
	providerName, ok := results["ProviderName"]
	if !ok {
		return "", "", errors.Errorf("Missing 'ProviderName' when expression '%s' was applied to '%s'",
			providerIDExpression, node.Spec.ProviderID)
	}
	providerSpecificNodeID, ok := results["ProviderSpecificNodeID"]
	if !ok {
		return "", "", errors.Errorf("Missing 'ProviderSpecificNodeID' when expression '%s' was applied to '%s'",
			providerIDExpression, node.Spec.ProviderID)
	}
	return providerName, providerSpecificNodeID, nil
}
