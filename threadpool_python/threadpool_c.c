// #cgo CFLAGS: -v
// +build darwin linux

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include <pthread.h>

#include "threadpool_c.h"

typedef enum {
    immediate_shutdown = 1,
    graceful_shutdown  = 2
} threadpool_shutdown_t;

typedef enum {
    threadpool_graceful       = 1
} threadpool_destroy_flags_t;

typedef enum {
    threadpool_invalid        = -1,
    threadpool_lock_failure   = -2,
    threadpool_queue_full     = -3,
    threadpool_shutdown       = -4,
    threadpool_thread_callback_failure = -5
} threadpool_error_t;

struct threadpool_t {
  int size;
  int shutdown;
  int count;

  // queue
  int front;
  int queue_count;
  pthread_t *threads;

  pthread_mutex_t lock;
  pthread_cond_t notify;
  pthread_mutex_t queue_lock;
  pthread_cond_t queue_notify;
};

struct threadpool_callback_t {
    void (*function)(void *);
    void *args;
	threadpool_t *threadpool;
};

int threadpool_free(threadpool_t *threadpool) {
    if(threadpool == NULL || threadpool->count > 0) {
        return -1;
    }

    free(threadpool->threads);

    pthread_mutex_lock(&(threadpool->lock));
    pthread_mutex_destroy(&(threadpool->lock));
    pthread_cond_destroy(&(threadpool->notify));

	pthread_mutex_lock(&(threadpool->queue_lock));
	pthread_mutex_destroy(&(threadpool->queue_lock));
	pthread_cond_destroy(&(threadpool->queue_notify));

    free(threadpool);
    return 0;
}

int threadpool_destroy(threadpool_t *threadpool, int flags)
{
    int i, err = 0;

    if(threadpool == NULL) {
        return threadpool_invalid;
    }

    if(pthread_mutex_lock(&(threadpool->lock)) != 0) {
        return threadpool_lock_failure;
    }

    do {
        if(threadpool->shutdown) {
            err = threadpool_shutdown;
            break;
        }

        threadpool->shutdown = (flags & threadpool_graceful) ?
            graceful_shutdown : immediate_shutdown;

        if((pthread_cond_broadcast(&(threadpool->notify)) != 0) ||
           (pthread_mutex_unlock(&(threadpool->lock)) != 0)) {
            err = threadpool_lock_failure;
            break;
        }

        for(i = 0; i < threadpool->size; i++) {
            if(pthread_join(threadpool->threads[i], NULL) != 0) {
                err = threadpool_thread_callback_failure;
            }
        }
    } while(0);

    if(!err) {
        threadpool_free(threadpool);
    }
    return err;
}

threadpool_t *threadpool_create(int size)
{
    threadpool_t *threadpool;
    int i;

    if((threadpool = (threadpool_t *)malloc(sizeof(threadpool_t))) == NULL) {
        goto err;
    }

    threadpool->count = 0;
    threadpool->size = size;
    threadpool->shutdown = 0;

    threadpool->threads = (pthread_t *)malloc(sizeof(pthread_t) * threadpool->size);

    threadpool->queue_count = threadpool->size;
    threadpool->front = 0;

    if((pthread_mutex_init(&(threadpool->lock), NULL) != 0) || (pthread_cond_init(&(threadpool->notify), NULL) != 0)) {
        goto err;
    }

    if((pthread_mutex_init(&(threadpool->queue_lock), NULL) != 0) || (pthread_cond_init(&(threadpool->queue_notify), NULL) != 0)) {
        goto err;
    }

    return threadpool;

 err:
    if(threadpool) {
        threadpool_free(threadpool);
    }
    return NULL;
}

void threadpool_join (pthread_t *t) {
    pthread_join(*t, 0);
    //pthread_detach(*t);

    return;
}

pthread_t *threadpool_get_thread_id(threadpool_t *threadpool) {
    pthread_mutex_lock(&(threadpool->queue_lock));

    if (threadpool->queue_count == 0) {
		goto err;
    }

    threadpool->queue_count--;
    threadpool->front++;
    if (threadpool->front == threadpool->size) {
        threadpool->front=0;
    }

    pthread_t *t = &(threadpool->threads[threadpool->front]);

    pthread_mutex_unlock(&(threadpool->queue_lock));
	return t;
err:
    pthread_mutex_unlock(&(threadpool->queue_lock));
	return NULL;
}

void threadpool_return_thread_id(threadpool_t *threadpool) {
    pthread_mutex_lock(&(threadpool->queue_lock));

    threadpool->queue_count++;

    pthread_mutex_unlock(&(threadpool->queue_lock));
    return;
err:
    pthread_mutex_unlock(&(threadpool->queue_lock));
	return;
}

static void *threadpool_thread_callback(void *args)
{
	threadpool_callback_t *callback = (threadpool_callback_t *)args;

    //fprintf(stderr, "99999999: %d\n", pthread_self());

    runGoCallback(callback->function, callback->args);

	callback->threadpool->count -= 1;

    threadpool_return_thread_id(callback->threadpool);

    pthread_mutex_lock(&(callback->threadpool->queue_lock));
    pthread_cond_signal(&(callback->threadpool->queue_notify));
    pthread_mutex_unlock(&(callback->threadpool->queue_lock));

	free(callback);
    pthread_exit(NULL);
}

pthread_t *start_thread (threadpool_t *threadpool, threadpool_callback_t *callback) {
    // if threadpool is shutdown, cancel.
    if ((threadpool->shutdown == immediate_shutdown) ||
			(threadpool->shutdown == graceful_shutdown)) {
        goto err;
    }

    // if the threadpool size reached, wait.
    pthread_t *t;

    while(1) {
		t = threadpool_get_thread_id(threadpool);
		if (t != NULL)
		{
			break;
		}

		pthread_mutex_lock(&(threadpool->queue_lock));
        pthread_cond_wait(&(threadpool->queue_notify), &(threadpool->queue_lock));
		pthread_mutex_unlock(&(threadpool->queue_lock));
    }

    pthread_mutex_lock(&(threadpool->lock));
	threadpool->count += 1;
    pthread_mutex_unlock(&(threadpool->lock));

    int _r = pthread_create(
        t,
        NULL,
        threadpool_thread_callback,
		callback
    );

    if(_r != 0) {
        threadpool_destroy(threadpool, 0);
        return NULL;
    }

    //threadpool_join(t);
    return t;

err:
	if(threadpool) {
	    threadpool_free(threadpool);
	}
	return NULL;
}


threadpool_callback_t *create_callback (threadpool_t *threadpool, void *function, void *args) {
    threadpool_callback_t *c;
    
    c = (threadpool_callback_t *)malloc(sizeof(threadpool_callback_t));
    c->function = function;
    c->args = args;
    c->threadpool = threadpool;

    return c;
}
