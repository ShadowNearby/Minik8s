#!/usr/bin/python3
import paramiko
from time import sleep
from scp import SCPClient
from os import getenv


NREAD = 100000
hostname = 'pilogin.hpc.sjtu.edu.cn'  # 远程服务器地址
username = 'stu091'  # SSH用户名

ssh = paramiko.SSHClient()
ssh.load_system_host_keys()
ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
ssh.connect(hostname,username=username)


job_submit_tag = "Submitted batch job"
line_finish_tag = "[stu091@"
PENDING = "PENDING"
COMPLETED = "COMPLETED"
FAILED = "FAILED"

source_path = getenv("source-path")
job_name = getenv("job-name")
partition= getenv("partition")
N = getenv("N")
ntasks_per_node = getenv("ntasks-per-node")
cpus_per_task = getenv("cpus-per-task")
gres = getenv("gres")
if not source_path:
    source_path = "/matrix-mul"
if not job_name:
    job_name = "matrix-mul"
if not job_name or not source_path:
    print("env error")
    exit(0)
if source_path[-1] == "/":
    source_path = source_path[:-1]

if not partition:
    partition = "dgx2"
if not N:
    N = 1
if not ntasks_per_node:
    ntasks_per_node = 1
if not cpus_per_task:
    cpus_per_task = 6
if not gres:
    gres = "gpu:1"

def generate_slurm():
    print("=>\tgenerating slurm")
    with open(f"./{job_name}.slurm","w") as f:
        f.write("#!/bin/bash\n")
        f.write(f"#SBATCH --job-name={job_name}\n")
        f.write(f"#SBATCH --partition={partition}\n")
        f.write(f"#SBATCH -N {N}\n")
        f.write(f"#SBATCH --ntasks-per-node={ntasks_per_node}\n")
        f.write(f"#SBATCH --cpus-per-task={cpus_per_task}\n")
        f.write(f"#SBATCH --gres={gres}\n")
        # result must exist . is the same dir as .slurm
        f.write(f"#SBATCH --output=output.txt\n")
        f.write(f"#SBATCH --error=error.txt\n")
        f.write(f"ulimit -s unlimited\n")
        f.write(f"ulimit -l unlimited\n")
        f.write("module load gcc cuda/12 \n")
        f.write("make build\n")
        f.write("make run\n")

def upload_source():
    print("=>\tuploading source")
    scp = SCPClient(ssh.get_transport(),socket_timeout=16)
    scp.put(source_path,recursive=True,remote_path=f"~/")
    scp.put(f"./{job_name}.slurm",f"~/{job_name}.slurm")
    scp.close()

def download_result(job_id):
    print("=>\tdownloading result")
    scp = SCPClient(ssh.get_transport(),socket_timeout=16)
    scp.get(f"~/{job_name}/output.txt",recursive=True,local_path=f"{source_path}/")
    scp.get(f"~/{job_name}/error.txt",recursive=True,local_path=f"{source_path}/")
    scp.close()

def submit_job():
    t = 3
    while t:
        s = ssh.invoke_shell()
        print("=>\tstarting ssh")
        sleep(2)
        recv = s.recv(NREAD).decode('utf-8')
        if recv.find("stu091") == -1:
            print("start ssh failed,retrying")
            t -= 1
            sleep(5)
            continue

        print("=>\tstart ssh success")
        print("=>\tsending sbatch")
        s.send(f"cd ~/{job_name} && sbatch ./{job_name}.slurm\n")
        sleep(5)

        recv = s.recv(NREAD).decode('utf-8')
        index = recv.find(job_submit_tag)
        if index ==-1:
            print(recv)
            print("sbatch failed,retrying")
            t -= 1
            sleep(5)
            continue
        print("=>\tsbatch success")
        job_id = recv[index+len(job_submit_tag)+1:recv.index(line_finish_tag)-2]
        print(f"{job_id=}")
        print("start checking job status")
        check_status_cmd = f"sacct | grep {job_id} | awk '{{print $6}}'"

        while True:
            s.send(check_status_cmd+"\n")
            sleep(2)
            recv = s.recv(NREAD).decode('utf-8')
            status = recv[recv.index(check_status_cmd)+len(check_status_cmd)+2:recv.index(line_finish_tag)-2]
            print(f"{status=}")
            if status.find(FAILED)!=-1:
                print("job failed")
                return job_id
            if status.find(COMPLETED)==-1:
                sleep(10)
            else:
                return job_id

generate_slurm()
upload_source()
job_id = submit_job()
if job_id:
    download_result(job_id)
print("finish")
