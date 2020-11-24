<p align="center">
    <a alt="GoReport" href="https://goreportcard.com/report/github.com/otaviof/hosts">
        <img alt="GoReport" src="https://goreportcard.com/badge/github.com/otaviof/hosts">
    </a>
    <a href="https://pkg.go.dev/github.com/otaviof/hosts">
        <img alt="GoDoc Reference" src="https://godoc.org/github.com/otaviof/hosts/pkg/hosts?status.svg">
    </a>
    <a alt="CI Status" href="https://travis-ci.com/otaviof/hosts">
        <img alt="CI Status" src="https://travis-ci.com/otaviof/hosts.svg?branch=master">
    </a>
</p>

# `hosts`

Command line utility to generate your `/etc/hosts` (and DNSMasq files), based on a combination of
configuration files. It also supports reading an external HTTP resource to populate local
definitions, like block-lists for instance.

The objective of `hosts` is to allow keeping individual projects, or contexts, hosts definitions
into their own context directory and files, and also move this type of data back to user home.

## Installing

The easiest way to install `hosts` is via `go get`:

```sh
go get -u github.com/otaviof/hosts/cmd/hosts
```

Alternatively, when cloning this repository, execute:

```sh
make install
```

## Usage

To start, copy the example configuration:

```sh
mkdir ~/.hosts
cp -v test/hosts-dir/hosts.yaml ~/.hosts
${EDITOR} ~/.hosts/hosts.yaml
```

For the command line parameters, please use `hosts --help` and inspect the sub-commands `update` and
`apply` as well.

The daily basis usage would be to edit files under `~/.hosts` (as the
[example here](https://github.com/otaviof/hosts/tree/master/test/hosts-dir)), and run:

```sh
# updating external data sources
hosts update

# creating the target (output) file(s)
sudo hosts apply
```

In daily basis, running `hosts update` would a optional task, while `hosts apply` must take place
every time the output files must be updated.

## Configuration

The configuration file is named `hosts.yaml` and must be located inside the base-directory
(`--base-dir`) informed to `hosts`. By default, `hosts` expect to find a `~/.hosts` base directory.

The following describes each property in the configuration file, and how to use them:

- `.hosts`:
  - `.input`:
    - `.sources[]`: list of external data sources;
      - `.name`: block name, an optional short description;
      - `.url`: resource location, URL to data to be loaded;
      - `.file`: file name to store data retrieved from URL;
    - `.transformations[]`: list of transformations for external data sources;
      - `.search`: regular-expression applied to each line;
      - `.replace`: string to replace `.search` findings with. When `.replace` is empty, the line is
      skipped;
  - `.output[]`: list of output files created by this application;
    - `.name`:  block name, an optional short description;
    - `.path`: fullpath to file, including extension;
    - `.with`: regular-expression, matches file names that will be included on output;
    - `.without`: regular-expression, matches file names that will be excluded from output;

A example of `hosts.yaml` file is:

```yml
---
hosts:
  input:
    sources:
      - name: uBlockOrigin
        url: https://github.com/uBlockOrigin/uAssets/raw/master/thirdparties/www.malwaredomainlist.com/hostslist/hosts.txt
        file:  99-blocks.host
    transformations:
      - search: "127.0.0.1"
        replace: "0.0.0.0"
      - search: ^.*?(local|localhost|broadcasthost|ip6).*?$
      - search: ^\s+#.*?$
  output:
    - name: etc-hosts
      path: /etc/hosts
      without: 99-blocks.*
      mode: 0644
    - name: dnsmasq-blocks
      path: /etc/dnsmasq.d/blocks.conf
      dnsmasq: true
      with: 99-blocks.*
```

After editing the configuration file run `hosts config --validate`, so you can informed about
possible errors early on.

### Host Files

This application will look for `.host` files in the base directory location. You can find example
of those files in [`test/hosts-dir`](https://github.com/otaviof/hosts/tree/master/test/hosts-dir),
the formatting is the same [than `/etc/hosts` file](https://man7.org/linux/man-pages/man5/hosts.5.html).

For instance:

```
127.0.0.1 hostname.local hostname
```

### DNSMasq

Alternatively, instead of generating `/etc/hosts` format, this application can generate the format
employed on [`dnsmasq.conf` files ](http://www.thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html).
The flag to turn this formatter on is under `hosts.output` as the following example:

```yml
---
hosts:
  output:
    - name: dnsmasq-blocks
      dnsmasq: true
      path: /etc/dnsmasq.d/blocks.conf
```

### External Resources

It's a common use-case to map malicious or advertising related addresses in `/etc/hosts` to
`0.0.0.0` or `127.0.0.1`, therefore you block any communication from your device to those endpoints.

Online communities like for instance [SomeOneWhoCares.org](https://someonewhocares.org) and
[uBlock Assets](https://github.com/uBlockOrigin/uAssets), are providing a up-to-date list of
[hosts](https://github.com/uBlockOrigin/uAssets/tree/master/thirdparties) which users can adopt,
although, in many cases you may need modifications, and may want to skip certain entries as well.

Therefore, `hosts` provide a way to load the external resource and apply regular expression based
transformation, which can modify contents and skip certain lines. Please consider
[configuration](#configuration) section.
