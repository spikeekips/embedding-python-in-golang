import sys
sys.stdout.write('`%s` imported\n' % __file__, )

import datetime
import random
import json
import threading

threads_ids = list()

def application (environ, start_response, ) :
    assert callable(start_response, )
    assert type(environ, ) is dict

    _tid = threading.currentThread().ident
    assert _tid not in threads_ids
    threads_ids.append(_tid, )

    # start_response
    _response_headers = [
            ('Content-Type', 'application/json'),
            ('X-py', datetime.datetime.now().isoformat(), ),
        ]

    _status = '%d OK' % random.choice(range(500))
    start_response(_status, _response_headers, )

    return json.dumps(
            [
                    environ,
                    datetime.datetime.now().isoformat(),
                    _tid,
                    len(threads_ids, ),
                ],
        )


