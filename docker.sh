# !/bin/bash
xz -d rbq_anonymous_bot.xz -k
md5sum rbq_anonymous_bot
docker stop rbq_anonymous_bot_c
docker rm rbq_anonymous_bot_c
docker rmi rbq_anonymous_bot_i
docker build -t rbq_anonymous_bot_i .
docker run -it --name rbq_anonymous_bot_c --net work --ip 172.18.0.15 -d rbq_anonymous_bot_i
rm rbq_anonymous_bot
