typedef struct threadpool_t threadpool_t;
typedef struct threadpool_callback_t threadpool_callback_t;

threadpool_t *threadpool_create(int size);
pthread_t *start_thread (threadpool_t *threadpool, threadpool_callback_t *callback);
threadpool_callback_t *create_callback (threadpool_t *threadpool, void *function, void *args);
void threadpool_join (pthread_t *t);

