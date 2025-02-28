resource "armis_policy" "example_policy" {
  name                = "BMS Security Alert Policy"
  description         = "This is an example security policy with all options."
  enabled             = true
  rule_type           = "ACTIVITY"
  labels              = ["Security"]
  mitre_attack_labels = ["Enterprise.TA0009.T1056.001", "Enterprise.TA0009.T1056.004"]

  actions = [
    {
      type = "alert"
      params = {
        severity = "high"
        title    = "BMS Security Alert"
        type     = "Security - Threat"
        consolidation = {
          amount = 2
          unit   = "Hours"
        }
      }
    }
  ]

  rules = {
    and = [
      "protocol:BMS",
    ]
  }
}
