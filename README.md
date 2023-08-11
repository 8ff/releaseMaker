# ReleaseMaker
A simple command-line tool to manage GitHub releases. With this tool, you can create new releases, upload assets to existing releases, and delete releases by tag name.

## Installation
Grab a binary from releases or clone the repository and build the project using Go:
```bash
git clone https://github.com/8ff/releaseMaker.git
cd releaseMaker
go build
```

## Usage
Make sure to set your GitHub token as an environment variable:
```bash
export GITHUB_TOKEN=your_github_token
```

### Create a New Release
Create a new release for the specified repository:
```bash
./your-binary create [owner/repo] [tag] [name] [body]
```

- `owner/repo`: The owner and repository name, separated by a slash.
- `tag`: The tag name for the release.
- `name`: The name of the release.
- `body`: The body text of the release.

### Upload an Asset to a Release
Upload a file as an asset to an existing release by tag name:
```bash
./your-binary upload [owner/repo] [tag] [file] [assetName]
```

- `owner/repo`: The owner and repository name, separated by a slash.
- `tag`: The tag name of the release to upload the asset to.
- `file`: The path to the file to upload.
- `assetName`: The name of the asset in the release.

### Delete a Release
Delete an existing release by tag name:
```bash
./your-binary delete [owner/repo] [tag]
```

- `owner/repo`: The owner and repository name, separated by a slash.
- `tag`: The tag name of the release to delete.

## Contributing
Feel free to open issues or submit pull requests. All contributions are welcome!

## License
[GNU Affero General Public License v3.0](LICENSE)