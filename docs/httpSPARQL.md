# SPARQL over HTTP

## About

looking at different patterns for loading large files over
http for SPARQL.

### GraphDB POST

```bash
curl -X POST -H 'Content-Type:application/n-quads'  --data-binary @May4Buildings.nq  http://192.168.86.45:32774/repositories/loadtest/rdf-graphs/services
```

### Blazegraph

```bash
curl -X POST -H 'Content-Type:text/x-nquads'  --data-binary @May4Buildings.nq  http://192.168.86.45:32772/blazegraph/namespace/loadtest/sparql
```

## Test SPARQL

```sparql
SELECT (COUNT(?s) AS ?triples) 
WHERE { 
  ?s ?p ?o 
}
```

```sparql
SELECT *

WHERE { 
  graph ?g {
  ?s ?p ?o 
     }
}
```