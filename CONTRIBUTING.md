### Contributing
Feel free to open issues.

### Adding a matcher — definition of done

A new matcher ships with all of the following:

1. **Doc-comment** with one-line semantics and a copy-pasteable example, plus a
   "prefer this over ..." line when the matcher supersedes a raw idiom or a
   `be.Not(...)` composition. Doc-comments are the primary discovery surface —
   agents and editors read `go doc`, not the README.
2. **Tests**: positive match, negative match, and (for matchers with custom
   failure output) the failure message.
3. **Catalog entry**: run `./generate-docs.sh` to regenerate `MATCHERS.md`
   (new matchers are picked up automatically; assign the matcher to an intent
   section in `internal/docgen/main.go` — unassigned ones land under "Other").
   If the matcher supersedes a raw idiom, add it to the `insteadOf` table there.
4. **belint rule** (`x/belint`), if the matcher replaces a raw expression that
   people would otherwise wrap in `be.True(...)` — point-of-use feedback is
   what makes both humans and LLM agents actually discover it.
5. **Root alias** only if the name is collision-free across packages AND the
   matcher is high-frequency. `be_time` is never aliased at root.

### Releasing
This repo is multi-module: the core (`.`) plus the opt-in `x/mock` and
`x/belint` plugins. Cut them all at once with:

```sh
make release VERSION=v1.0.0-rc.6
```

It tags the core, points `x/mock`'s core requirement at that tag (committing the
bump only if it changed), and tags `x/mock` and `x/belint` (the latter has no
core requirement to bump). Tags are created locally; the target prints the
`git push` command for you to run after review.

### TODO: stabilize with contributing guidelines