package service

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type Env string

const (
	envStAuth Env = "ST_AUTH"
	envStUser Env = "ST_USER"
	envStKey  Env = "ST_KEY"

	envAwsDefaultRegion   Env = "AWS_DEFAULT_REGION"
	envAwsRegion          Env = "AWS_REGION"
	envAwsAccessKeyId     Env = "AWS_ACCESS_KEY_ID"
	envAwsSecretAccessKey Env = "AWS_SECRET_ACCESS_KEY"
	envAwsSessionToken    Env = "AWS_SESSION_TOKEN"
)

type Envs []Env

func (e Envs) String() string {
	sb := &strings.Builder{}
	w := tabwriter.NewWriter(sb, 0, 8, 4, '\t', 0)
	for _, env := range e {
		_, _ = fmt.Fprintf(w, "\t%s\t%s\n", env, env.Desc())
	}
	_ = w.Flush()
	return sb.String()
}

var AllEnvs = Envs{
	envStAuth, envStUser, envStKey,
	envAwsDefaultRegion, envAwsRegion, envAwsAccessKeyId, envAwsSecretAccessKey, envAwsSessionToken,
}

var envDescriptions = map[Env]string{
	envStAuth: "Swift authentication URL (v1.0)",
	envStUser: "Swift username",
	envStKey:  "Swift key (password)",

	envAwsRegion:          fmt.Sprintf("AWS region (if %s not set)", envAwsDefaultRegion),
	envAwsDefaultRegion:   "AWS region",
	envAwsAccessKeyId:     "AWS access key ID (AKID)",
	envAwsSecretAccessKey: "AWS secret",
	envAwsSessionToken:    "AWS session token",
}

func (e Env) Get() string {
	return os.Getenv(string(e))
}

func (e Env) Desc() string {
	return envDescriptions[e]
}
