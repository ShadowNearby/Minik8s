package core

/* an basic example of a job apiobject:
apiVersion: v1
kind: Pod
metadata:
  name: gpu-job
  namespace: gpu
spec:
  containers:
    - name: gpu-server
      image: gpu-server
      command:
        - "./job.py"
      env:
        - name: source-path
          value: /gpu
        - name: job-name
          value: gpu-matrix
        - name: partition
          value: dgx2
        - name: "N"
          value: "1"
        - name: ntasks-per-node
          value: "1"
        - name: cpus-per-task
          value: "6"
        - name: gres
          value: gpu:1
      volumeMounts:
        - name: share-data
          mountPath: /gpu
  volumes:
    - name: share-data
      hostPath:
        path: /minik8s-sharedata/gpu/matrix
*/

type Job struct {
	ApiVersion string    `json:"apiVersion" yaml:"apiVersion"`
	MetaData   MetaData  `json:"metadata" yaml:"metadata"`
	Spec       JobSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status     PodStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type JobSpec struct {
	NodeSelector            map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	Containers              []Container       `json:"containers,omitempty" yaml:"containers,omitempty"`
	Volumes                 []Volume          `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	BackoffLimit            int               `json:"backoffLimit,omitempty" yaml:"backoffLimit,omitempty"`
	TtlSecondsAfterFinished int               `json:"ttlSecondsAfterFinished" yaml:"ttlSecondsAfterFinished"`
}
