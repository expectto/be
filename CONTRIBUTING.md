### Contributing
Feel free to open issues.

### Releasing
This repo is multi-module: the core (`.`) plus the opt-in `x/mock` plugin. Cut
both at once with:

```sh
make release VERSION=v1.0.0-rc.6
```

It tags the core, points `x/mock`'s core requirement at that tag (committing the
bump only if it changed), and tags `x/mock`. Tags are created locally; the target
prints the `git push` command for you to run after review.

### TODO: stabilize with contributing guidelines