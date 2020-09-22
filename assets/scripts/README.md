# Scripts

## codegen.py

The purpose of `codegen.py` script is generate the code necessary
for fully unpack JSON string into a data structure. The script
ensures that Go data structures fully cover the expected API
server output.

As of writing this `README.md`, the library have Go data structures containing
less attributes than the data being passed from an API server. The script
executes the following steps:

* Loads JSON input containing a data strcuture received from an API server
* Reads each key in the JSON's dictionaries and checks whether there is
  a corresponding attribute in a target Go data structure
* If necessary, generates `struct` Go attributes and test actions

For example, suppose the `ComputerSystem` object does not have `AssetTag`
attribute.

Run the following command to detect the lack of coverage:

```bash
assets/scripts/codegen.py \
  --json-input ./assets/responses/computer_system_1.json \
  --go-source ./pkg/client/computer_system.go \
  --parser-go-struct computerSystemResponse \
  --output-go-struct ComputerSystem \
  --debug
```

However, the attribute exists in the JSON output (see `assets/responses/computer_system_1.json`):

```json
{
  "@odata.context": "/redfish/v1/$metadata#ComputerSystem.ComputerSystem",
  "@odata.id": "/redfish/v1/Systems/System.Embedded.1",
  "@odata.type": "#ComputerSystem.v1_5_1.ComputerSystem",
  "AssetTag": "KN23N857Z"
}
```

The script generates the following Go code for `pkg/client/computer_system.go`:

```golang
type computerSystemResponse struct {
    AssetTag    string
}

type ComputerSystem struct {
    AssetTag    string    `yaml:"asset_tag" json:"asset_tag" xml:"asset_tag"`
}

func newComputerSystemFromBytes(s []byte) (*ComputerSystem, error) {
    cs.AssetTag = response.AssetTag
}
```

Additionally, the script generates the following Go code for `pkg/client/computer_system_test.go`:

```golang
if (resource.AssetTag != test.exp.AssetTag) && !test.shouldFail {
    t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
        i, fp, "AssetTag", resource.AssetTag, test.exp.AssetTag)
    testFailed++
    continue
}
```
