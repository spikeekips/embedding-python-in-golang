import json
import threading
import datetime

def run (*a) :
    return json.dumps(
            [
                    threading.currentThread().ident,
                    datetime.datetime.now().isoformat(),
                ] + list(a),
        )


