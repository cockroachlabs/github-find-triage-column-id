# github-find-triage-column-id

## Setup

First, export GITHUB_API_KEY on your local environment:

```sh
export GITHUB_API_KEY = '<api key>'

```
## Usage

For organization based projects:

```sh
github-get-column-id --project "Spatial" --column "Backlog"
```

For repo based projects:

```sh
github-get-column-id --repo "cockroach" --project "Bazel" --column "To do"
```

