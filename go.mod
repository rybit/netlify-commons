module github.com/netlify/netlify-commons

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/armon/go-metrics v0.3.3
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bugsnag/bugsnag-go v1.5.3
	github.com/bugsnag/panicwrap v1.2.0 // indirect
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/joho/godotenv v1.3.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/nats-io/nats-streaming-server v0.15.1 // indirect
	github.com/nats-io/nats.go v1.8.1
	github.com/nats-io/stan.go v0.5.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pelletier/go-toml v1.3.0 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rs/cors v1.6.0
	github.com/satori/go.uuid v1.2.0
	github.com/sebest/xff v0.0.0-20160910043805-6c115e0ffa35
	github.com/shopify/logrus-bugsnag v0.0.0-20171204204709-577dee27f20d
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.4-0.20190321000552-67fc4837d267
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
	github.com/tinylib/msgp v1.1.0 // indirect
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.12.1
	gopkg.in/yaml.v3 v3.0.0-20190924164351-c8b7dadae555
)

go 1.13

replace github.com/armon/go-metrics v0.3.3 => github.com/mheffner/go-metrics v0.3.4-0.20200625203819-e0b1d18b0cdc
