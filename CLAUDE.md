# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

bookget is a Go-based digital ancient book download tool supporting 50+ digital libraries worldwide. It downloads high-resolution images from various academic and cultural institutions' digital collections.

## Build Commands

```bash
# Build for current platform
make build

# Cross-platform builds
make linux-amd64      # Linux x64
make darwin-amd64      # macOS Intel
make darwin-arm64      # macOS Apple Silicon
make windows-amd64     # Windows x64
make release           # All platforms

# Clean build artifacts
make clean
```

## Testing

Run tests with:
```bash
go test ./pkg/hash/
go test ./pkg/quickxorhash/
```

## Development Commands

```bash
# Run directly during development
go run ./cmd/

# Run with specific URL
go run ./cmd/ -u "https://example.com/book-url"

# Interactive mode (default)
go run ./cmd/

# Image downloader mode
go run ./cmd/ -d 1
```

## Architecture

### Core Components

- **cmd/bookget.go**: Main entry point with four run modes (single URL, batch URLs, interactive, interactive image)
- **router/**: Factory pattern for routing URLs to appropriate site handlers based on domain matching
- **app/**: Site-specific downloaders for 50+ digital libraries (e.g., harvard.go, nlc.go, waseda.go)
- **pkg/**: Reusable packages for HTTP handling, file operations, progress bars, and concurrent downloading
- **model/**: Data structures for site-specific API responses
- **config/**: Configuration management with INI file support

### URL Processing Flow

1. URLs are parsed and routed through `router.FactoryRouter()`
2. Domain-based matching determines the appropriate app handler
3. Site handlers implement `RouterInit` interface with `GetRouterInit()` method
4. Each handler manages site-specific authentication, pagination, and image extraction

### Site Support

The app/ directory contains dedicated handlers for institutions including:
- **Academic**: Harvard, Princeton, Library of Congress, Oxford Bodleian
- **Asian**: National Diet Library Japan, Chinese University of Hong Kong, Waseda University
- **European**: Berlin State Library, Austrian National Library, British Library
- **Special formats**: IIIF, DZI, image downloaders

### HTTP and Download Infrastructure

- **pkg/gohttp/**: Multi-threaded downloading with chunking and resume support
- **pkg/chttp/**: Cookie and header management for authenticated sessions
- **pkg/queue/**: Concurrent job processing with configurable thread pools
- **pkg/progressbar/**: Terminal progress visualization

### Configuration

Configuration via INI files in config/ with support for:
- Proxy settings (respects HTTP_PROXY/HTTPS_PROXY environment variables)
- Custom user agents and request headers
- Cookie persistence for authenticated sessions
- Download concurrency and thread management

### Key Design Patterns

- **Factory Pattern**: `router.FactoryRouter()` creates appropriate site handlers
- **Strategy Pattern**: Each site handler implements domain-specific download logic
- **Concurrent Processing**: Queue-based job management for parallel downloads
- **Auto-detection**: Automatic site identification from URL patterns and content types