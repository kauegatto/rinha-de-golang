#!/bin/sh

# Use este script para executar testes locais

RESULTS_WORKSPACE="$(pwd)/load-test/user-files/results"
GATLING_WORKSPACE="$(pwd)/load-test/user-files"
GATLING_DIR=$HOME/gatling/gatling-charts-highcharts-bundle-3.12.0
    

runGatling() {
     $GATLING_DIR/mvnw gatling:test \
        -DsimulationsFolder="$GATLING_WORKSPACE/simulations"
}

startTest() {
    for i in {1..20}; do
        # 2 requests to wake the 2 api instances up :)
        curl --fail http://localhost:9999/clientes/1/extrato && \
        echo "" && \
        curl --fail http://localhost:9999/clientes/1/extrato && \
        echo "" && \
        runGatling && \
        break || sleep 2;
    done
}

startTest