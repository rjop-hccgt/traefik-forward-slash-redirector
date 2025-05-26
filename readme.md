This repository includes traefik plugin which will add a final forward slash to a URL if not present. It will do a
temporary (HTTP 302 Found) or permanent (HTTP 301 Moved Permanently) redirect.

[![Build Status](https://github.com/rjop-hccgt/traefik-forward-slash-redirector/workflows/Main/badge.svg?branch=master)](https://github.com/rjop-hccgt/traefik-forward-slash-redirector/actions)

## Usage

For a plugin to be active for a given Traefik instance, it must be declared in the static configuration.

Plugins are parsed and loaded exclusively during startup, which allows Traefik to check the integrity of the code and
catch errors early on.
If an error occurs during loading, the plugin is disabled.

For security reasons, it is not possible to start a new plugin or modify an existing one while Traefik is running.

Once loaded, middleware plugins behave exactly like statically compiled middlewares.
Their instantiation and behavior are driven by the dynamic configuration.

Plugin dependencies must be [vendored](https://golang.org/ref/mod#vendoring) for each plugin.
Vendored packages should be included in the plugin's GitHub
repository. ([Go modules](https://blog.golang.org/using-go-modules) are not supported.)

### Configuration

For each plugin, the Traefik static configuration must define the module name (as is usual for Go packages).

The following declaration (given here in YAML) defines a plugin:

```yaml
# Static configuration

experimental:
  plugins:
    forward-slash-redirector:
      moduleName: github.com/rjop-hccgt/traefik-forward-slash-redirector
      version: v1.0.0
```

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the
`http.middlewares` section:

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - forward-slash-redirector

  services:
    service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000

  middlewares:
      plugin:
        forward-slash-redirector:
          permanent: false
```

### Local Mode

Traefik also offers a developer mode that can be used for temporary testing of plugins not hosted on GitHub.
To use a plugin in local mode, the Traefik static configuration must define the module name (as is usual for Go
packages) and a path to a [Go workspace](https://golang.org/doc/gopath_code.html#Workspaces), which can be the local
GOPATH or any directory.

The plugins must be placed in `./plugins-local` directory,
which should be in the working directory of the process running the Traefik binary.
The source code of the plugin should be organized as follows:

```
./plugins-local/
    └── src
        └── github.com
            └── traefik
                └── plugindemo
                    ├── demo.go
                    ├── demo_test.go
                    ├── go.mod
                    ├── LICENSE
                    ├── Makefile
                    └── readme.md
```

```yaml
# Static configuration

experimental:
  localPlugins:
    example:
      moduleName: github.com/rjop-hccgt/traefik-forward-slash-redirector
```

(In the above example, the `plugindemo` plugin will be loaded from the path
`./plugins-local/src/rjop-hccgt/traefik-forward-slash-redirector`.)

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-plugin

  services:
    service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000

  middlewares:
    my-plugin:
      plugin:
        forward-slash-redirector:
          permanent: false
```

