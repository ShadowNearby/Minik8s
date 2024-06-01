package serverless

import (
	"encoding/json"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/serverless/activator"
	"minik8s/utils"
	"testing"
)

const ImagePath = "shadownearby"

func generateImage(name string) string {
	return fmt.Sprintf("%s/%s:v1", ImagePath, name)
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
	err = utils.CreateObject(core.ObjReplicaSet, replica.MetaData.Namespace, replica)
	if err != nil {
		t.Errorf("GenerateReplicaSet failed, error creating object: %s", err)
	}
	t.Logf("ReplicaSet created successfully，%s", replica.MetaData.Namespace)
}
