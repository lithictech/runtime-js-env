[![](https://godoc.org/github.com/lithictech/runtime-js-env?status.svg)](http://godoc.org/github.com/lithictech/runtime-js-env)

![Build](https://github.com/lithictech/runtime-js-env/workflows/Build/badge.svg)

# runtime-js-env

Write dynamic config into a static index.html for use in Single Page Applications,
**without** requiring additional requests from the client or runtime complexity.

## What does runtime-js-env do?

Usually SPAs have some configuration written into the JS files
that is read from `process.env` at build time.
This is usually desirable because it allows you to put your SPA into a CDN and cache forever.

Sometimes, though, you want to configure your frontend from the environment at runtime.

There are many approaches to this, but most of them involve the client loading a `config.js` file
before or alongside the main JS bundle. Or config is templated in for each request.

Both of these are too complex for a thing that can be solved in just a bit of code.

## Installation

`runtime-js-env` is a Golang program, so you need some version of Go.
Installation is the same as any other Golang CLI:

```
$ go get github.com/lithictech/runtime-js-env

$ runtime-js-env --help
NAME:
   runtime-js-env - A new cli application

USAGE:
   runtime-js-env [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --index value, -i value            Path to the index.html file. Default to index.html in pwd. (default: "index.html")
   --window-var-name value, -w value  Attribute name for the config object. (default: "_jsenv")
   --env-prefixes value, -p value     Environment variable prefixes to copy into the config object. (default: "REACT_APP_", "NODE_", "HEROKU_")
   --indent value, -t value           Indentation for each line in the config script tag. (default: "  ")
   --help, -h                         show help (default: false)
```

## Usage

After installation, there are two things you need to do to use `runtime-js-env`:

First, when your app (whatever is serving your HTML) boots,
call `runtime-js-env` with the path to your index.html file
(and any other options you need).

Then just call `runtime-js-env` with the location of your `index.html` file:

```
$ REACT_APP_MYVAR=hello runtime-js-env -i public/index.html
$ cat public/index.html
<!DOCTYPE html><html><head><script id="jsenv">
  window._jsenv = {
    "REACT_APP_MYVAR": "hello",
  };
</script></head><body><div id="___gatsby"></div><script src="/commons.js"></script></body></html>
```

This will overwrite your index.html file with a version
that injects a `<script>` tag in `<head>` which includes your config
(it's safe to call multiple times, and uses Go's HTML5 parser so should be valid for whatever you throw at it).

Then, in your JavaScript, wherever you use `process.env` you should prefer `window._jsenv`:

```js
// config.js
export function readConfig() {
    const env = window._jsenv || process.env;
    return {
      myVar: env.REACT_APP_MY_VAR || "default",
    };
}
```

That's all- basically wherever you use `process.env`,
you should use `window._jsenv` and fall back to `process.env`.
You can write some helpers to do this, or have a centralized config file, whatever you prefer.

## Development

Clone it down and check out the Makefile. It should be pretty self-explanatory.
