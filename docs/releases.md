## Releases

### Creating a release

The project uses a GitHub Actions workflow (`.github/workflows/release.yml`) that triggers on tags matching `v*`.

Example:

```bash
git tag -a v1.2.0 -m "v1.2.0"
git push origin v1.2.0
```

### What gets built

The release workflow builds (in a single job) the following targets:

- Linux: `amd64`, `arm64`
- macOS (Darwin): `amd64`, `arm64`
- Windows: `amd64`, `386`

Artifacts are uploaded to the GitHub Release along with `sha256sums.txt`.

