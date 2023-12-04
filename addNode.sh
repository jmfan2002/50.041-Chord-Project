if [ -z "$1" ]; then
    echo "ERROR: Please pass server number, should be higher than current server count"
    exit
fi

SERVER_ID=$1
TOLERANCE=1

DOCKER_TEMPLATE='
version: "3.4"

services:
  server_nodeSERVER_ID:
    build: ./ServerNode
    image: server_node
    entrypoint: ["./ServerNode", "PORT_NUM", "TOLERANCE", "server_nodeSERVER_ID"]
    ports:
      - "PORT_NUM:PORT_NUM"
    deploy:
      mode: replicated
      replicas: 1
    networks:
      - chord-network

networks:
  chord-network:
'
COMPOSE_STR="${DOCKER_TEMPLATE//SERVER_ID/${SERVER_ID}}"
COMPOSE_STR="${COMPOSE_STR//TOLERANCE/${TOLERANCE}}"
COMPOSE_STR="${COMPOSE_STR//PORT_NUM/$((SERVER_ID + 4000))}"
echo "${COMPOSE_STR}">"./docker-compose.yml"

docker compose -f "./docker-compose.yml" -p chord-network up -d --build

# cleanup
rm "./docker-compose.yml"
