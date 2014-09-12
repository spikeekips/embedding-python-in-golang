import sys
sys.stdout.write('`%s` imported\n' % __file__, )

from full import wsgi as wsgi_full

def application (*a, **kw) :
    try :
        return wsgi_full.get_wsgi_application()(*a, **kw)
    except :
        import traceback
        traceback.print_exc()
        return None


