# pbc
The Passbook Compiler

# Import Certs and create profile

pbc profile ls

This will list the cert profiles installed on the local machine

pbc profile add -p [profile-name] -c [path-to-cert+key-bundle]

This will ask you to enter the password to decrypt the cert bundle


pbc build -n [filePathAndNameOfPassbookFile] -p [nameOfProfileUsedToBuildPass]
