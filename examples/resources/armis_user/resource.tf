resource "armis_user" "manager" {
  name = "Lab Manager"

  phone    = "867-5309"
  location = "Houston"
  username = "lab.manager@lab.com"
  email    = "lab.manager@lab.com"

  role_assignments = [{
    name  = ["Asset Manager", "User Manager", "Integrations Manager"]
    sites = ["Lab"]
  }]
}

