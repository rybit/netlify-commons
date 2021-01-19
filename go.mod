module github.com/netlify/netlify-commons

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/DataDog/datadog-go v4.2.0+incompatible
	github.com/armon/go-metrics v0.3.3
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bugsnag/bugsnag-go v1.5.3
	github.com/bugsnag/panicwrap v1.2.0 // indirect
	github.com/confluentinc/confluent-kafka-go v1.4.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/google/go-cmp v0.4.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/launchdarkly/eventsource v1.6.1 // indirect
	github.com/launchdarkly/go-test-helpers v1.2.0 // indirect
	github.com/nats-io/nats-streaming-server v0.15.1 // indirect
	github.com/nats-io/nats.go v1.8.1
	github.com/nats-io/stan.go v0.5.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.6.0
	github.com/satori/go.uuid v1.2.0
	github.com/sebest/xff v0.0.0-20160910043805-6c115e0ffa35
	github.com/shopify/logrus-bugsnag v0.0.0-20171204204709-577dee27f20d
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.4-0.20190321000552-67fc4837d267
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.0
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.4.1
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc
	gopkg.in/DataDog/dd-trace-go.v1 v1.26.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/ghodss/yaml.v1 v1.0.0 // indirect
	gopkg.in/launchdarkly/go-sdk-common.v1 v1.0.0-20200401173443-991b2f427a01 // indirect
	gopkg.in/launchdarkly/go-server-sdk.v4 v4.0.0-20200729232655-2a44fb361895
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace github.com/armon/go-metrics v0.3.3 => github.com/mheffner/go-metrics v0.3.4-0.20210119031049-ce109c593e89

go 1.14
