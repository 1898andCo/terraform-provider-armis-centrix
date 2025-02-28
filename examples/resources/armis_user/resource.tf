resource "armis_user" "test_user" {
  name = "Test User"

  phone    = "8675309"
  location = "Houston"
  username = "test.user@test.com"
  email    = "test.user@test.com"

  role_assignments = [{
    name  = "Read Only"
    sites = ["Lab"]
  }]
}
