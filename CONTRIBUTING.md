# Contributing Guidelines

## Branches

- **main** – Stable production-ready code.
- **develop** – Active development branch.
- **feature/<name>** – New features.
- **fix/<name>** – Bug fixes.

## Commit Message Convention

    Please follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard. The commit message format:

```<type>(<scope>): <short description>```


Examples:

- `feat(api-gateway): add new routing rule`
- `fix(streaming-service): correct environment variable naming`

## Pull Requests and Code Reviews

- All changes must be submitted via a Pull Request (PR) from a feature/fix branch into `develop`.
- The PR must include a detailed description of the changes.
- Direct commits to `main` are not allowed; changes must be merged after a successful code review.
