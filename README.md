<p align="center">
    <a alt="GoReport" href="https://goreportcard.com/report/github.com/otaviof/hosts">
        <img alt="GoReport" src="https://goreportcard.com/badge/github.com/otaviof/hosts">
    </a>
    <a href="https://godoc.org/github.com/otaviof/hosts/pkg/hosts">
        <img alt="GoDoc Reference" src="https://godoc.org/github.com/otaviof/hosts/pkg/hosts?status.svg">
    </a>
    <a alt="CI Status" href="https://travis-ci.com/otaviof/hosts">
        <img alt="CI Status" src="https://travis-ci.com/otaviof/hosts.svg?branch=master">
    </a>
</p>

# `hosts`

Command line utility to generate your `/etc/hosts`, based in a combination files. It also supports
reading an external HTTP resource to populate local definitions, like block-lists for instance.

The objective of `hosts` is to allow keeping individual projects, or contexts, hosts definitions
into their own files, and also move this type of data back to user home.

## Installing

The easiest way to install `hosts` is via `go get`:

``` bash
go get -u github.com/otaviof/hosts/cmd/hosts
```

Alternatively, when cloning this repository, execute:

``` bash
make install
```

## Configuration

The following example configuration has fields description as comments, please consider.

``` yaml
hosts:
  # base directory where to look for `.host` files
  baseDirectory: ~/.hosts
  # final output file
  output: /etc/hosts

# external resources
external:
  # read "body" from external URL
  - url: https://someonewhocares.org/hosts/hosts
    # file name to save contents, under `hosts.baseDirectory`
    output: 99-blocks.host
    # transform downloaded contents line-by-line
    transform:
      # search using a regular expression
      - search: 127.0.0.1
        # replace with
        replace: 0.0.0.0
      # regular expression using a match group
      - search: (\w+.*?)#.*?$
        # replacing with match group
        replace: $1
      # when `replace` is empty, the line is skipped
      - search: ^#.*?$
      # skipping localhost related entries
      - search: ^.*?(local|localhost|broadcasthost|ip6).*?$
```

The default location the configuration file is at `/usr/local/etc/hosts.yaml`, or alternatively
you can employ `--config` parameter to inform a different location.

To start, copy the example configuration:

``` bash
cp -v configs/hosts.yaml /usr/local/etc/
```

### Host Files

This application will look for `.host` files in the `hosts.baseDirectory` location. You can find
example of those files in
[`test/hosts-dir`](https://github.com/otaviof/hosts/tree/master/test/hosts-dir), the formatting is
the same than `/etc/hosts` file.

For instance:

```
127.0.0.1 hostname.local hostname
```

### External Resource

It's a common use-case to map malicious or advertising related addresses in `/etc/hosts` to
`0.0.0.0` or `127.0.0.1`, therefore you block any communication from your device to those endpoints.

Online communities like for instance [SomeOneWhoCares.org](https://someonewhocares.org) and
[uBlock Assets](https://github.com/uBlockOrigin/uAssets), are providing a up-to-date list of
[hosts](https://github.com/uBlockOrigin/uAssets/tree/master/thirdparties) which users can adopt,
although, in many cases you may need modifications, and may want to skip certain entries as well.

Therefore, `hosts` provide a way to load the external resource and apply regular expression based
transformation, which can modify contents and skip certain lines. Please consider
[configuration](#configuration) section.

## Usage

This command-line utility will inspect `hosts.baseDirectory`, and the `*.host` files found over
there are combined to create a new `/etc/hosts` files, accordingly to configuration.

The sequence of files in this directory is kept based on alpha-numeric ordering, therefore it's
encouraged to name files starting with numbers, like `00-first.host`, `10-second.host` and so forth.

The following parameters are applicable to all sub-commands:

- `--config`: alternative location of the configuration file;
- `--dry-run`: do not apply changes;
- `--help`: inline help message;

To generate your hosts file, use `apply` command:

``` bash
hosts apply --dry-run
```

And to read from `external` resource, use:

``` bash
hosts update --dry-run
```

Rinse and repeat.
