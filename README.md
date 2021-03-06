# gha

Research project of go-faster.

> Content based on
> <a href="https://www.gharchive.org/">www.gharchive.org</a>
> used under the
> <a href="https://creativecommons.org/licenses/by/4.0/">CC-BY-4.0</a>
> license.</a>

Utilities to work with [GH Archive](https://www.gharchive.org/), project to record the public GitHub timeline, 
archive it, and make it easily accessible for further analysis.

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

## Incomplete repo language data

Language data is not included in events.

There is incomplete (only 3 million repos) public dataset:
```sql
SELECT * FROM `bigquery-public-data.github_repos.languages`;
```

However, many popular repositories are missing and manual data retrieval is required.

### Source

Programming languages by repository as reported by 
GitHub's https://developer.github.com/v3/repos/#list-languages API

### Properties
* No repo id, just name
* Probably no removed or renamed repos
* ~3 million entries
* Language data is in array (language name, bytes)