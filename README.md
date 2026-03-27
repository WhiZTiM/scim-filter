## SCIM-Filter

A simple golang library that implements the SCIM filters .

### Rationale
I like the semantics and the ergonomics of SCIM filters, and I find it really useful in APIs,
specifically, when implementing a tree walker to generate SQL queries.

### Usage  

```go
expr, err := scim_filter.Parse("userName eq \"john\" and (emails[type eq \"work\"] or emails[type eq \"home\"])")
if err != nil {
	panic(err)
}

switch e := expr.(type) {
case *scim_filter.AndExpr:
	// do something 
case *scim_filter.OrExpr:
	// do something
case *scim_filter.AttrExpr:
    // do something
}
```

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.