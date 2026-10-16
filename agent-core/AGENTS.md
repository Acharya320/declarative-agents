<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

# Agent Instructions

## Documentation authority

Markdown documents under `docs/` and `docs/constitutions/` in this repository are the source of truth for agent-core behavior contracts. `docs/constitutions/design.yaml` stays YAML because it declares every document type and its Markdown layout. When workspace-wide editor rules or other guidance disagrees with that corpus, follow the repository docs.

Place a new document at the path for its type. Types and locations are declared in `docs/constitutions/design.yaml`. Required fields live there; do not copy them here.

| Type | Location |
|------|----------|
| vision | `docs/VISION.md` |
| architecture | `docs/ARCHITECTURE.md` |
| specifications | `docs/SPECIFICATIONS.md` |
| roadmap | `docs/road-map.md` |
| pattern_language | `docs/pattern-language.md` |
| engineering_guideline | `docs/engineering/engNN-short-name.md` |
| srd | `docs/specs/software-requirements/srdNNN-short-name.md` |
| use_case | `docs/specs/use-cases/relNN.N-ucNNN-short-name.md` |
| test_suite | `docs/specs/test-suites/test-relNN.N-short-name.md` |
| audit_register | `docs/specs/audits/name-register.md` |
| audit_report | `docs/specs/audits/name.md` |
| semantic_model | `docs/specs/semantic-models/name.md` |
| config_format | `docs/specs/config-formats/name.md` |
| constitution | `docs/constitutions/name.md` |
| migration | `docs/migrations/name.md` |
| guide | `docs/guides/name.md` |

Each corpus document is Markdown with frontmatter and one-line marker comments above its headings; `../docs/engineering/eng10-markdown-spec-format.md` defines the format. After editing a document, run `mage fmt <path>` from the repository root to put it in canonical form, then `mage audit` to check it.

Do not create a document that matches no declared type. Register the type in `docs/constitutions/design.yaml` first, or do not add the file. `mage audit` enforces this.

## Package layout

`cmd` packages hold entry points and adapter wiring only, per `docs/constitutions/go-style.md`. Runtime logic belongs in `internal` packages.

## Non-Interactive Shell Commands

**ALWAYS use non-interactive flags** with file operations to avoid hanging on confirmation prompts.

Shell commands like `cp`, `mv`, and `rm` may be aliased to include `-i` (interactive) mode on some systems, causing the agent to hang indefinitely waiting for y/n input.

**Use these forms instead:**
```bash
# Force overwrite without prompting
cp -f source dest           # NOT: cp source dest
mv -f source dest           # NOT: mv source dest
rm -f file                  # NOT: rm file

# For recursive operations
rm -rf directory            # NOT: rm -r directory
cp -rf source dest          # NOT: cp -r source dest
```

**Other commands that may prompt:**
- `scp` - use `-o BatchMode=yes` for non-interactive
- `ssh` - use `-o BatchMode=yes` to fail instead of prompting
- `apt-get` - use `-y` flag
- `brew` - use `HOMEBREW_NO_AUTO_UPDATE=1` env var

## Session Completion

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

### Must-follow rules
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
