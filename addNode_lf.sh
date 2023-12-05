if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then

    echo "ERROR: Usage: ./addNode.sh [NODE_NAME] [NODE_ID] [TOLERANCE]"

    exit

fi



SERVER_NAME=$1

SERVER_ID=$2

TOLERANCE=$(($3+1))



DOCKER_TEMPLATE='

version: "3.4"



services:

  SERVER_NAME:

    build: ./ServerNode

    image: server_node

    entrypoint: ["./ServerNode", "PORT_NUM", "TOLERANCE", "SERVER_NAME", "entry_node:3000"]

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

COMPOSE_STR="${DOCKER_TEMPLATE//SERVER_NAME/${SERVER_NAME}}"

COMPOSE_STR="${COMPOSE_STR//TOLERANCE/${TOLERANCE}}"

COMPOSE_STR="${COMPOSE_STR//PORT_NUM/$((SERVER_ID + 4000))}"

echo "${COMPOSE_STR}">"./docker-compose.yml"



docker compose -f "./docker-compose.yml" -p chord-network up -d --build



# cleanup

rm "./docker-compose.yml"

