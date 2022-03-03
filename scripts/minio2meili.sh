#!/bin/bash
# A wrapper script for loading RDF into Jena from Minio
# usage:  load2jena.sh souceBucket targetDataBase targetGraph
# example:  load2jena.sh local/gleaner-milled/run1 index common 
#Notes:
# mc cat nas/issue31/summoned/obis/ffc4657bd5ef66bcae561f9e63482fa787bab86a.jsonld
#| jsonld flatten 
#| curl  -H 'Content-Type: application/json' -X POST -d @- http://127.0.0.1:7700/indexes/movies/documents

# mc cat nas/issue31/summoned/obis/ffc4657bd5ef66bcae561f9e63482fa787bab86a.jsonld |  jq '. |= {"id": "1234"} + .' | curl  -H 'Content-Type: application/json' -X POST -d @- http://127.0.0.1:7700/indexes/movies/documents


mc_cmd() {
        mc ls $1 | awk '{print $6}'
}

# If you use this for ntriples, be sure to add in a graph in the URL target
for i in $(mc_cmd $1); do
    echo $i
    mc cat $1/$i | curl  -H 'Content-Type: application/json' -X POST -d @- http://127.0.0.1:7700/indexes/movies/documents 
done

