docker run -d --name rabbitmq_sayur \
-p 5672:5672 -p 15672:15672 rabbitmq:3-management-alpine

docker run -d --name elastic_sayur \
-e "discovery.type=single-node"
-p 9200:9200 \
-p 9300:9300 docker.elastic.co/elasticsearch/elasticsearch:7.17.28

docker run -d --name redis_sayur -p 6379:6379 redis:latest

migrate -database "postgres://postgres:sayur@localhost:5434/sayur_user_service?sslmode=disable" -path database/migrations up
