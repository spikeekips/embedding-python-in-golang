import json
import threading

def run (*a) :
    return json.dumps(
            [
                    threading.currentThread().ident,
                ] + list(a),
        )


