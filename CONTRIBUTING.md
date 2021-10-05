# Contributing to helm-docs

## Testing

### Benchmarks

If you are working on a feature that is likely to impact performance, consider running benchmarks
and comparing the results before and after your change.

To run benchmarks, run the command:

```
go test -run=^$ -bench=. ./cmd/helm-docs
```
