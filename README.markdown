demo
----
http://awwimage.herokuapp.com/random/pug

heroku deploy
-------------
* heroku config:add BUILDPACK_URL=https://github.com/kr/heroku-buildpack-go.git

* .godir
awwimage

* Procfile
web: awwimage

appengine
~/code/appengine/appcfg.py --oauth2 update .
