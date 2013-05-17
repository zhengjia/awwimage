heroku deploy
1. heroku config:add BUILDPACK_URL=https://github.com/kr/heroku-buildpack-go.git

2. .godir
awwimage

3. Procfile
web: awwimage

appengine
~/code/appengine/appcfg.py --oauth2 update .
