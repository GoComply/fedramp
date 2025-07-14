# Go Version Requirement

## Minimum Required Version: Go 1.19

This project requires **Go 1.19 or higher** due to dependencies that use features introduced in newer Go versions.

## CI Configuration Update Required

The current CI pipeline is using Go 1.14.15, which causes build failures with the error:
```
io/fs: package io/fs is not in GOROOT
```

This is because the `io/fs` package was introduced in Go 1.16, and our dependencies (particularly `fsnotify`) require it.

## Fix for GitHub Actions

The GitHub Actions workflow needs to be updated to use Go 1.19 or higher. Update the workflow file with:

```yaml
- name: Set up Go
  uses: actions/setup-go@v4
  with:
    go-version: '1.19'
```

## Build Error Resolution

If you encounter the build error locally, ensure you have Go 1.19+ installed:
```bash
go version  # Should show go1.19 or higher
```

## Dependencies Requiring Newer Go

The following dependencies require Go 1.16+:
- `github.com/fsnotify/fsnotify` (requires `io/fs`)
- Various server implementation dependencies added for R5 Balance features 