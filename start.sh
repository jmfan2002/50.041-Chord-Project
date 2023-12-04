if [ -z "$1" ] || [ -z "$2" ]; then
    echo "ERROR: Usage: ./start.sh [NUM_SERVERS] [TOLERANCE]"
    exit
fi

TOLERANCE=$2

# set up docker compose string
NUM_SERVERS=$1
SERVER_NODE_STR='
  server_nodeSERVER_NUM:
    build: ./ServerNode
    image: server_node
    entrypoint: ["/bin/sh", "-c", "sleep SERVER_NUM && ./ServerNode PORT_NUM TOLERANCE server_nodeSERVER_NUM entry_node:3000"]
    ports:
      - "PORT_NUM:PORT_NUM"
    deploy:
      mode: replicated
      replicas: 1
    networks:
      - chord-network
'
COMPOSE_STR='
version: "3.4"

services:
  entry_node:
    build: ./EntryNode
    image: entry_node
    ports:
      - "3000:3000"
    entrypoint: ["./EntryNode", "-port=3000", "-k=TOLERANCE"]
    networks:
      - chord-network
'
END_STR='
networks:
  chord-network:
'
for i in $(seq 1 $NUM_SERVERS); do
    STR="${SERVER_NODE_STR//SERVER_NUM/${i}}"
    COMPOSE_STR+="${STR//PORT_NUM/$((i+4000))}"
done
COMPOSE_STR+="${END_STR}"
COMPOSE_STR="${COMPOSE_STR//TOLERANCE/${TOLERANCE}}"
echo "${COMPOSE_STR}">./docker-compose.yml

# optional: remove -d --build to keep docker compose active in the terminal,
#       to show live container logs in the terminal
docker compose -f "./docker-compose.yml" -p chord-network up -d --build

# cleanup
rm "./docker-compose.yml"
