# tempulate

Template out a file from parmeter file inputs (YAML and JSON) using Golang templates. Useful for DevOps work e.g. generate terraform from YAML/JSON parameter files. Could have other uses - happy to hear about suggestions. PRs welcome.

## Build

There is a Makefile. This version uses Golang modules, so will require Golang V1.11 or above. To buld use the command

````make````

This will build a binary in the bin directory.

All of this should work with earlier versions of Golang, but you'll need to handle dependencies e.g. use go dep or similar tools.

## Usage

### Binary

Let's use the example of bootstrapping an AWS account with users. There are some who would do this using the AWS console for each new user, far easier to do this with code. You know that each user is created in exactly the correct way, and it is easy to add new users/take users out, and manage changes to group membership etc. Also makes it very easy to recreate exactly the same users in a new account.

We accomplish this with Hashicorp's Terraform, but you could use other tools if you wish.

Start with a template file (account.tf.tmpl), this sets the account password policy, and creates user accounts. Keybase is used to encrypt the users console password. This can be decrypted by the user securely.


```
resource "aws_iam_account_password_policy" "strict" {
  minimum_password_length        = "{{ config "account.min_password_length" }}"
  require_lowercase_characters   = true
  require_numbers                = true
  require_uppercase_characters   = true
  require_symbols                = true
  allow_users_to_change_password = true
}

/* Use groups as AWS limits to 10 policies per group - change these to match your policies */
resource "aws_iam_group" "devops" {
  name = "devops"
  path = "/devops/"
}

resource "aws_iam_group_policy_attachment" "devops_policy_full_access" {
  group      = "${aws_iam_group.devops.name}"
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}


{{ range $user := value "$.users" }}
resource "aws_iam_user" "{{ $user.user_name }}" {
        name = "{{ $user.user_name }}"
}

resource "aws_iam_user_login_profile" "{{ $user.user_name }}" {
         user = "${aws_iam_user.{{ $user.user_name }}.name}"
         pgp_key = "{{ $user.pgp_key }}"
}

resource "aws_iam_access_key" "{{ $user.user_name }}" {
        user = "${aws_iam_user.{{ $user.user_name }}.name}"
}

resource "aws_iam_user_group_membership" "{{ $user.user_name }}" {
  user = "${aws_iam_user.{{ $user.user_name }}.name}"
  groups = [
    "${aws_iam_group.devops.name}",
  ]
}

output "{{ $user.user_name }}_name" {
       value = "${aws_iam_user.{{ $user.user_name }}.name}"
}

output "{{ $user.user_name }}_arn" {
       value = "${aws_iam_user.{{ $user.user_name }}.arn}"
}

output "{{ $user.user_name }}_key_fingerprint" {
       value = "${aws_iam_user_login_profile.{{ $user.user_name }}.key_fingerprint}"
}

output "{{ $user.user_name }}_encrypted_password" {
       value = "${aws_iam_user_login_profile.{{ $user.user_name }}.encrypted_password}"
}

{{ end }}

```

This template is driven from a YAML file of paramters. As new users are required/old users need to be removed, then this file can be editied and held in source control for audit/managment.

```
account:
    min_password_length: 12


users:
    - my.user:
      user_name: lorcan_user
      pgp_key: key_base_user
    - another.user:
      user_name: another_user
      pgp_key: another_key_base_user
````

This could also be done with JSON if you'd prefer e.g:

```
{
  "account": {
    "min_password_length": 12
  },
  "users": {
    "my.user": {
      "user_name": "lorcan_user",
      "pgp_key": "key_base_user"
    },
    "aother.user": {
      "user_name": "another_user",
      "pgp_key": "another_key_base_user"
    }
  }
}
```

Finally to test this out, run:

```
bin/tempulate -t account.tf.tmpl -p account.yaml
```

This will output to STDOUT, to output to a file for running with Terraform, use:

```
bin/tempulate -t account.tf.tmpl -p account.yaml -o account.tf
```

The resulting file can be 'applied' using Terraform.

### Package

This functionality can also be used from within your own code using the tempulate/munge package.

The function MungeFile will allow you to pass a template (as string) and YAML/JSON parameter files. See the munch_test.go file for examples.
