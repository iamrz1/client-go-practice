### Custom Resources in GO
First create a custom resource definition [file](/yaml/crontabClient.yaml) in yaml format.
Now an object of that type can me made using yaml [file](/yaml/crontabDefination.yaml)s.

To create a go client for custom object creation,  a few things are needed to be done.
1. A code generator is needed to be cloned in vendor/k8s.io.
2. Create [register.go](pkg/apis/examplecrd.com/register.go)  in pkg/apis/<group-name>/ directory.
3. Create [doc.go](pkg/apis/examplecrd.com/v1/doc.go), [register.go](pkg/apis/examplecrd.com/v1/register.go),
[types.go](pkg/apis/examplecrd.com/v1/types.go) in pkg/apis/<group-name>/<api-version>/ directory.
4. Now create a shell script [update-codegen.sh](hack/update-codegen.sh) in /hack directory.
run the script from root directory of the project:
``
$hack/update-codegen.sh 
``
5. Now objects of newly created custom type can be created from go program.

