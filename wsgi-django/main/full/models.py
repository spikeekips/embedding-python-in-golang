# -*- coding: utf-8 -*-

from django.db import models


class Mandu (models.Model, ) :
    restaurant = models.TextField(default=u"자청룡", )


class OhDaeSu (models.Model, ) :
    message = models.TextField(default=u"너는 누구냐?", )

    mandu = models.ManyToManyField(Mandu, )

    def __unicode__ (self, ) :
        return u"매일 만두 %d개를 먹고 있다." % (
                self.mandu.count(),
            )

    def eat (self, mandus, ) :
        self.mandu = mandus
        self.save()

        return


