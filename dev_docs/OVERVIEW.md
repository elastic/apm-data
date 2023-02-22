# Overview

## Directories

- `input/`
  - Struct definitions of input models, e.g. Intake v2, RUM v3, OTel
  - Decoding and translation logic from these input models to the internal model `APMEvent`
- `model/`
  - Struct definitions of the internal model `APMEvent`
  - JSON serialization logic for indexing into Elasticsearch