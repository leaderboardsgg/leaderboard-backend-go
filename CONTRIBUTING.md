# Contributing to leaderboard-backend
We appreciate your help!

## Before filing an issue
If you are unsure whether you have found a bug, please consider asking in [our discord](https://discord.gg/TZvfau25Vb) first.

Similarly, if you have a question about a potential feature, [the discord](https://discord.gg/TZvfau25Vb) can be a fantastic resource for first comments.

## Filing issues
Filing issues is as simple as going to [the issue tracker](https://github.com/speedrun-website/leaderboard-backend/issues), and adding an issue using one of the below templates.

### Feature Request / Task
```
{Feature Request/Task}: {short description}
---
{Detailed description}

### Affected Functionality
{Any known functionality impacted, or Unknown if further research needs to be done}

### Other Relevant Issues
{Links to Relevant Issues}
```

### Bugs
```
Bug: {short description}
---
{Summary of bug}

### Step(s) to Reproduce
{Numbered list of step(s)}

### Expected Result
{Summary of expected result}

### Actual Outcome
{Description of actual outcome}
```

## Contributing code
### Example code contribution flow
1. Make a fork of this repo.
1. Name a branch on your fork something descriptive for this change (eg. `UpdateMakefile`).
1. Commit your changes (Tip! Please read our [Style guide](#style-guide) to help the pull request process go smoothly).
1. Verify your changes work with `make test`.
1. Push your branch.
1. Open a pull request to `speedrun-website/leaderboard-backend`.
1. Get your pull request approved.
1. Get someone to click `Rebase and merge`.
1. Celebrate your amazing changes! 🎉

## Style guide
### General
- Be inclusive, this is a project for everyone.
- Be descriptive, it can be hard to understand abbreviations or short-hand.

### GoLang
- Add tests for any new feature or bug fix, to ensure things continue to work.
- Comments should be full sentences, starting with a capital letter and ending with punctuation.
- Comments above a func or struct should start with the name of the thing being described.
- Wrap errors before returning with `oops.Wrapf(` to help build useful stacks for debugging.
- Early returns are great, they help reduce nesting!
- Avoid `interface{}` where possible, if we need such a generic please add a comment explaining what the real type is.

### Git
- Try to have an informative branch name for others eg. `LB-{issue number}-{ghusername}`.
  - Do not make pull requests from `main`.
  - Do not include slashes in your branch name.
    - Nested paths can act strange when other people start looking at your branch.
- Try to keep commit line length below 80 characters.
- All commit titles should be of the format `{area} {optional sub-area}: commit description`.
  - This will help people reading through commits quickly find the relevant ones.
  - Some examples might include:
    - `go data: add support for sorting by high-score`
    - `makefile: add docker build commands`
- Commits should be as [atomic](https://www.freshconsulting.com/insights/blog/atomic-commits/) as possible.