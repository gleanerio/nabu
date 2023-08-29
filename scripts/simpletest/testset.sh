#!/bin/bash

NABU='go run ../../cmd/nabu/main.go'
CFG='../../iow/iow_dev_v3.yaml'
BLAZE='http://coreos.lan:9090/blazegraph/namespace/iow/sparql'
GRAPHDB='http://coreos.lan:7200/repositories/testing/statements'

echo "-------  BLAZEGRAPH"
echo "----------  clear, bulk, prune "
${NABU} clear --cfg ${CFG} --dangerous --endpoint blazegraph
${NABU} bulk --cfg ${CFG}  --prefix summoned/test --endpoint blazegraph
${NABU} prune --cfg ${CFG}  --prefix summoned/test --endpoint blazegraph

curl -H 'Accept: application/sparql-results+json' ${BLAZE} --data-urlencode 'query=SELECT (COUNT(DISTINCT ?graph) AS ?namedGraphsCount)(COUNT(*) AS ?triplesCount)WHERE {GRAPH ?graph {?subject ?predicate ?object}}' | jq '.results.bindings[0] | { namedGraphsCount: .namedGraphsCount.value, triplesCount: .triplesCount.value }' input.json -


echo "----------  clear, release, object, prune "
${NABU} clear --cfg ${CFG} --dangerous --endpoint blazegraph
${NABU} release --cfg ${CFG}  --prefix summoned/test --endpoint blazegraph
${NABU} object --cfg ${CFG}  graphs/latest/test_release.nq --endpoint blazegraph
${NABU} prune --cfg ${CFG}  --prefix summoned/test --endpoint blazegraph

curl -H 'Accept: application/sparql-results+json' ${BLAZE} --data-urlencode 'query=SELECT (COUNT(DISTINCT ?graph) AS ?namedGraphsCount)(COUNT(*) AS ?triplesCount)WHERE {GRAPH ?graph {?subject ?predicate ?object}}' | jq '.results.bindings[0] | { namedGraphsCount: .namedGraphsCount.value, triplesCount: .triplesCount.value }' input.json -


echo "----------  clear, prefix, prune "
${NABU} clear --cfg ${CFG} --dangerous --endpoint blazegraph
${NABU} prefix --cfg ${CFG}  --prefix summoned/test --endpoint blazegraph
${NABU} prune --cfg ${CFG}  --prefix summoned/test --endpoint blazegraph

curl -H 'Accept: application/sparql-results+json' ${BLAZE} --data-urlencode 'query=SELECT (COUNT(DISTINCT ?graph) AS ?namedGraphsCount)(COUNT(*) AS ?triplesCount)WHERE {GRAPH ?graph {?subject ?predicate ?object}}' | jq '.results.bindings[0] | { namedGraphsCount: .namedGraphsCount.value, triplesCount: .triplesCount.value }' input.json -


echo "-------  GraphDB"
echo "----------  clear, bulk, prune "
${NABU} clear --cfg ${CFG} --dangerous --endpoint graphdb
${NABU} bulk --cfg ${CFG}  --prefix summoned/test --endpoint graphdb
${NABU} prune --cfg ${CFG}  --prefix summoned/test --endpoint graphdb

curl -H 'Accept: application/sparql-results+json' ${GRAPHDB} --data-urlencode 'query=SELECT (COUNT(DISTINCT ?graph) AS ?namedGraphsCount)(COUNT(*) AS ?triplesCount)WHERE {GRAPH ?graph {?subject ?predicate ?object}}' | jq '.results.bindings[0] | { namedGraphsCount: .namedGraphsCount.value, triplesCount: .triplesCount.value }' input.json -


echo "----------  clear, release, object, prune "
${NABU} clear --cfg ${CFG} --dangerous --endpoint graphdb
${NABU} release --cfg ${CFG}  --prefix summoned/test --endpoint graphdb
${NABU} object --cfg ${CFG}  graphs/latest/test_release.nq --endpoint graphdb
${NABU} prune --cfg ${CFG}  --prefix summoned/test --endpoint graphdb

curl -H 'Accept: application/sparql-results+json' ${GRAPHDB} --data-urlencode 'query=SELECT (COUNT(DISTINCT ?graph) AS ?namedGraphsCount)(COUNT(*) AS ?triplesCount)WHERE {GRAPH ?graph {?subject ?predicate ?object}}' | jq '.results.bindings[0] | { namedGraphsCount: .namedGraphsCount.value, triplesCount: .triplesCount.value }' input.json -


echo "----------  clear, prefix, prune "
${NABU} clear --cfg ${CFG} --dangerous --endpoint graphdb
${NABU} prefix --cfg ${CFG}  --prefix summoned/test --endpoint graphdb
${NABU} prune --cfg ${CFG}  --prefix summoned/test --endpoint graphdb

curl -H 'Accept: application/sparql-results+json' ${GRAPHDB} --data-urlencode 'query=SELECT (COUNT(DISTINCT ?graph) AS ?namedGraphsCount)(COUNT(*) AS ?triplesCount)WHERE {GRAPH ?graph {?subject ?predicate ?object}}' | jq '.results.bindings[0] | { namedGraphsCount: .namedGraphsCount.value, triplesCount: .triplesCount.value }' input.json -
