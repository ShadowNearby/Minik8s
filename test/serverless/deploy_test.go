package serverless

import (
	"encoding/json"
	"minik8s/config"
	"minik8s/pkgs/serverless/activator"
	"testing"
)

func generateImage(name string) string {
	return config.LocalServerIp + ":5000/" + name + ":latest"
}

func TestGenerateReplicaSet(t *testing.T) {
	replica := activator.GenerateReplicaSet("test", "serverless", generateImage("test"), 0)
	if replica.MetaData.Name != "test" {
		t.Errorf("GenerateReplicaSet failed, expected %s, got %s", "test", replica.MetaData.Name)
	}
	if replica.MetaData.Namespace != "serverless" {
		t.Errorf("GenerateReplicaSet failed, expected %s, got %s", "serverless", replica.MetaData.Namespace)
	}
	if replica.Spec.Replicas != 0 {
		t.Errorf("GenerateReplicaSet failed, expected %d, got %d", 0, replica.Spec.Replicas)
	}
	// print the replicaSet
	replicaJson, err := json.MarshalIndent(replica, "", "    ")
	if err != nil {
		t.Errorf("GenerateReplicaSet failed, error marshalling replicas: %s", err)
	}

	t.Logf("replicaSet: %s", replicaJson)
}
