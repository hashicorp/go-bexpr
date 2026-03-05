## 0.1.16 (March 5, 2026)

### Improvements
- Adds "is nil" and "is not nil" selector. [[GH-129](https://github.com/hashicorp/go-bexpr/pull/129)]

### Bug Fixes
- Fixed a bug where using "is empty" or "is not empty" with a non-slice or non-map value would panic. [[GH-129](https://github.com/hashicorp/go-bexpr/pull/129)]

### Security

## 0.1.15 (October 17, 2025)

### Improvements
- Adds a default of 2 million for max evaluated expressions. [[GH-112](https://github.com/hashicorp/go-bexpr/pull/112)]

### Bug Fixes
- Fixes incorrect struct tag in README example. [[GH-76](https://github.com/hashicorp/go-bexpr/pull/76)]

### Security
