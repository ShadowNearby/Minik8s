#include "readfile.h"
#include <cuda_runtime.h>
#define BLOCK_SIZE 32
#define N 10000
#define CHECK
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
__global__ void matrixAddGlobalKernel(double * pfMatrixA, double * pfMatrixB, double * pfMatrixC, int w)
{
    int nRow = blockIdx.y * blockDim.y + threadIdx.y;
    int nCol = blockIdx.x * blockDim.x + threadIdx.x;
    if (nRow < w && nCol < w)
        pfMatrixC[nRow * w + nCol] = pfMatrixA[nRow * w + nCol] + pfMatrixB[nRow * w + nCol];
}
void matrixAddCPU( double * a, double * b, double * c )
{

  for( int row = 0; row < N; ++row )
    for( int col = 0; col < N; ++col )
    {
      c[row * N + col] = a[row*N+col]+b[row*N+col];
    }
}
int main(){
    generate_matrix_data();
    cudaError_t cudaStatus;
    unsigned long long size = (unsigned long long)N * N * sizeof (double );
    // Allocate input vectors h_A and h_B in host memory
    double * h_A = (double *)malloc(size);
    double * h_B = (double *)malloc(size);
    double * h_C = (double *)malloc(size);
    double * h_C_cpu = (double *)malloc(size);
    // Initialize input vectors
    read_values_from_file("matrix_a_data",h_A,size);
    read_values_from_file("matrix_b_data",h_B,size);
    // Allocate vectors in device memory
    double *d_A;
    cudaMalloc(&d_A, size);
    double *d_B;
    cudaMalloc(&d_B, size);
    double *d_C;
    cudaMalloc(&d_C, size);
    double gpuStartTime = clock(); // 记录开始时间
    // Copy vectors from host memory to device memory
    cudaMemcpy(d_A, h_A, size, cudaMemcpyHostToDevice);
    cudaMemcpy(d_B, h_B, size, cudaMemcpyHostToDevice);

    // Invoke kernel
    dim3 threadsPerBlock(BLOCK_SIZE, BLOCK_SIZE);
    dim3 blocksPerGrid((N + BLOCK_SIZE - 1) / BLOCK_SIZE, (N + BLOCK_SIZE - 1) / BLOCK_SIZE);
    double gpuMemcpyTime = clock(); // 记录开始时间
    matrixAddGlobalKernel<<<blocksPerGrid, threadsPerBlock>>>(d_A, d_B, d_C, N);
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
    matrixAddCPU( h_A, h_B, h_C_cpu);
    double cpuEndTime = clock();
    printf("GPU copy memory time: %lf\n", (gpuMemcpyTime - gpuStartTime) / CLOCKS_PER_SEC);
    printf("GPU computation time: %lf\n", (gpuEndTime - gpuMemcpyTime) / CLOCKS_PER_SEC);
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
    // Free host memory
    write_values_to_file("matrix_c_data",h_C,size);
    free(h_A);
    free(h_B);
    free(h_C);
    free(h_C_cpu);
}