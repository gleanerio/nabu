# Some SPARQL used

Count triples: result 11062

```SPARQL
prefix schema: <https://schema.org/> 
SELECT  (count(?s) as ?scount)
WHERE {   
  ?s ?p ?o 
}
```

Count all graphs: result 697

```SPARQL
prefix schema: <https://schema.org/> 
SELECT  (count(distinct ?g) as ?gcount)
WHERE {   
  graph ?g {?s ?p ?o} 
}
```

See if a graph exists

```SPARQL
ASK WHERE { GRAPH <graph> { ?s ?p ?o } }
```

Delete all

```SPARQL
DELETE  { ?s ?p ?o } WHERE { ?s ?p ?o }
DELETE {  GRAPH ?g {  } } WHERE { ?s ?p ?o  }
```