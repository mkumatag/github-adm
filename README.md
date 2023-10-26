# github-adm

github-adm is a command-line tool written in Go that helps you manage and synchronize labels in your GitHub repositories.

## Installation

To install github-adm, you'll need to have Go installed. You can install it using Go's package manager, `go get`:

```bash
go get -u github.com/mkumatag/github-adm
```

## Usage

### `sync-labels`

The `sync-labels` command allows you to synchronize GitHub labels with labels defined in a JSON file. This can be useful for maintaining a consistent set of labels across multiple repositories or simply ensuring that your labels are in sync with a predefined set.

#### Usage

```bash
github-adm sync-labels [options]
```

#### Options
`github-adm sync-labels --help` can be used to list all the options.


#### Example

To synchronize labels for the repository "example/repo" with labels defined in a JSON file and remove any labels not in the JSON file:

```bash
GH_TOKEN=<GH_TOKEN> github-adm sync-labels --base-url https://github.ibm.com/api/v3 --upload-url https://uploads.github.ibm.com/ --org org --repo repo --manifest labels.json --delete-out-of-sync`
```

The `labels.json` file should contain label definitions in the following format:

```json
[
  {
    "name": "bug",
    "color": "d73a4a"
  },
  {
    "name": "enhancement",
    "color": "a2eeef"
  },
  {
    "name": "documentation",
    "color": "0075ca"
  }
  // Add more labels as needed
]
```

## Building from Source

If you prefer to build github-adm from source, you can use the following steps:

1. Clone the repository:

```bash
git clone https://github.com/mkumatag/github-adm.git
```

2. Build the project:

```bash
cd github-adm
go build
```

3. Install the binary:

```bash
go install
```

## Contributing

If you'd like to contribute to this project, please check out the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines.

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](LICENSE) file for details.

## Contact

If you have any questions or need assistance, feel free to contact us at [manjunath.cse@gmail.com](mailto:manjunath.cse@gmail.com).

Happy labeling with github-adm!
```
