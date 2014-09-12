import datetime
import json
import threading

from django.http import HttpResponse

from . import models as models_oldboy


threads_ids = list()

def home (request, ) :
    _tid = threading.currentThread().ident
    if request.META.get('X_FROM', ) in ('go', ) :
        assert _tid not in threads_ids
        threads_ids.append(_tid, )
    
    #print "> environ in django:\n", pprint.pprint(request.META, )

    _daesu = models_oldboy.OhDaeSu.objects.create()
    _daesu.eat(
            map(lambda x : models_oldboy.Mandu.objects.create(), range(15, ), ),
        )

    _response = HttpResponse(
            json.dumps(
                [
                        unicode(_daesu, ),
                        datetime.datetime.now().isoformat(),
                        _tid,
                        len(threads_ids, ),
                    ],
            ),
            content_type='application/json',
        )
    _response['X-Now'] = datetime.datetime.now().isoformat()

    return _response


