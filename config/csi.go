package config

const CsiSockAddr = "/run/csi/csi.sock"

const CsiStagingTargetPath = "/mnt/staging"

const CsiMntPath = "/mnt/minik8s"

var CsiServerIP = ClusterMasterIP

const CsiStorageClassName = "nfs-csi"
