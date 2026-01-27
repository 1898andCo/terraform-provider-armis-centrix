# MITRE ATT&CK Label Format Validation

## Summary

The MITRE ATT&CK label validator in [internal/verify/validate.go](internal/verify/validate.go) is **CORRECT** for the Armis-specific format used in the Terraform provider.

## Findings

### Armis API Behavior

The Armis API has two different representations for MITRE labels:

**Input Format (Write Operations):**
- Accepts `[]string` with combined format: `"Enterprise.TA0009.T1056.001"`
- This is an Armis-specific format that combines Matrix + Tactic ID + Technique ID + Sub-technique ID

**Output Format (Read Operations):**
- Returns `[]armis.MitreAttackLabel` structs with separate fields:
  ```go
  type MitreAttackLabel struct {
      Matrix       string // "Enterprise", "Mobile", or "ICS"
      Tactic       string // Human-readable name, e.g., "Initial Access"
      Technique    string // Human-readable name, e.g., "Phishing"
      SubTechnique string // Full ID format, e.g., "T1566.001"
  }
  ```

### Format Specification

The Armis input format combines standard MITRE ATT&CK components:

**Format:** `<Matrix>.<TacticID>.<TechniqueID>[.<SubTechniqueID>]`

**Components:**
- **Matrix**: One of `Enterprise`, `Mobile`, or `ICS` (MITRE ATT&CK framework domains)
- **Tactic ID**: `TA` followed by 4 digits (e.g., `TA0009` for Collection tactic)
- **Technique ID**: `T` followed by 4 digits (e.g., `T1056` for Input Capture)
- **Sub-technique ID**: Optional 3 digits (e.g., `001` for Keylogging)

**Examples:**
- `Enterprise.TA0009.T1056.001` = Enterprise Matrix > Collection Tactic > Input Capture Technique > Keylogging Sub-technique
- `Enterprise.TA0009.T1056.004` = Enterprise Matrix > Collection Tactic > Input Capture Technique > Credential API Hooking
- `Mobile.TA0001.T1234` = Mobile Matrix > Tactic TA0001 > Technique T1234 (no sub-technique)
- `ICS.TA0010.T0800.123` = ICS Matrix > Tactic TA0010 > Technique T0800 > Sub-technique 123

### Validator Regex

```go
^(Enterprise|Mobile|ICS)\.TA\d{4}\.T\d{4}(\.\d{3})?$
```

**Breakdown:**
- `^(Enterprise|Mobile|ICS)` - Matches one of the three MITRE ATT&CK matrices
- `\.TA\d{4}` - Dot followed by TA and exactly 4 digits (tactic ID)
- `\.T\d{4}` - Dot followed by T and exactly 4 digits (technique ID)
- `(\.\d{3})?` - Optional: dot followed by exactly 3 digits (sub-technique ID)
- `$` - End of string

**Validation:** ✅ The regex correctly validates the Armis-specific format.

### Evidence from Codebase

**Test File:** [internal/provider/policy_resource_test.go:31-32](internal/provider/policy_resource_test.go#L31-L32)
```go
resource.TestCheckResourceAttr(resourceName, "mitre_attack_labels.0", "Enterprise.TA0009.T1056.001"),
resource.TestCheckResourceAttr(resourceName, "mitre_attack_labels.1", "Enterprise.TA0009.T1056.004"),
```

**Example File:** [examples/resources/armis_policy/resource.tf:7](examples/resources/armis_policy/resource.tf#L7)
```hcl
mitre_attack_labels = ["Enterprise.TA0009.T1056.001", "Enterprise.TA0009.T1056.004"]
```

**Conversion Function:** [internal/utils/policy_utils.go:645](internal/utils/policy_utils.go#L645)
```go
MitreAttackLabels: ConvertListToStringSlice(model.MitreAttackLabels),
```

This converts the Terraform `[]types.String` to `[]string` for the API.

## Verification Against MITRE ATT&CK Standard

### Official MITRE Format

According to the [official MITRE ATT&CK documentation](https://attack.mitre.org/):

- **Tactic IDs**: `TA####` format (e.g., [TA0009 - Collection](https://attack.mitre.org/tactics/TA0009/))
- **Technique IDs**: `T####` format (e.g., [T1056 - Input Capture](https://attack.mitre.org/techniques/T1056/))
- **Sub-technique IDs**: `T####.###` format (e.g., [T1056.001 - Keylogging](https://attack.mitre.org/techniques/T1056/001/))
- **Matrices**: Enterprise, Mobile, ICS

### Armis vs. Standard MITRE Format

**Standard MITRE:**
- Stores components separately (not concatenated)
- References like: "Enterprise > TA0009 > T1056 > T1056.001"

**Armis Format:**
- Combines components with dots: `"Enterprise.TA0009.T1056.001"`
- This is an Armis-specific serialization format for API input

The Armis format uses standard MITRE IDs but combines them into a single string for convenience.

## Recommendations

### ✅ No Changes Required to Validator

The regex pattern is correct and accurately validates the Armis-specific MITRE label format used in the API.

### 📝 Documentation Enhancement

Update the validator function comment in [internal/verify/validate.go](internal/verify/validate.go) to be more detailed:

```go
// ValidMitreAttackLabel validates MITRE ATT&CK label format used by the Armis API.
// Expected format: <Matrix>.<TacticID>.<TechniqueID>[.<SubTechniqueID>]
// - Matrix: One of "Enterprise", "Mobile", or "ICS"
// - TacticID: "TA" followed by 4 digits (e.g., TA0009 for Collection)
// - TechniqueID: "T" followed by 4 digits (e.g., T1056 for Input Capture)
// - SubTechniqueID: Optional 3 digits (e.g., 001 for Keylogging)
//
// Examples:
//   - Enterprise.TA0009.T1056.001 (with sub-technique)
//   - Mobile.TA0001.T1234 (without sub-technique)
//   - ICS.TA0010.T0800.123
//
// Reference: https://attack.mitre.org/
func ValidMitreAttackLabel() validator.String {
```

### ✅ Unit Tests Still Required

As noted in the PR review, [internal/verify/validate_test.go](internal/verify/validate_test.go) should be created to test the validator logic.

## Sources

- [MITRE ATT&CK - Collection Tactic (TA0009)](https://attack.mitre.org/tactics/TA0009/)
- [MITRE ATT&CK - Input Capture Technique (T1056)](https://attack.mitre.org/techniques/T1056/)
- [MITRE ATT&CK - Keylogging Sub-technique (T1056.001)](https://attack.mitre.org/techniques/T1056/001/)
- [Armis Developer Documentation](https://docs.ic.armis.com/)
- [Armis MITRE ATT&CK Support](https://www.armis.com/solutions/mitre-attck-for-ics)

## Conclusion

**Status:** ✅ **VALIDATED**

The MITRE ATT&CK label validator regex is correct for the Armis API's expected input format. The format combines standard MITRE ATT&CK component IDs into a dot-separated string following the pattern `Enterprise.TA0009.T1056.001`.

No changes to the validation logic are required. The only improvement needed is enhanced documentation as outlined in the PR review.
