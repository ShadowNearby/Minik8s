#!/bin/bash

#SBATCH --job-name=add-job
#SBATCH --partition=dgx2
#SBATCH -N 1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=6
#SBATCH --gres=gpu:1
#SBATCH --partition=dgx2
#SBATCH --output=output.txt
#SBATCH --error=error.txt
#SBATCH --mail-type=end
#SBATCH --mail-user=starcloud.0416@sjtu.edu.cn
ulimit -s unlimited
ulimit -l unlimited
module load cuda/12 gcc
make build
make run