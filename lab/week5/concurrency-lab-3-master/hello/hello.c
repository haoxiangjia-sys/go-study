#include <pthread.h>
#include <stdio.h>

void *hello_world(void *args) {
    int *n = args;
    printf("Hello from thread %d\n", *n);
    pthread_exit(NULL);
}

int main(int argc, char const *argv[]) {
    pthread_t threads[5];   // 5 个线程句柄

    // 第一个循环：启动 5 个线程
    for (int i = 0; i < 5; i++) {
        int n = i + 1;
        pthread_create(&threads[i], NULL, hello_world, &n);
    }

    // 第二个循环：等 5 个线程都结束
    for (int i = 0; i < 5; i++) {
        pthread_join(threads[i], NULL);
    }
}
