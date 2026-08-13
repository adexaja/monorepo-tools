# Template source tree

This directory is the project skeleton used by `../create.sh`. If you are
working directly here, bootstrap the Bun workspace before running Moon:

```sh
bun install
moon run :test
moon run :lint
moon run :build
```

For a new project, use `../create.sh` so placeholders are rendered first.
