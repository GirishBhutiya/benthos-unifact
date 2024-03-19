package opcua_plugin

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/gopcua/opcua/ua"
)

func randomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		randInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[randInt.Int64()]
	}
	return string(result)
}
func join(a, b string) string {
	if a == "" {
		return b
	}
	return a + "." + b
}

func ParseNodeIDs(incomingNodes []string) []*ua.NodeID {
	// Parse all nodeIDs to validate them.
	// loop through all nodeIDs, parse them and put them into a slice
	parsedNodeIDs := make([]*ua.NodeID, len(incomingNodes))

	for _, id := range incomingNodes {
		parsedNodeID, err := ua.ParseNodeID(id)
		if err != nil {
			return nil
		}

		parsedNodeIDs = append(parsedNodeIDs, parsedNodeID)
	}

	return parsedNodeIDs
}
func ParseTriggerNodeIDs(incomingTNodes []string) []*ua.NodeID {

	// Parse all nodeIDs to validate them.
	// loop through all nodeIDs, parse them and put them into a slice
	parsedTNodeIDs := make([]*ua.NodeID, len(incomingTNodes))

	for _, tNodeElements := range incomingTNodes {

		var tNodeObj map[string][]map[string]string
		err := json.Unmarshal([]byte(tNodeElements), &tNodeObj)
		if err != nil {
			return nil
		}

		for _, values := range tNodeObj {
			for _, obj := range values {

				fmt.Println(obj)

				nodeID := obj["node"]
				parsedTNodeID, err := ua.ParseNodeID(nodeID)
				if err != nil {
					return nil
				}
				parsedTNodeIDs = append(parsedTNodeIDs, parsedTNodeID)
			}
		}
	}
	return parsedTNodeIDs
}
