#! /bin/bash

cursh=`readlink -f $0`
curdir=`dirname $cursh`

cd $curdir/releasechk && go build && ./releasechk && rm -f ./releasechk && cd $curdir || cd $curdir