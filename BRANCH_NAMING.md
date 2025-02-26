# Nexlayer CLI Branch Naming Guide

This document outlines the branch naming conventions for the Nexlayer CLI repository.

## Branch Structure

```
main              # Production-ready code
├── develop       # Integration branch for features
├── feature/*     # Feature branches 
├── fix/*         # Bug fixes
├── refactor/*    # Code refactoring
├── docs/*        # Documentation updates
├── test/*        # Test additions or updates
├── chore/*       # Maintenance tasks
└── release/*     # Release preparation
```

## Branch Naming Conventions

All branch names should follow this pattern:
```
<type>/<short-description>
```

Where:
- `<type>` is one of the predefined types below
- `<short-description>` is a brief, hyphenated description of the change

### Branch Types

- `feature/`: New functionality
- `fix/`: Bug fixes or issues
- `refactor/`: Code changes that neither add features nor fix bugs
- `docs/`: Documentation changes only
- `test/`: Adding or updating tests
- `chore/`: Maintenance tasks, build changes, etc.
- `release/`: Branches for release preparation

### Examples

✅ Good branch names:
- `feature/custom-domains`
- `fix/deployment-url-handling`
- `refactor/api-client`
- `docs/improve-installation-guide`
- `test/add-deployment-tests`
- `chore/update-dependencies`
- `release/v1.0.0`

❌ Poor branch names:
- `my-branch`
- `bug-fix`
- `updates`
- `john-work`
- `quick-fix`

## Branch Lifecycle

1. **Creation**: Branch from `develop` (or `main` for hotfixes)
2. **Development**: Make changes, commit regularly
3. **PR**: Submit PR back to `develop` (or `main` for hotfixes)
4. **Review & Merge**: After approval, merge via PR
5. **Cleanup**: Delete branch after successful merge

## Stale Branch Policy

Branches that haven't been updated in 30 days will be considered stale. Team members will be notified, and stale branches may be archived or deleted after 60 days of inactivity unless marked for preservation. 