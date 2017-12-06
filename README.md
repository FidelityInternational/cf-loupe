# CF Loupe

![loupe](assets/loupe.jpg "loupe")

CF Loupe is an App and Buildpack status dashboard

A 'Loupe' is a small magnifying glass used by jewelers and watchmakers to identify quality/defects of precious stones and materials.

## Requirements

- Go version 1.9 or higher
- [Ginkgo](https://onsi.github.io/ginkgo/)

## Checking out the code

```
git clone https://github.com/FidelityInternational/cf-loupe.git "${GOPATH}/src/github.com/FidelityInternational/cf-loupe"
cd "${GOPATH}/src/github.com/FidelityInternational/cf-loupe"
```

## Running the tests

```
export CF_USERNAME_1="YOUR-CF-USERNAME"
export CF_PASSWORD_1="YOUR-CF-PASSWORD"
export CF_API_1="YOUR-CF-API"
export CF_FOUNDATION_1="NAME-OF-FOUNDATION" (e.g. NP1, P1 etc)
ginkgo -r -cover
```

## Running on Cloud Foundry

Look in pass and use the admin test credentials (the paas service account)

```
cf push cf-loupe --no-start
cf set-env cf-loupe "CF_USERNAME_1" "YOUR-CF-USERNAME"
cf set-env cf-loupe "CF_PASSWORD_1" "YOUR-CF-PASSWORD"
cf set-env cf-loupe "CF_API_1" "YOUR-CF-API"
cf set-env cf-loupe "CF_FOUNDATION_1" "NAME-OF-FOUNDATION"
cf start cf-loupe
```

## Running Locally

Ensure that the repo is cloned into your `GOPATH`

```
export CF_USERNAME_1="YOUR-CF-USERNAME"
export CF_PASSWORD_1="YOUR-CF-PASSWORD"
export CF_API_1="YOUR-CF-API"
export CF_FOUNDATION_1="NAME-OF-FOUNDATION"
go run main.go router.go
```

## Multiple Cloud Foundries

`cf-loupe` searches for cloud foundry credentials in the environment. To see multiple Cloud Foundries on the dashboard set environment variables of the format `CF_FOUNDATION_X` containing the name of the foundation eg `CF_FOUNDATION_1=dev, CF_FOUNDATION_2=test, CF_FOUNDATION_3=prod`. Make sure that the credentials are also set for each foundation in the format `CF_USERNAME_X`, `CF_PASSWORD_X`, `CF_API_X`.

The number convention is such that we will have variables suffixed "_1", "_2" up till "_n" where n is the total number of foundations.
