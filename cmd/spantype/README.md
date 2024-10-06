```
$ gcloud spanner databases execute-sql ${SPANNER_DATABASE} \
    --format="json" --query-mode=PLAN \
    --sql 'SELECT 1 AS n, ARRAY(SELECT AS STRUCT 1 AS n, ARRAY(SELECT AS STRUCT 1 AS n) AS `inner`) AS `outer`' \
    | jq .metadata.rowType | ./spantype 
n INT64, outer ARRAY<STRUCT<n INT64, inner ARRAY<STRUCT<n INT64>>>>
```