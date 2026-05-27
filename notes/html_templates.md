
# Dynamic HTML templates

Use `html/template` package, and not `text/template` package.

## custom template functions

Two steps:
- create a `template.FuncMap` map, containing the custom `humanDate()` function
- use the `template.Funcs()` method to register
