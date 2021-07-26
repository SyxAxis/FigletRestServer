# FigletRestServer
Produce figlet output from a REST server 

Yet another example of me wasting my time while learning Go.

Fire it up 

Either get the list of fonts:

```
http://localhost:8888/v1/getfontlist
```

Or send a JSON request to generate

```
http://localhost:7777/v1/genmsg
```

```
{
    "fontname":"red_phoenix",
    "message":"REST call"
}
```
