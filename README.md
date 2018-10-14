[![Build Status](https://travis-ci.com/otaviof/hosts.svg?branch=master)](https://travis-ci.com/otaviof/hosts)

# `hosts`

Command line utility to generate your `/etc/hosts`, based in a combination files. It also
supports reading a external HTTP resource to populate local definitions, like block-lists for
instance.

The objective of `hosts` is to allow keeping individual projects, or contexts, hosts definitions
into their own files, and also move this type of data back to user home.

## Installing

The easiest way to install `hosts` is via `go get`:

``` bash
go get -u github.com/otaviof/hosts/cmd/hosts
```

Alternatively, when you cloning the repository:

``` bash
make bootstrap
make
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
    # search and replace values
    mappings:
      # search for
      - search: 127.0.0.1
        # replace with
        replace: 0.0.0.0
    # skip lines from returned body, based in a list of regular expressions
    skip:
      - ^#.*?$
      - ^\s+#.*?$
      - ^.*?(local|localhost|broadcasthost|ip6).*?$
```

The default place for this configuration is `/usr/local/etc/hosts.yaml`, or alternatively you
maybe use `--config` to inform a different location, in command-line.

To initialize your configuration, use:

``` bash
cp -v configs/hosts.yaml /usr/local/etc/
```

### Host Files

This application will look for `.host` files in the `hosts.baseDirectory` location. You can find
example of those files in `example/hosts-dir`, the formatting is the same than `/etc/hosts` file.

For instance:

```
127.0.0.1 hostname.local hostname
```

### External Resource

It's a common use-case to map malicious or advertising related addresses in `/etc/hosts` to
`0.0.0.0` or `127.0.0.1`, therefore you block any communication from your device to those endpoints.

Online communities like for instance [SomeOneWhoCares.org](https://someonewhocares.org) and
[uBlock Assets](https://github.com/uBlockOrigin/uAssets), are providing a up-to-date list of
[hosts](https://github.com/uBlockOrigin/uAssets/tree/master/thirdparties) that users can adopt,
although, in many cases you may need modifications, and may want to skip certain entries as well.

Therefore, `hosts` provide a way to load the external resource and apply `mappping`s and `skip`
certain lines. Please consider [configuration](#configuration) section.

## Usage

This command-line utility will inspect `hosts.baseDirectory`, and the `*.host` files found over
there are combined to create a new `/etc/hosts` files, accordingly to configuration.

The sequence of files in this directory is kept based on alpha-numeric ordering, therefore it's
encouraged to name files starting with numbers, like `00-first.host`, `10-second.host` and so
forth, since the sequence will be maintained.

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

When update external resources it will save date to configured file, and to take part of
`/etc/hosts` file, you must use `apply` command again.
