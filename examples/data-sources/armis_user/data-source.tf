# Read in user information
data "armis_user" "user" {
  email = armis_user.test_user.email
}
