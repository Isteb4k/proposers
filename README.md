## Usage

```bash
go run . -case=proposers -from=0
go run . -case=slashes -from=8162625
go run . -case=passes -from=8162625
```

## Flags

-case=proposers - to display a table of validators with the number of mass passes when they are proposed, run the command
-case=passes - to display information at whose proposer the validator skips blocks
-case=slashes - to display validators who would have received a slash if there was no grace period

-from=0 - to get statistics for the entire history of collecting statistics
-from=8162625 - to get statistics starting from block 8162625
