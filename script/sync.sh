rsync -azP . root@k8s-2:/root/minik8s/ --exclude build --exclude .git --exclude log --exclude data
rsync -azP . root@k8s-3:/root/minik8s/ --exclude build --exclude .git --exclude log --exclude data
