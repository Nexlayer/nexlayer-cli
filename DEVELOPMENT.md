# Development Guide for Nexlayer CLI

This document outlines the development workflow and processes for the Nexlayer CLI project.

## Branch Structure

We follow a modified GitFlow workflow with the following branches:

```
main                # Production-ready code
├── develop         # Integration branch for features
├── release/*       # Release candidate branches (e.g., release/v0.1.0)
├── feature/*       # New features
├── fix/*           # Bug fixes 
└── refactor/*      # Code refactoring
```

### Branch Naming

All branches should follow the naming convention:

- `feature/short-description` - For new features
- `fix/issue-description` - For bug fixes
- `refactor/component-name` - For code refactoring
- `release/vX.Y.Z` - For release candidates

Examples:
- `feature/api-authentication`
- `fix/deployment-url-handling`
- `refactor/api-client`

## Development Workflow

### Starting New Work

1. Always start from the latest `develop` branch:
   ```bash
   git checkout develop
   git pull
   ```

2. Create a new branch with the appropriate prefix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Make changes, commit frequently with descriptive messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Testing Your Changes

1. Run tests locally before pushing:
   ```bash
   make test
   make lint
   ```

2. Verify functionality using the E2E test script:
   ```bash
   /Users/salstagroup/CursorProjects/nexlayer-test-projects/test-e2e.sh
   ```

### Submitting Changes

1. Push your branch to GitHub:
   ```bash
   git push -u origin feature/your-feature-name
   ```

2. Open a Pull Request against the `develop` branch.
3. Ensure all CI checks pass.
4. Request a code review.
5. Address any feedback.

### Merging Strategy

1. All PRs should be merged via the GitHub UI using "Squash and Merge"
2. Delete the branch after merging

## Release Process

1. Create a release branch from `develop`:
   ```bash
   git checkout develop
   git checkout -b release/vX.Y.Z
   ```

2. Perform final testing on the release branch.
3. Update version numbers and CHANGELOG.md.
4. Open a PR from `release/vX.Y.Z` to `main`.
5. After merging to `main`, tag the release:
   ```bash
   git checkout main
   git pull
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   git push origin vX.Y.Z
   ```
6. Merge `main` back to `develop`:
   ```bash
   git checkout develop
   git merge main
   git push
   ```

## Code Review Guidelines

### For Authors

- Keep PRs focused and reasonably sized (under 500 lines when possible)
- Write clear descriptions explaining the changes
- Link to relevant issues
- Run tests locally before requesting review
- Respond to reviewer comments promptly

### For Reviewers

- Review PRs within 48 hours
- Be constructive and specific in feedback
- Comment inline for specific issues
- Approve only when all concerns are resolved

## Branch Cleanup

Regularly clean up merged branches:

```bash
# List branches that have been merged to main
git checkout main
git branch --merged

# Delete merged branches (after verifying the list)
git branch -d branch-name
```

For remote branches:
```bash
git push origin --delete branch-name
```

## Hotfix Process

For urgent fixes to production:

1. Create a hotfix branch from `main`:
   ```bash
   git checkout main
   git checkout -b hotfix/critical-bug
   ```

2. Fix the issue and open a PR to `main`
3. After merging to `main`, tag a patch release
4. Also merge the changes back to `develop` 