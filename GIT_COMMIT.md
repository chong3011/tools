# Git Commit Commands

This document describes common git commit commands used when contributing to this repository.

## Basic Commit

Stage and commit all changes:

```bash
git add .
git commit -m "your commit message"
```

Or stage specific files:

```bash
git add <file>
git commit -m "your commit message"
```

## Commit Message Guidelines

- Use a short, descriptive subject line (50 characters or less)
- Use the imperative mood: "Fix bug" not "Fixed bug"
- Reference related issues or pull requests where applicable

## Amending a Commit

To amend the most recent commit message or add forgotten changes:

```bash
git add <forgotten-file>
git commit --amend
```

## Viewing Commit History

```bash
git log --oneline
git log --oneline -10   # show last 10 commits
```

## Undoing Changes

Undo the last commit but keep changes staged:

```bash
git reset --soft HEAD~1
```

Undo the last commit and unstage changes:

```bash
git reset HEAD~1
```

## Pushing Changes

Push your branch to the remote:

```bash
git push origin <branch-name>
```

For more details on contributing to this project, see [CONTRIBUTING.md](CONTRIBUTING.md).
