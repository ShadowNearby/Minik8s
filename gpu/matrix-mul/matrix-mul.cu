#include "readfile.h"
#include <cuda_runtime.h>
#define BLOCK_SIZE 32
#define N 10000
// #define CHECK
#define RAND 200
#define INITIAL_SEED 12
void generate_matrix_data(){
    unsigned long long size = (unsigned long long)N*N*sizeof(double);
    double *a = (double*)malloc(size);
    double *b = (double*)malloc(size);
    srand(INITIAL_SEED);
    for( int row = 0; row < N; ++row ){
        for( int col = 0; col < N; ++col ){
            a[row*N + col] = (double)(rand() % RAND);
            b[row*N + col] = (double)(rand() % RAND);
        }
    }
    write_values_to_file("matrix_a_data",a,size);
    write_values_to_file("matrix_b_data",b,size);
    printf("generate data success\n");
}
__global__ void matrixMulGlobalKernel(double * pfMatrixA, double * pfMatrixB, double * pfMatrixC, int w)
{
    int nRow = blockIdx.y * blockDim.y + threadIdx.y;
    int nCol = blockIdx.x * blockDim.x + threadIdx.x;
    double fCVal = 0.0f;
    for(int i =0; i < w; i++)
    {
        fCVal += pfMatrixA[nRow * w + i] * pfMatrixB[i * w + nCol];
    }
    pfMatrixC[nRow * w + nCol] = fCVal;
}
void matrixMulCPU(double * A, double * B, double* C, int w) {
    for (int i = 0; i < w; ++i) {
        for (int j = 0; j < w; ++j) {
            double sum = 0.0f;
            for (int k = 0; k < w; ++k) {
                sum += A[i * w + k] * B[k * w + j];
            }
            C[i * w + j] = sum;
        }
    }
}

int main(){
    cudaError_t cudaStatus;
    unsigned long long size = (unsigned long long)N * N * sizeof (float );
    // Allocate input vectors h_A and h_B in host memory
    double* h_A = (double*)malloc(N * N * sizeof(double));
    double* h_B = (double*)malloc(N * N * sizeof(double));
    double* h_C = (double*)malloc(N * N * sizeof(double));
    double* h_C_cpu = (double*)malloc(N * N * sizeof(double));
    // Initialize input vectors
    read_values_from_file("matrix_a_data",h_A,size);
    read_values_from_file("matrix_b_data",h_B,size);
    // Allocate vectors in device memory

    double* d_A, * d_B, * d_C;
    cudaMalloc(&d_A, N * N * sizeof(double));
    cudaMalloc(&d_B, N * N * sizeof(double));
    cudaMalloc(&d_C, N * N * sizeof(double));

    cudaMemcpy(d_A, h_A, N * N * sizeof(double), cudaMemcpyHostToDevice);
    cudaMemcpy(d_B, h_B, N * N * sizeof(double), cudaMemcpyHostToDevice);

    // Invoke kernel
    // 定义线程块和网格大小
    dim3 threadsPerBlock(BLOCK_SIZE, BLOCK_SIZE);
    dim3 blocksPerGrid((N + BLOCK_SIZE - 1) / BLOCK_SIZE, (N + BLOCK_SIZE - 1) / BLOCK_SIZE);
    // 调用 GPU 核函数，使用 double 类型
    double gpuStartTime = clock();
    matrixMulGlobalKernel<<<blocksPerGrid, threadsPerBlock>>>(d_A, d_B, d_C, N);
    cudaStatus = cudaGetLastError();
    if (cudaStatus != cudaSuccess) {
        fprintf(stderr, "matrixAddGlobalKernel launch failed: %s\n", cudaGetErrorString(cudaStatus));
        return -1;
    }
    // Copy result from device memory to host memory
    // h_C contains the result in host memory
    cudaMemcpy(h_C, d_C, size, cudaMemcpyDeviceToHost);
    // Wait for the GPU to finish before proceeding
    cudaDeviceSynchronize();
    double gpuEndTime = clock();
#ifdef CHECK
    double cpuStartTime = clock();
    matrixMulCPU( h_A, h_B, h_C_cpu,N);
    double cpuEndTime = clock();
    printf("GPU computation time: %lf\n", (gpuEndTime - gpuStartTime) / CLOCKS_PER_SEC);
    printf("CPU computation time: %lf\n", (cpuEndTime - cpuStartTime) / CLOCKS_PER_SEC);
    for (int i = 0; i < N * N; i++) {
        if (fabs(h_C_cpu[i] - h_C[i]) > 1e-10) {
            fprintf(stderr, "CPU and GPU results differ at element %d!\n", i);
            exit(EXIT_FAILURE);
        }
    }
#endif
    // Free device memory
    cudaFree(d_A);
    cudaFree(d_B);
    cudaFree(d_C);
    free(h_A);
    free(h_B);
    free(h_C);
    free(h_C_cpu);
    return 0;
}