1.下载redis镜像<br>
docker pull redis<br>
2.下载mysql镜像<br>
docker pull mysql<br>

3.运行redis&mysql容器<br>
docker run --name redis16379 -p 16739:6379  -d -v /Users/mengfanzhen/docker/redis:/data  redis:latest --requirepass  "abc@0912" <br>
docker run --name mysql13306 -p 13306:3306 -e MYSQL_ROOT_PASSWORD=abc@198920 -v /Users/mengfanzhen/docker/mysql:/var/lib/mysql -d mysql:latest<br>

4.应用镜像构建<br>
 docker build . -t bluebell<br>
 
5.启动并关联相关容器<br>
docker run --name bluebell -d -p 19000:9000 --link=redis16379:redis16379 --link=mysql13306:mysql13306 bluebell:latest

