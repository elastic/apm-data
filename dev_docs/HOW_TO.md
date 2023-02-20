## How to add a new intake field

Adding a new intake field requires changes to both apm-data and apm-server repositories. 

### 1. In apm-data (this repo)

1. Add the new field to modeldecoder `model.go` such that the field is parsed from JSON.
   - Intake v2: [input/elasticapm/internal/modeldecoder/v2/model.go](/input/elasticapm/internal/modeldecoder/v2/model.go)
   - RUM v3: [input/elasticapm/internal/modeldecoder/rumv3/model.go](/input/elasticapm/internal/modeldecoder/rumv3/model.go)
2. Run `make generate` to generate the corresponding `model_generated.go` from the modified `model.go` in step 1.
3. Run `make update-licenses` to add license header to `model_generated.go` generated in step 2.
4. Add the new field to the corresponding file in `model/`, e.g. `span.go`.
   1. Add the field to the struct.
   2. Modify the `fields` function.
   3. Add test to the corresponding `*_test.go` file.
5. Modify the modeldecoder `decoder.go` to map the modeldecoder model in step 1 to the internal model in step 4.
   - Intake v2: [input/elasticapm/internal/modeldecoder/v2/decoder.go](/input/elasticapm/internal/modeldecoder/v2/decoder.go)
   - RUM v3: [input/elasticapm/internal/modeldecoder/rumv3/decoder.go](/input/elasticapm/internal/modeldecoder/rumv3/decoder.go)
6. Run `make test`

### 2. In [apm-server](https://github.com/elastic/apm-server/)

1. Use the modified apm-data by replacing the apm-data dependency with the local one.
   - Run

         go mod edit -replace=github.com/elastic/apm-data=/path/to/your/apm-data
         make update
2. Modify [apmpackage](https://github.com/elastic/apm-server/tree/main/apmpackage) to add the field to Elasticsearch mapping.
   1. Find the corresponding data stream directory in `apmpackage/apm/data_stream/`.
   2. Add the field under the YAML file under `fields/` in the data stream directory, e.g. `apmpackage/apm/data_stream/traces/fields/fields.yml`.
   3. Update apmpackage changelog `apmpackage/apm/changelog.yml`
3. Modify system test to ensure the field works end-to-end.
   1. Identify a `.ndjson` file to change under `TestIntake` in `systemtest/intake_test.go`, e.g. `events.ndjson` 
   2. Add the field to the NDJSON file.
   3. Run `make system-test` or `cd systemtest && go test ./... -run TestIntake`
   4. System test above should fail because there's a new field in the Elasticsearch documents. If it doesn't fail, check the code.
   5. Run `make check-approvals` to review and accept the changes in the Elasticsearch documents.

### 3. Finalize PRs

1. Create a PR in apm-data, and have it reviewed and merged.
2. In apm-server, bump apm-data dependency.
   - Run

         go mod edit -dropreplace=github.com/elastic/apm-data
         go get github.com/elastic/apm-data@<merged_git_commit_hash>
         make update
3. Create a PR in apm-server.

### Example PRs to add a field:
1. [apm-data PR](https://github.com/elastic/apm-data/pull/3)
2. [apm-server PR](https://github.com/elastic/apm-server/pull/9850)

