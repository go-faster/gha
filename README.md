# gha

github archive utilities

```json
[
  {
    "input": "1693 GB",
    "content": "13 TB",
    "output": "1191 GB"
  }
]
```

```json
[
  {
    "state": "NotFound",
    "count": 319
  },
  {
    "state": "Ready",
    "count": 68952
  }
]
```

## Results

## Missing chunks
319 of 68952 chunks are missing, not sure about restore, not critical.

## No language data
This should be retrieved separately.

```sql
SELECT * FROM `bigquery-public-data.github_repos.languages`;
```

Can be exported to json, loaded into database and joined with events.


### Source

Programming languages by repository as reported by 
GitHub's https://developer.github.com/v3/repos/#list-languages API

### Properties
* No repo id, just name
* Probably no removed or renamed repos
* ~3 million entries
* Language data is in array (language name, bytes)