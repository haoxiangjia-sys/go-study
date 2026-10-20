#include <pthread.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

struct question {
    char *question;
    char *answer;
};
typedef struct question question_t;

struct ask_args {
    int score;
    question_t *question;
};
typedef struct ask_args ask_args_t;

// 1 把锁 + 1 个条件变量，保护下面两个标志
pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
pthread_cond_t cond = PTHREAD_COND_INITIALIZER;
bool time_up = false;   // 计时器：时间到了
bool answered = false;  // 答题线程：这题答完了

void *timer(void *args) {
    sleep(5);

    pthread_mutex_lock(&mutex);
    printf("\n5s have passed\n");
    time_up = true;
    pthread_cond_signal(&cond);  // 按门铃：叫醒主线程
    pthread_mutex_unlock(&mutex);

    pthread_exit(NULL);
}

void *ask(void *args) {
    ask_args_t *ask_args = args;

    printf("%s? ", ask_args->question->question);

    char input[128];
    fgets(input, 128, stdin);

    // strip newline
    input[strcspn(input, "\n")] = '\0';

    pthread_mutex_lock(&mutex);
    if (!time_up) {  // 时间到了才答的不算分
        if (!strcmp(input, ask_args->question->answer)) {
            printf("Correct!\n");
            ask_args->score++;
        } else {
            printf("Incorrect :-(\n");
        }
    }
    answered = true;
    pthread_cond_signal(&cond);  // 按门铃：叫醒主线程
    pthread_mutex_unlock(&mutex);

    pthread_exit(NULL);
}

int main(int argc, char const *argv[]) {
    question_t questions[] = {
        {.question = "3*2", .answer = "6"},
        {.question = "50/10", .answer = "5"},
        {.question = "2+1+1+1", .answer = "5"},
        {.question = "3^3", .answer = "27"},
        {.question = "3+3", .answer = "6"},
        {.question = "4/2", .answer = "2"}};

    pthread_t timer_thread;
    if (pthread_create(&timer_thread, NULL, timer, NULL)) {
        printf("Error creating timer thread.\n");
    }

    int score = 0;
    ask_args_t ask_args[6];  // 放在循环外，main 结束前一直有效（Q1 学到的）

    for (int i = 0; i < 6; ++i) {
        pthread_mutex_lock(&mutex);
        if (time_up) {  // 出题前先看看时间到了没
            pthread_mutex_unlock(&mutex);
            break;
        }
        answered = false;
        ask_args[i].score = score;
        ask_args[i].question = &questions[i];

        pthread_t ask_thread;
        if (pthread_create(&ask_thread, NULL, ask, &ask_args[i])) {
            printf("Error creating asker thread.\n");
        }
        pthread_detach(ask_thread);  // 不再 join 它

        // 睡觉，直到 "答完了" 或 "时间到了"
        while (!answered && !time_up) {
            pthread_cond_wait(&cond, &mutex);
        }

        score = ask_args[i].score;
        bool stop = time_up;
        pthread_mutex_unlock(&mutex);

        if (stop) {
            break;
        }
    }

    pthread_mutex_lock(&mutex);
    printf("End of questions, final score %d\n", score);
    pthread_mutex_unlock(&mutex);
    return 0;  // main 返回，整个进程结束，还卡在 fgets 的答题线程也一起结束
}
